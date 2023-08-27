package nixpacks

import (
	"os"
	"os/exec"
)

type Phase struct {
	DependsOn      []string `json:"dependsOn"`
	NixPkgs        []string `json:"nixPkgs"`
	NixLibraries   []string `json:"nixLibraries"`
	NixOverlays    []string `json:"nixOverlays"`
	NixpkgsArchive string   `json:"nixpkgsArchive"`
	AptPackages    []string `json:"aptPackages"`

	Cmds             []string `json:"cmds"`
	OnlyIncludeFiles []string `json:"onlyIncludeFiles"`
	CacheDirectories []string `json:"cacheDirectories"`
	Paths            []string `json:"paths"`
}

type StartPhase struct {
	Cmd             string   `json:"cmd"`
	RunImage        string   `json:"runImage"`
	OnlyIncludeFile []string `json:"onlyIncludeFiles"`
}

type BuildPlan struct {
	Providers    []string          `json:"providers"`
	BuildImage   string            `json:"buildImage"`
	Variables    map[string]string `json:"variables"`
	Phases       map[string]Phase  `json:"phases"`
	Start        StartPhase        `json:"start"`
	StaticAssets map[string]string `json:"staticAssets"`
}

func Plan(dir string) ([]byte, error) {
	// Run command `nixpacks` to plan the build.
	cmd := exec.Command("nixpacks", "plan", dir)
	// Collect all output from the command into a string.
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	return out, nil
}

func BuildDockerfile(dir string) string {
	// Create a temp dir
	tmpDir, err := os.MkdirTemp("", "zyreva")
	if err != nil {
		// If the command fails, print the error and exit.
		panic(err)
	}

	// Run command `nixpacks` to build the Dockerfile and save it to the temp dir.
	cmd := exec.Command("nixpacks", "build", dir, "-o", tmpDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()

	if err != nil {
		// If the command fails, print the error and exit.
		panic(err)
	}

	return tmpDir
}
