mod router;

use axum::{
    error_handling::HandleErrorLayer,
    http::{Method, StatusCode},
    routing::get,
};
use clap::Parser;
use router::create_router;
use std::{
    net::{Ipv4Addr, SocketAddr},
    path::PathBuf,
};
use tower_http::cors::{Any, CorsLayer};
use tokio::signal;

/// Search for a pattern in a file and display the lines that contain it.
#[derive(Parser)]
struct Cli {
    #[arg(long)]
    generate_bindings: bool,
}

#[tokio::main]
async fn main() {
    let router = create_router().arced();

    #[cfg(all(debug_assertions, not(feature = "k8s")))] // Only export in development builds
    router
        .export_ts(PathBuf::from(env!("CARGO_MANIFEST_DIR")).join("../../packages/rspc/index.ts"))
        .unwrap();

    let cors = CorsLayer::new()
        // allow `GET` and `POST` when accessing the resource
        .allow_methods([Method::GET, Method::POST])
        // allow requests from any origin
        .allow_origin(Any)
        .allow_headers(Any);

    let app = axum::Router::new()
        .route("/", get(|| async { "Hello 'rspc~'!" }))
        .nest("/rspc", router.endpoint(|| ()).axum())
        .layer(cors)
        .layer(HandleErrorLayer::new(|error| async move {
            (
                StatusCode::INTERNAL_SERVER_ERROR,
                format!("Unhandled internal error: {}", error),
            )
        }));

    let addr = SocketAddr::from((Ipv4Addr::UNSPECIFIED, 3000));

    axum::Server::bind(&addr)
        .serve(app.into_make_service())
        .with_graceful_shutdown(shutdown_signal())
        .await
        .unwrap();
}

async fn shutdown_signal() {
    let ctrl_c = async {
        signal::ctrl_c()
            .await
            .expect("failed to install Ctrl + C handler");
    };

    #[cfg(unix)]
    let terminate = async {
        signal::unix::signal(signal::unix::SignalKind::terminate())
            .expect("failed to install signal handler")
            .recv()
            .await;
    };

    #[cfg(not(unix))]
    let terminate = std::future::pending::<()>();

    tokio::select! {
        _ = ctrl_c => {},
        _ = terminate => {},
    }

    println!("signal received, starting graceful shutdown");
}