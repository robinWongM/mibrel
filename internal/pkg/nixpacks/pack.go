package nixpacks

import (
	binding "github.com/zyreva/zyreva/internal/gen/uniffi/nixpacks"
)

func Build(dir string) binding.BuildPlanWithHashMap {
	return binding.GenerateBuildPlan(dir, []string{})
}
