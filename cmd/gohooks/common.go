package main

import (
	"context"
	"os"
	"path"

	"github.com/pkg/errors"

	"github.com/UsingCoding/gohooks/pkg/config"
	"github.com/UsingCoding/gohooks/pkg/executor"
	"github.com/UsingCoding/gohooks/pkg/vcs"
)

// stopErr to stop execution
type stopErr struct{}

func (err stopErr) Error() string {
	return "execution stopped intentionally with no err"
}

func initVcs(ctx context.Context) (context.Context, error) {
	e, err := executor.New("git")
	if err != nil {
		return nil, err
	}

	git := vcs.NewGit(e)

	ctx = vcs.WithVCS(ctx, git)

	gitRepo, err := git.CurrentRepo()
	if err != nil {
		return nil, err
	}

	return vcs.WithRepo(ctx, gitRepo), nil
}

func loadConfig() (*config.Config, error) {
	homeDir := os.Getenv("HOME")
	if homeDir == "" {
		return nil, errors.New("$HOME env is not set")
	}

	return config.LoadConfig(path.Join(homeDir, configPath))
}
