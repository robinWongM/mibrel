use std::path::{PathBuf, Path};

use git2::{build::RepoBuilder, FetchOptions};
use rspc::{Router, RouterBuilder, Type};
use serde::Deserialize;
use tempfile;

use crate::entities::{prelude::*, *};
use sea_orm::*;

use super::Context;
use nixpacks::{generate_build_plan, nixpacks::plan::generator::GeneratePlanOptions};

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
        .mutation("analyze", |t| {
            t(|ctx: Context, input: AnalyzeRequest| async move {
                let entity = Application::find_by_id(input.id)
                    .one(ctx.db.as_ref())
                    .await
                    .expect("Failed to execute query")
                    .expect("No application found");

                let dir = tempfile::tempdir().expect("Failed to create temp directory");
                println!("Created temp directory, {}", dir.path().display());

                clone(dir.path(), &entity.git_url).expect("Failed to clone repository");

                generate_build_plan(
                    dir.path().to_str().unwrap(),
                    vec![],
                    &GeneratePlanOptions::default(),
                )
                .and_then(|p| p.to_toml())
                .expect("Failed to generate build plan")
            })
        })
}

fn clone(dir: &Path, url: &str) -> Result<String, rspc::Error> {
    let mut remote_callbacks = git2::RemoteCallbacks::new();
    remote_callbacks.transfer_progress(|progress| {
        println!(
            "{}: {}/{}, {} bytes received",
            progress.total_objects(),
            progress.indexed_objects(),
            progress.received_objects(),
            progress.received_bytes(),
        );
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
