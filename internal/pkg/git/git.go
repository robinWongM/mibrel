package git

import (
	"os"

	"github.com/go-git/go-git/v5"
)

func Clone(url string) string {
	// create a temp dir
	dir := os.TempDir()

	_, err := git.PlainClone(dir, false, &git.CloneOptions{
		URL: url,
	})

	if err != nil {
		panic(err)
	}

	return dir
}
