package pprof

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/pprof"
	"time"

	"github.com/nekomeowww/insights-bot/pkg/logger"
	"go.uber.org/fx"
)

type NewPprofParams struct {
	fx.In

	Lifecycle fx.Lifecycle
	Logger    *logger.Logger
}

type Pprof struct {
	srv    *http.Server
	logger *logger.Logger
}

func NewPprof() func(NewPprofParams) *Pprof {
	return func(params NewPprofParams) *Pprof {
		srvMux := http.NewServeMux()

		srvMux.HandleFunc("/debug/pprof/", pprof.Index)
		srvMux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		srvMux.HandleFunc("/debug/pprof/profile", pprof.Profile)
		srvMux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		srvMux.HandleFunc("/debug/pprof/trace", pprof.Trace)

		srv := &http.Server{
			Addr:              ":6060",
			Handler:           srvMux,
			ReadHeaderTimeout: time.Second * 15,
		}

		params.Lifecycle.Append(fx.Hook{
			OnStop: func(ctx context.Context) error {
				closeCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()

				if err := srv.Shutdown(closeCtx); err != nil && err != http.ErrServerClosed {
					return err
				}

				return nil
			},
		})

		return &Pprof{
			srv:    srv,
			logger: params.Logger,
		}
	}
}

func Run() func(*Pprof) error {
	return func(srv *Pprof) error {
		listener, err := net.Listen("tcp", srv.srv.Addr)
		if err != nil {
			return fmt.Errorf("failed to listen %s: %v", "0.0.0.0:6060", err)
		}

		go func() {
			if err := srv.srv.Serve(listener); err != nil && err != http.ErrServerClosed {
				srv.logger.Fatalf("failed to serve pprof: %v", err)
			}
		}()

		return nil
	}
}
