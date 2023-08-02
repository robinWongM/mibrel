pub fn get_plan_providers(path: String, envs: Vec<String>) -> Vec<String> {
    return nixpacks::get_plan_providers(
        path.as_str(),
        envs.iter().map(|x| x.as_str()).collect(),
        &nixpacks::nixpacks::plan::generator::GeneratePlanOptions {
            plan: None,
            config_file: None,
        },
    ).unwrap();
}

uniffi::include_scaffolding!("nixpacks");
