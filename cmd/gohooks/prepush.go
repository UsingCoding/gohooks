package main

import (
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"

	"gohooks/pkg/vcs"
)

func executePrePush(ctx *cli.Context) error {
	if ctx.Args().Len() != 2 {
		return errors.New("not enough args to call pre-push hook")
	}

	gitRepo, err := vcs.ContextRepo(ctx.Context)
	if err != nil {
		return err
	}

	branch, err := gitRepo.GetCurrentBranch()
	if err != nil {
		return err
	}

	if branch != "master" {
		return nil
	}

	// Allow master push set up env var GOHOOK_UNPROTECT_MASTER
	if !ctx.Bool("unprotect-master") {
		return errors.New("cannot push to the protected master")
	}

	return nil
}
