package main

import (
	"context"
	"os"
	"os/exec"
	"path"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"

	"gohooks/pkg/common"
	"gohooks/pkg/config"
	"gohooks/pkg/vcs"
)

func shouldProtectRepo(ctx context.Context) (bool, error) {
	repo, err := vcs.ContextRepo(ctx)
	if err != nil {
		return false, err
	}

	conf, err := config.ContextConfig(ctx)
	if err != nil {
		return false, err
	}

	remoteBranchesMap, err := repo.GetRemotes()
	if err != nil {
		return false, err
	}

	for remoteName, remoteAddress := range remoteBranchesMap {
		regexpsToCheck, ok := conf.ProtectedRepoRegExps[remoteName]
		if !ok {
			// For this remote there is no record in config file
			continue
		}

		for _, regexp := range regexpsToCheck {
			matched := regexp.MatchString(remoteAddress)
			if matched {
				// If we matched any regexp by this remote address we should protect this repo
				return true, nil
			}
		}
	}

	return false, nil
}

func runSpecifiedHookFromRepo(ctx *cli.Context) error {
	if ctx.Args().Len() == 0 {
		return errors.New("invalid usage of runSpecifiedHookFromRepo without any args")
	}

	repo, err := vcs.ContextRepo(ctx.Context)
	if err != nil {
		return err
	}

	hookPath := path.Join(repo.HooksPath(), ctx.Args().First())

	exists, err := common.PathExists(hookPath)
	if err != nil {
		return err
	}

	if !exists {
		// Skip launching custom hook
		return nil
	}

	cmd := exec.CommandContext(ctx.Context, hookPath, ctx.Args().Tail()...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
