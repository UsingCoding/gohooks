package vcs

import (
	"context"

	"github.com/pkg/errors"
)

type vcsKey struct{}
type repoKey struct{}

func WithRepo(ctx context.Context, repo Repo) context.Context {
	return context.WithValue(ctx, repoKey{}, repo)
}

func ContextRepo(ctx context.Context) (Repo, error) {
	v, ok := ctx.Value(repoKey{}).(Repo)
	if !ok {
		return nil, errors.New("no repo in context")
	}

	return v, nil
}

func WithVCS(ctx context.Context, vcs VCS) context.Context {
	return context.WithValue(ctx, vcsKey{}, vcs)
}

func ContextVCS(ctx context.Context) (VCS, error) {
	v, ok := ctx.Value(vcsKey{}).(VCS)
	if !ok {
		return nil, errors.New("no vcs in context")
	}

	return v, nil
}
