package tasks

import (
	"context"
	"log"

	"github.com/hibiken/asynq"
)

// A list of task types.
const (
	TypeGitClone = "git:clone"

	redisAddr = "redis:6379"
)

func NewClient() *asynq.Client {
	return asynq.NewClient(
		asynq.RedisClientOpt{Addr: redisAddr},
	)
}

func RunServer() {
	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisAddr},
		asynq.Config{
			// Specify how many concurrent workers to use
			Concurrency: 10,
			// Optionally specify multiple queues with different priority.
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
			// See the godoc for other configuration options
		},
	)

	// mux maps a type to a handler
	mux := asynq.NewServeMux()
	mux.HandleFunc(TypeGitClone, func(ctx context.Context, t *asynq.Task) error {
		log.Printf("Handling task: type=%s, payload=%s", t.Type(), t.Payload())

		handler, err := NewGitCloneTaskFromPayload(t.Payload())
		if err != nil {
			return err
		}

		return handler.Handle()
	})

	if err := srv.Run(mux); err != nil {
		log.Fatalf("could not run server: %v", err)
	}

	log.Printf("server stopped")
}
