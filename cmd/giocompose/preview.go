package main

import (
	"bytes"
	"context"
	"fmt"
	"go/build"
	"path/filepath"
	"strings"

	"github.com/octohelm/gio-compose/cmd/giocompose/internal"
	"mvdan.cc/sh/v3/interp"

	"github.com/innoai-tech/infra/pkg/cli"
)

func init() {
	cli.AddTo(App, &Preview{})
}

type Preview struct {
	cli.C
	Previewer
}

type Previewer struct {
	PkgPath string `arg:""`
	Name    string `arg:""`

	WithGuideline          bool `flag:"with-guideline,omitempty"`
	WithBoundingRectStdout bool `flag:"with-bounding-rect-stdout,omitempty"`
}

func (s *Previewer) Run(ctx context.Context) error {
	importPath := "."
	dir := "."

	if strings.HasPrefix(s.PkgPath, ".") {
		dir = workspace.WorkDir(s.PkgPath)
	} else {
		importPath = s.PkgPath
	}

	p, err := build.Import(importPath, dir, build.FindOnly)
	if err != nil {
		return err
	}

	previewCode := bytes.NewBuffer(nil)

	_, _ = fmt.Fprintf(previewCode, `package main

import p %q

type Preview = p.%s
`, p.ImportPath, s.Name)

	if s.WithGuideline {
		_, _ = fmt.Fprintf(previewCode, `
func init() {
	flags.WithGuideline = true
}
`)
	}

	if s.WithBoundingRectStdout {
		_, _ = fmt.Fprintf(previewCode, `
func init() {
	flags.WithBoundingRectStdout = true
}
`)
	}

	previewRoot := workspace.CacheDir("preview")

	if err := internal.WriteFile(filepath.Join(previewRoot, "main.go"), internal.PreviewDefaultTpl); err != nil {
		return err
	}
	if err := internal.WriteFile(filepath.Join(previewRoot, "preview.go"), previewCode.Bytes()); err != nil {
		return err
	}

	return internal.Run(ctx, "go run -tags=preview .", interp.Dir(previewRoot))
}
