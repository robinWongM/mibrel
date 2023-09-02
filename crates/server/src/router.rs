use git2::Repository;
use rspc::{Router, RouterBuilder, Type};
use serde::Deserialize;
use tempfile;

#[derive(Type, Deserialize)]
struct CreateReq {
    git_url: String,
}

fn to_rspc_error(err: impl std::fmt::Display) -> rspc::Error {
    rspc::Error::new(rspc::ErrorCode::InternalServerError, err.to_string())
}

fn create_user_router() -> RouterBuilder {
    Router::new().mutation("create", |t| {
        t(|ctx, input: CreateReq| {
            let dir = tempfile::tempdir().map_err(to_rspc_error)?;
            println!("Created temp directory, {}", dir.path().display());

            let repo = Repository::clone(&input.git_url, dir.path()).map_err(to_rspc_error)?;

            println!("Cloned Git repository, {}", dir.path().display());

            // Get the latest commit message
            let head = repo.head().map_err(to_rspc_error)?;
            let head = head.peel_to_commit().map_err(to_rspc_error)?;
            let message = head.message().unwrap_or("No commit message");

            Ok(String::from(message))
        })
    })
}

pub fn create_router() -> Router {
    <Router>::new()
        .query("version", |t| t(|ctx, input: ()| env!("CARGO_PKG_VERSION")))
        .merge("apps.", create_user_router())
        .build()
}
