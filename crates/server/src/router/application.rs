use std::{
    fs,
    path::{Path, PathBuf},
    sync::Arc, io::Stdout,
};

use git2::{build::RepoBuilder, FetchOptions};
use rspc::{Router, RouterBuilder, Type};
use serde::Deserialize;
use tempfile;
use tokio::{process::Command, io::AsyncBufReadExt, sync::mpsc::Sender};

use crate::entities::{prelude::*, *};
use sea_orm::*;

use super::Context;
use nixpacks::{
    create_docker_image, generate_build_plan,
    nixpacks::{
        builder::docker::DockerBuilderOptions,
        plan::{generator::GeneratePlanOptions, BuildPlan},
    },
};

#[derive(Type, Deserialize)]
struct CreateReq {
    git_url: String,
}

#[derive(Type, Deserialize)]
struct AnalyzeRequest {
    id: i32,
}

fn to_rspc_error(err: impl std::fmt::Display) -> rspc::Error {
    rspc::Error::new(rspc::ErrorCode::InternalServerError, err.to_string())
}

pub fn create_app_router() -> RouterBuilder<Context> {
    Router::new()
        .query("list", |t| {
            t(|ctx: Context, _: ()| async move {
                let res = Application::find()
                    .all(ctx.db.as_ref())
                    .await
                    .map_err(to_rspc_error)?;

                Ok(res)
            })
        })
        .mutation("create", |t| {
            t(|ctx: Context, input: CreateReq| async move {
                let entity = application::ActiveModel {
                    name: ActiveValue::Set(input.git_url.clone()),
                    git_url: ActiveValue::Set(input.git_url.clone()),
                    ..Default::default()
                };

                let res = Application::insert(entity)
                    .exec(ctx.db.as_ref())
                    .await
                    .map_err(to_rspc_error)?;

                Ok(res.last_insert_id)
            })
        })
        .query("analyze", |t| {
            t(|ctx: Context, input: AnalyzeRequest| async move {
                let entity = Application::find_by_id(input.id)
                    .one(ctx.db.as_ref())
                    .await
                    .map_err(to_rspc_error)?
                    .ok_or_else(|| {
                        rspc::Error::new(rspc::ErrorCode::NotFound, "Not found".to_string())
                    })?;

                let dir = tempfile::tempdir().or_else(|e| Err(to_rspc_error(e)))?;
                println!("Created temp directory, {}", dir.path().display());

                clone(dir.path(), &entity.git_url, None).map_err(to_rspc_error)?;

                generate_build_plan(
                    dir.path().to_str().unwrap(),
                    vec![],
                    &GeneratePlanOptions {
                        plan: Some(BuildPlan {
                            providers: Some(vec!["...".to_string(), "staticfile".to_string()]),
                            ..Default::default()
                        }),
                        ..Default::default()
                    },
                )
                .and_then(|p| p.to_toml())
                .map_err(to_rspc_error)
            })
        })
        .subscription("build", |t| {
            t(|ctx: Context, input: AnalyzeRequest| {
                async_stream::stream! {
                    let (tx, mut rx) = tokio::sync::mpsc::channel(32);

                    // Spawn the long-running operation as a separate task
                    tokio::spawn(async move {
                        let entity = Application::find_by_id(input.id)
                            .one(ctx.db.as_ref())
                            .await
                            .map_err(to_rspc_error)?
                            .ok_or_else(|| {
                                rspc::Error::new(rspc::ErrorCode::NotFound, "Not found".to_string())
                            })?;

                        tx.send(format!("Git URL: {}", entity.git_url)).await.unwrap();

                        let dir = tempfile::tempdir().or_else(|e| Err(to_rspc_error(e)))?;
                        tx.send(format!("Created temp directory, {}", dir.path().display())).await.unwrap();

                        clone(dir.path(), &entity.git_url, Some(&tx)).map_err(to_rspc_error)?;
                        tx.send(format!("Cloned")).await.unwrap();

                        let _ = create_docker_image(
                            dir.path().to_str().unwrap(),
                            vec![],
                            &GeneratePlanOptions::default(),
                            &DockerBuilderOptions {
                                out_dir: Some(dir.path().display().to_string()),
                                ..Default::default()
                            },
                        )
                        .await
                        .map_err(to_rspc_error)?;

                        tx.send(format!("Generated Dockerfile")).await.unwrap();

                        // Generate Docker image full URL
                        let image_url = format!("{}/{}:{}", "registry:5000", "test", "latest");

                        // Run command `buildctl`
                        let mut buildctl = Command::new("buildctl");
                        buildctl
                            .arg("--addr")
                            .arg("tcp://buildkitd:1234")
                            .arg("build")
                            .arg("--frontend")
                            .arg("dockerfile.v0")
                            .arg("--local")
                            .arg(format!("context={}", dir.path().display()))
                            .arg("--local")
                            .arg(format!("dockerfile={}", dir.path().join(".nixpacks").display()))
                            .arg("--output")
                            .arg(format!("type=image,name={},push=true,registry.insecure=true", image_url));

                        // Stream output from the command
                        let mut child = buildctl
                            .stdout(std::process::Stdio::piped())
                            .stderr(std::process::Stdio::piped())
                            .spawn()
                            .map_err(to_rspc_error)?;

                        let mut stdout = tokio::io::BufReader::new(child.stdout.take().unwrap()).lines();
                        let mut stderr = tokio::io::BufReader::new(child.stderr.take().unwrap()).lines();

                        loop {
                            tokio::select! {
                                // Read a line from stdout
                                Ok(line_result) = stdout.next_line() => {
                                    if let Some(line) = line_result {
                                        tx.send(line).await.expect("Failed to send to channel");
                                    } else {
                                        break;
                                    }
                                },
                                // Read a line from stderr
                                Ok(line_result) = stderr.next_line() => {
                                    if let Some(line) = line_result {
                                        tx.send(line).await.expect("Failed to send to channel");
                                    } else {
                                        break;
                                    }
                                },
                                // Break when both stdout and stderr are closed
                                else => {
                                    break;
                                }
                            }
                        }

                        tx.send(format!("Build complete")).await.unwrap();

                        Ok::<(), rspc::Error>(())
                    });

                    // In the main async stream, yield results from the channel
                    while let Some(progress) = rx.recv().await {
                        yield progress.to_string();
                    }

                    // The channel is closed, so the stream ends
                    println!("Build complete");
                }
            })
        })
}

fn clone(dir: &Path, url: &str, sender: Option<&Sender<String>>) -> Result<String, rspc::Error> {
    let mut remote_callbacks = git2::RemoteCallbacks::new();
    remote_callbacks.transfer_progress(|progress| {
        let log = format!(
            "[Git] {}: {}/{}, {} bytes received",
            progress.total_objects(),
            progress.indexed_objects(),
            progress.received_objects(),
            progress.received_bytes(),
        );
        println!("{}", log);

        if let Some(sender) = sender {
            let sender_clone = sender.clone();
            let log_clone = log.clone();
            tokio::spawn(async move {
                sender_clone.send(log_clone).await.unwrap();
            });
        }

        true
    });

    let mut fetch_options = FetchOptions::new();
    fetch_options.depth(1);
    fetch_options.remote_callbacks(remote_callbacks);

    // Fetch the remote repository, with --depth 1
    let mut builder = RepoBuilder::new();
    builder.fetch_options(fetch_options);

    // Clone the repository
    let repo = builder.clone(url, dir).map_err(to_rspc_error)?;
    println!("Cloned Git repository, {}", dir.display());

    // Get the latest commit message
    let head = repo.head().map_err(to_rspc_error)?;
    let head = head.peel_to_commit().map_err(to_rspc_error)?;
    let message = head.message().unwrap_or("No commit message");

    Ok(message.to_owned())
}
