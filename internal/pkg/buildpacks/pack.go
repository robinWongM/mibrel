package pack

import (
	"context"

	"github.com/buildpacks/pack/pkg/client"
	"github.com/buildpacks/pack/pkg/image"
)

func Build(dir string) string {
	c, err := client.NewClient()
	if err != nil {
		panic(err)
	}

	err = c.Build(context.Background(), client.BuildOptions{
		Builder:    "paketobuildpacks/builder-jammy-base",
		Image:      "zyreva",
		AppPath:    dir,
		PullPolicy: image.PullIfNotPresent,
	})
	if err != nil {
		panic(err)
	}

	return "zyreva"
}
