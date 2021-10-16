//nolint
package executor

import (
	"os/exec"
)

type Executor interface {
	Output(args []string, opts ...Opt) ([]byte, error)
	Run(args []string, opts ...Opt) error
}

func New(executable string) (Executor, error) {
	_, err := exec.LookPath(executable)
	if err != nil {
		return nil, err
	}

	return &executor{executable: executable}, nil
}

type executor struct {
	executable string
}

func (e *executor) Run(args []string, opts ...Opt) error {
	cmd := exec.Command(e.executable, args...)
	for _, opt := range opts {
		opt.apply(cmd)
	}
	return cmd.Run()
}

func (e *executor) Output(args []string, opts ...Opt) ([]byte, error) {
	cmd := exec.Command(e.executable, args...)
	for _, opt := range opts {
		opt.apply(cmd)
	}
	return cmd.Output()
}
