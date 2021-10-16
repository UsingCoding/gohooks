package notifications

import (
	"context"

	"github.com/pkg/errors"
)

type serviceKey struct{}

func WithService(ctx context.Context, service Service) context.Context {
	return context.WithValue(ctx, serviceKey{}, service)
}

func ContextService(ctx context.Context) (Service, error) {
	v, ok := ctx.Value(serviceKey{}).(Service)
	if !ok {
		return nil, errors.New("no notifications.Service in context")
	}

	return v, nil
}
