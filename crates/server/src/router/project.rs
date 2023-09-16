use k8s_openapi::api::core::v1::Namespace;
use kube::{api::ListParams, Api};
use rspc::{Router, RouterBuilder, Type};
use serde::Serialize;

use super::Context;

fn to_rspc_err(err: impl std::fmt::Display) -> rspc::Error {
    rspc::Error::new(rspc::ErrorCode::InternalServerError, err.to_string())
}

#[derive(Serialize, Type)]
struct Project {
    id: String,
    name: Option<String>,
}

pub fn create_router() -> RouterBuilder<Context> {
    Router::new().query("list", |t| {
        t(|ctx: Context, _: ()| async move {
            let lp = ListParams::default();
            let ns = Api::<Namespace>::all(ctx.kube_client)
                .list(&lp)
                .await
                .map_err(to_rspc_err)?;

            let ns: Vec<_> = ns
                .items
                .into_iter()
                .filter_map(|item| {
                    let id = item.metadata.name?;
                    if id.starts_with("kube-") {
                        return None;
                    }

                    let name = item
                        .metadata
                        .annotations
                        .and_then(|a| a.get("mibrel.dev/name").cloned());
                    Some(Project { id, name })
                })
                .collect();

            Ok(ns)
        })
    })
}
