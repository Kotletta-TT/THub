package main

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"os"
	"github.com/Kotletta-TT/THub/internal/controller/ws/node"
	"github.com/Kotletta-TT/THub/internal/controller/ws/user"
	"github.com/Kotletta-TT/THub/internal/usecases"
	"github.com/Kotletta-TT/THub/internal/usecases/repo"
	"github.com/Kotletta-TT/THub/logger"

	"go.uber.org/fx"
)

func main() {
	fx.New(CreateApp()).Run()
}

func CreateApp() fx.Option {
	return fx.Options(
		fx.Provide(
			NewHTTPServer,
			NewRouter,
			node.NewNodeRoutes,
			user.NewUserRoutes,
			usecases.NewConnectUseCase,
			usecases.NewDisconnectUseCase,
			usecases.NewTransferUseCase,
			usecases.NewListNodesUseCase,
			fx.Annotate(
				slog.New,
				fx.As(new(logger.Logger)),
			),
			fx.Annotate(
				NewSlogHandler,
				fx.As(new(slog.Handler)),
			),
			fx.Annotate(
				repo.NewNodeRepo,
				fx.As(new(usecases.NodeRepoGet)),
				fx.As(new(usecases.NodeRepoStore)),
				fx.As(new(usecases.NodeRepoList)),
				fx.As(new(usecases.NodeRepoRemove)),
			),
		),
		fx.Invoke(func(*http.Server) {}),
	)
}

func NewSlogHandler() *slog.JSONHandler {
	return slog.NewJSONHandler(os.Stdout, nil)
}

func NewRouter(nr *node.NodeRoutes, ur *user.UserRoutes) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("./client/")))
	mux.HandleFunc("/node/connect", nr.Connect)
	mux.HandleFunc("/user/connect", ur.TransferToNode)
	mux.HandleFunc("/user/node/list", ur.ListNodes)
	return mux
}

func NewHTTPServer(lc fx.Lifecycle, router *http.ServeMux) *http.Server {
	srv := &http.Server{Addr: ":8080", Handler: router}
	lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				ln, err := net.Listen("tcp", srv.Addr)
				if err != nil {
					return err
				}
				go srv.Serve(ln)
				return nil
			},
			OnStop: func(ctx context.Context) error {
				return srv.Shutdown(ctx)
			},
		})
	return srv
}
