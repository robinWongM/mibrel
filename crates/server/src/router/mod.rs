mod application;
use rspc::Router;

use std::sync::Arc;
use sea_orm::DatabaseConnection;

pub struct Context {
    pub db: Arc<DatabaseConnection>,
}

pub fn create_router() -> Router<Context> {
    Router::<Context>::new()
        .query("version", |t| t(|ctx, input: ()| env!("CARGO_PKG_VERSION")))
        .merge("apps.", application::create_app_router())
        .build()
}
