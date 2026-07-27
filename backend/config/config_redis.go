package config

import (
	"context"
	"crypto/tls"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func NewRedisConfig(addr, password string, userTls bool) (*redis.Client, error) {
	var tlsConfig *tls.Config
	if userTls {
		tlsConfig = &tls.Config{
			MinVersion: tls.VersionTLS13,
		}
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:      addr,
		Password:  password,
		DB:        0,
		TLSConfig: tlsConfig,
	})

	ctx := context.Background()
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("fail connecting redis: %w", err)
	}

	return rdb, nil

}
