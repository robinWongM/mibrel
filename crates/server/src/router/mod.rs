mod application;
use bollard::Docker;
use rspc::Router;

use std::sync::Arc;
use sea_orm::DatabaseConnection;

pub struct Context {
    pub db: Arc<DatabaseConnection>,
    pub docker: Arc<Docker>,
}

pub fn create_router() -> Router<Context> {
    Router::<Context>::new()
        .query("version", |t| t(|ctx, input: ()| env!("CARGO_PKG_VERSION")))
        .merge("apps.", application::create_app_router())
        .build()
}
