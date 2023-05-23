package datastore

import (
	"context"
	"crypto/tls"
	"net"
	"time"

	"github.com/nekomeowww/insights-bot/internal/configs"
	"github.com/redis/rueidis"
	"go.uber.org/fx"
)

type NewRedisParams struct {
	fx.In

	Configs *configs.Config
}

type Redis struct {
	rueidis.Client
}

func NewRedis() func(NewRedisParams) (*Redis, error) {
	return func(params NewRedisParams) (*Redis, error) {
		var tlsConfig *tls.Config
		if params.Configs.Redis.TLSEnabled {
			tlsConfig = &tls.Config{
				MinVersion: tls.VersionTLS12,
			}
		}

		client, err := rueidis.NewClient(rueidis.ClientOption{
			TLSConfig:    tlsConfig,
			Username:     params.Configs.Redis.Username,
			Password:     params.Configs.Redis.Password,
			InitAddress:  []string{net.JoinHostPort(params.Configs.Redis.Host, params.Configs.Redis.Port)},
			SelectDB:     int(params.Configs.Redis.DB),
			DisableCache: !params.Configs.Redis.ClientCacheEnabled,
		})
		if err != nil {
			return nil, err
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
		defer cancel()

		err = client.Do(ctx, client.B().Ping().Build()).Error()
		if err != nil {
			return nil, err
		}

		return &Redis{Client: client}, nil
	}
}
