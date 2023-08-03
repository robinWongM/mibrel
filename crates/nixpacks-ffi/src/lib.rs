use std::collections::{BTreeMap, HashMap};

use nixpacks::nixpacks::{plan::{phase::Phase, phase::StartPhase}, environment::EnvironmentVariables, app::StaticAssets};

type BTreeMapSS = BTreeMap<String, String>;

impl UniffiCustomTypeConverter for BTreeMapSS {
    type Builtin = HashMap<String, String>;

    fn into_custom(val: Self::Builtin) -> uniffi::Result<Self> {
        Ok(val.into_iter().collect())
    }

    fn from_custom(obj: Self) -> Self::Builtin {
        obj.into_iter().collect()
    }
}

pub struct BuildPlan {
    pub providers: Option<Vec<String>>,
    pub build_image: Option<String>,
    pub variables: Option<EnvironmentVariables>,
    pub static_assets: Option<StaticAssets>,
    pub phases: Option<HashMap<String, Phase>>,
    pub start_phase: Option<StartPhase>, 
}

pub fn get_plan_providers(path: String, envs: Vec<String>) -> Vec<String> {
    return nixpacks::get_plan_providers(
        path.as_str(),
        envs.iter().map(|x| x.as_str()).collect(),
        &nixpacks::nixpacks::plan::generator::GeneratePlanOptions {
            plan: None,
            config_file: None,
        },
    )
    .unwrap();
}

pub fn generate_build_plan(path: String, envs: Vec<String>) -> BuildPlan {
    let build_plan = nixpacks::generate_build_plan(
        path.as_str(),
        envs.iter().map(|x| x.as_str()).collect(),
        &nixpacks::nixpacks::plan::generator::GeneratePlanOptions {
            plan: None,
            config_file: None,
        },
    )
    .unwrap();

    let phases: Option<HashMap<String, Phase>> = Some(build_plan.phases.unwrap().into_iter().collect());

    return BuildPlan {
        providers: build_plan.providers,
        build_image: build_plan.build_image,
        variables: build_plan.variables,
        static_assets: build_plan.static_assets,
        phases,
        start_phase: build_plan.start_phase,
    };
}

uniffi::include_scaffolding!("nixpacks");
