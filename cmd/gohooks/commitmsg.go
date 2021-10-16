package main

import (
	"io"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"

	"gohooks/pkg/vcs"
)

func executeCommitMsg(ctx *cli.Context) error {
	if ctx.Args().Len() != 1 {
		return errors.New("not enough args to call commit-msg hook")
	}

	// Allow master push set up env var GOHOOK_UNPROTECT_MASTER
	if !ctx.Bool("unprotect-commit-message") {
		return errors.New("cannot push to the protected master")
	}

	commitMsgPath := ctx.Args().First()

	file, err := os.Open(commitMsgPath)
	if err != nil {
		return errors.Wrap(err, "failed to open commit msg file")
	}

	msgBytes, err := io.ReadAll(file)
	if err != nil {
		return errors.Wrap(err, "failed to read commit msg file")
	}

	msg := string(msgBytes)

	// Skip check for merge commit
	if strings.HasPrefix(msg, "Merge") {
		return nil
	}

	gitRepo, err := vcs.ContextRepo(ctx.Context)
	if err != nil {
		return err
	}

	branch, err := gitRepo.GetCurrentBranch()
	if err != nil {
		return err
	}

	if strings.HasPrefix(msg, branch) {
		return nil
	}

	return errors.Errorf("Commit message should start with branch name %s", branch)
}
