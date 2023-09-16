mod application;
mod project;
use bollard::Docker;
use rspc::Router;

use std::sync::Arc;
use sea_orm::DatabaseConnection;

pub struct Context {
    pub db: Arc<DatabaseConnection>,
    pub docker: Arc<Docker>,
    pub kube_client: kube::Client,
}

pub fn create_router() -> Router<Context> {
    Router::<Context>::new()
        .query("version", |t| t(|_ctx, _input: ()| env!("CARGO_PKG_VERSION")))
        .merge("apps.", application::create_app_router())
        .merge("projects.", project::create_router())
        .build()
}
