use git2::Repository;
use rspc::{Router, RouterBuilder, Type};
use serde::Deserialize;
use tempfile;

#[derive(Type, Deserialize)]
struct CreateReq {
    git_url: String,
}

fn create_user_router() -> RouterBuilder {
    <Router>::new().mutation("create", |t| {
        t(|ctx, input: CreateReq| {
            let dir = tempfile::tempdir().unwrap();
            println!("Created temp directory, {}", dir.path().display());

            let repo = match Repository::clone(&input.git_url, dir.path()) {
                Ok(repo) => repo,
                Err(err) => panic!("failed to clone: {}", err),
            };

            println!("Cloned Git repository, {}", dir.path().display());

            let head = repo.head().expect("No head found");
            println!("Head! {}", dir.path().display());

            let name = head.name().expect("No head name found");

            String::from(name)
        })
    })
}

pub fn create_router() -> Router {
    <Router>::new()
        .query("version", |t| t(|ctx, input: ()| env!("CARGO_PKG_VERSION")))
        .merge("apps.", create_user_router())
        .build()
}
