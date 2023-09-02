use rspc::{Router, RouterBuilder, Type};
use serde::Deserialize;

#[derive(Type, Deserialize)]
struct CreateReq {
    git_url: String,
}

fn create_user_router() -> RouterBuilder {
    <Router>::new().mutation("create", |t| t(|ctx, input: CreateReq| input.git_url))
}

pub fn create_router() -> Router {
    <Router>::new()
        .query("version", |t| t(|ctx, input: ()| env!("CARGO_PKG_VERSION")))
        .merge("apps.", create_user_router())
        .build()
}
