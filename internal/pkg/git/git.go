package git

import (
	"os"

	"github.com/go-git/go-git/v5"
)

func Clone(url string) string {
	// create a new empty random temp dir
	dir, err := os.MkdirTemp("", "zyreva-git-")
	if err != nil {
		panic(err)
	}

	_, err = git.PlainClone(dir, false, &git.CloneOptions{
		URL: url,
	})

	if err != nil {
		panic(err)
	}

	return dir
}
