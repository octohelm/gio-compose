package main

import (
	"context"
	"os"

	"github.com/innoai-tech/infra/pkg/cli"
	"github.com/octohelm/gio-compose/cmd/giocompose/internal"
	"github.com/octohelm/gio-compose/internal/version"
)

var App = cli.NewApp("giocompose", version.Version())

var workspace internal.Workspace

func main() {
	if w, err := internal.InitWorkspace(); err != nil {
		panic(err)
	} else {
		workspace = w
	}

	if err := cli.Execute(context.Background(), App, os.Args[1:]); err != nil {
		os.Exit(1)
	}
}
