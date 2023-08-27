package git

import (
	"io"
	"os"

	"github.com/go-git/go-git/v5"
)

func Clone(url string, log io.Writer) string {
	// create a new empty random temp dir
	dir, err := os.MkdirTemp("", "zyreva-git-")
	if err != nil {
		panic(err)
	}

	_, err = git.PlainClone(dir, false, &git.CloneOptions{
		URL:      url,
		Progress: log,
	})

	if err != nil {
		panic(err)
	}

	return dir
}
