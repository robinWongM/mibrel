package nixpacks

import (
	"fmt"

	binding "github.com/zyreva/zyreva/internal/gen/uniffi/nixpacks"
)

func Build(dir string) {
	result := binding.GetPlanProviders(dir, []string{})

	fmt.Printf("%v", result)
}
