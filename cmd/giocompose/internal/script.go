package internal

import (
	"context"
	"os"
	"strings"

	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"
)

func Run(ctx context.Context, script string, options ...interp.RunnerOption) (e error) {
	sh, err := syntax.NewParser().Parse(strings.NewReader(script), "")
	if err != nil {
		return err
	}

	runner, err := interp.New(append([]interp.RunnerOption{interp.StdIO(os.Stdin, os.Stdout, os.Stderr)}, options...)...)
	if err != nil {
		return err
	}
	return runner.Run(ctx, sh)
}
