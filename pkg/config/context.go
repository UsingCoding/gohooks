package config

import (
	"context"

	"github.com/pkg/errors"
)

type configKey struct{}

func WithConfig(ctx context.Context, config *Config) context.Context {
	return context.WithValue(ctx, configKey{}, config)
}

func ContextConfig(ctx context.Context) (*Config, error) {
	v, ok := ctx.Value(configKey{}).(*Config)
	if !ok {
		return nil, errors.New("no *config.Config in context")
	}

	return v, nil
}
