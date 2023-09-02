use git2::{build::RepoBuilder, FetchOptions};
use rspc::{Router, RouterBuilder, Type};
use serde::Deserialize;
use tempfile;

use crate::entities::{prelude::*, *};
use sea_orm::*;

use super::Context;

#[derive(Type, Deserialize)]
struct CreateReq {
    git_url: String,
}

fn to_rspc_error(err: impl std::fmt::Display) -> rspc::Error {
    rspc::Error::new(rspc::ErrorCode::InternalServerError, err.to_string())
}

pub fn create_user_router() -> RouterBuilder<Context> {
    Router::new().mutation("create", |t| {
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
}

fn clone(input: CreateReq) -> Result<String, rspc::Error> {
    let dir = tempfile::tempdir().map_err(to_rspc_error)?;
    println!("Created temp directory, {}", dir.path().display());

    // Fetch the remote repository, with --depth 1
    let mut builder = RepoBuilder::new();
    let mut fetch_options = FetchOptions::new();
    fetch_options.depth(1);
    builder.fetch_options(fetch_options);

    // Clone the repository
    let repo = builder
        .clone(&input.git_url, dir.path())
        .map_err(to_rspc_error)?;
    println!("Cloned Git repository, {}", dir.path().display());

    // Get the latest commit message
    let head = repo.head().map_err(to_rspc_error)?;
    let head = head.peel_to_commit().map_err(to_rspc_error)?;
    let message = head.message().unwrap_or("No commit message");

    Ok(String::from(message))
}
