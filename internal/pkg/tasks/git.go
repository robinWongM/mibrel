package tasks

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/hibiken/asynq"
)

type GitClonePayload struct {
	URL string
}

type GitCloneTask struct {
	Payload GitClonePayload
	LogPath string

	logger *log.Logger
}

func NewGitCloneTask(url string) (*asynq.Task, error) {
	payload, err := json.Marshal(GitClonePayload{URL: url})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeGitClone, payload, asynq.Retention(24*time.Hour)), nil
}

func NewGitCloneTaskFromPayload(payload []byte) (*GitCloneTask, error) {
	var p GitClonePayload
	if err := json.Unmarshal(payload, &p); err != nil {
		return nil, err
	}

	logFile, err := os.CreateTemp("/tmp/zyreva/logs/", "git-clone-*.log")
	if err != nil {
		return nil, err
	}

	fmt.Println("logFile.Name()", logFile.Name())

	return &GitCloneTask{
		Payload: p,
		LogPath: logFile.Name(),

		logger: log.New(logFile, "zyreva: ", log.LstdFlags),
	}, nil
}

func (t *GitCloneTask) Handle() error {
	t.logger.Printf("Cloning repo: url=%s", t.Payload.URL)

	// create a new empty random temp dir
	dir, err := os.MkdirTemp("", "zyreva-git-")
	if err != nil {
		panic(err)
	}

	_, err = git.PlainClone(dir, false, &git.CloneOptions{
		URL:      t.Payload.URL,
		Progress: t.logger.Writer(),
	})

	if err != nil {
		panic(err)
	}

	return nil
}
