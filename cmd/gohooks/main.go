package main

import (
	"context"
	"log"
	"os"
	"os/exec"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"

	"github.com/UsingCoding/gohooks/pkg/config"
	"github.com/UsingCoding/gohooks/pkg/notifications"
)

const (
	configPath = ".go-hooks/config.yaml"
)

func main() {
	ctx := context.Background()

	app := &cli.App{
		Name:  "go-hooks",
		Usage: "Git hooks written in go",
		Before: func(c *cli.Context) error {
			service, err := notifications.NewNotifySendService()
			if err != nil {
				return err
			}

			c.Context = notifications.WithService(c.Context, service)

			conf, err := loadConfig()
			if err != nil {
				return err
			}

			c.Context = config.WithConfig(c.Context, conf)

			return nil
		},
		ExitErrHandler: errHandler,
		Commands: []*cli.Command{
			{
				Name: "hook",
				Before: func(c *cli.Context) (err error) {
					c.Context, err = initVcs(c.Context)
					if err != nil {
						return err
					}

					err = runSpecifiedHookFromRepo(c)
					if err != nil {
						return errors.Wrap(err, "hook from repo exited with non-zero code")
					}

					should, err := shouldProtectRepo(c.Context)
					if err != nil {
						return err
					}

					if !should {
						return &stopErr{}
					}

					return err
				},
				Subcommands: []*cli.Command{
					{
						Name:   "commit-msg",
						Usage:  "Check that commit message starts with branch name. Disable via GOHOOKS_UNPROTECT_COMMIT_MESSAGE=1",
						Action: executeCommitMsg,
						Flags: []cli.Flag{
							&cli.BoolFlag{
								Name:    "unprotect-commit-message",
								EnvVars: []string{"GOHOOKS_UNPROTECT_COMMIT_MESSAGE"},
							},
						},
					},
					{
						Name:   "pre-push",
						Usage:  "Protect pushing to master. Disable via GOHOOKS_UNPROTECT_MASTER=1",
						Action: executePrePush,
						Flags: []cli.Flag{
							&cli.BoolFlag{
								Name:    "unprotect-master",
								EnvVars: []string{"GOHOOKS_UNPROTECT_MASTER"},
							},
						},
					},
				},
			},
		},
	}

	err := app.RunContext(ctx, os.Args)
	if err != nil {
		// Ignore stopErr to be logged
		if _, ok := err.(*stopErr); ok {
			return
		}
		if exitErr, ok := err.(*exec.ExitError); ok {
			log.Fatal(string(exitErr.Stderr))
		}
		log.Fatal(err)
	}
}

func errHandler(ctx *cli.Context, err error) {
	// in some cases context may be nil
	if ctx.Context == nil {
		return
	}

	if _, ok := err.(*stopErr); ok {
		return
	}

	notificationsService, err2 := notifications.ContextService(ctx.Context)
	if err2 != nil {
		// We can`t report user about error
		return
	}

	// Ignore error cause we can`t report user about error
	_ = notificationsService.Send(err.Error())
}
