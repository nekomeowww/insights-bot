package pprof

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/pprof"
	"time"

	"github.com/nekomeowww/insights-bot/pkg/healthchecker"
	"github.com/nekomeowww/insights-bot/pkg/logger"
	"github.com/samber/lo"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type NewPprofParams struct {
	fx.In

	Lifecycle fx.Lifecycle
	Logger    *logger.Logger
}

var _ healthchecker.HealthChecker = (*Pprof)(nil)

type Pprof struct {
	srv        *http.Server
	srvStarted bool
	logger     *logger.Logger
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

func (p *Pprof) Check(ctx context.Context) error {
	return lo.Ternary(p.srvStarted, nil, fmt.Errorf("pprof server is not started yet"))
}

func Run() func(*Pprof) error {
	return func(srv *Pprof) error {
		listener, err := net.Listen("tcp", srv.srv.Addr)
		if err != nil {
			return fmt.Errorf("failed to listen %s: %v", "0.0.0.0:6060", err)
		}

		go func() {
			if err := srv.srv.Serve(listener); err != nil && err != http.ErrServerClosed {
				srv.logger.Fatal("failed to serve pprof", zap.Error(err))
			}
		}()

		srv.srvStarted = true

		return nil
	}
}
