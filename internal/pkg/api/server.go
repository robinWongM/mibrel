package api

import (
	"net/http"

	"github.com/rs/cors"
	"github.com/zyreva/zyreva/internal/gen/proto/zyreva/v1/zyrevav1connect"
	"github.com/zyreva/zyreva/internal/pkg/tasks"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type Server struct {
	Handler http.Handler
}

func NewServer() *Server {
	appService := &AppService{
		client: tasks.NewClient(),
	}

	mux := http.NewServeMux()
	path, handler := zyrevav1connect.NewAppServiceHandler(appService)
	mux.Handle(path, handler)

	return &Server{
		Handler: h2c.NewHandler(cors.New(cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowCredentials: true,
			AllowedHeaders:   []string{"*"},
		}).Handler(mux), &http2.Server{}),
	}
}
