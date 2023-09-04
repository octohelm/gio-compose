package main

import (
	"context"
	"log"
	"os"

	"gioui.org/app"

	"github.com/octohelm/gio-compose/pkg/component/m3/theming"
	"github.com/octohelm/gio-compose/pkg/component/m3/typescale"
	. "github.com/octohelm/gio-compose/pkg/compose"
	"github.com/octohelm/gio-compose/pkg/compose/modifier"
	"github.com/octohelm/gio-compose/pkg/compose/renderer"
)

var flags = &Flags{}

type Flags struct {
	WithGuideline          bool
	WithBoundingRectStdout bool
}

func (flags Flags) Options() (options []renderer.RootOptionFunc) {
	if flags.WithGuideline {
		options = append(options, renderer.WithGuideline(true))
	}

	if flags.WithBoundingRectStdout {
		options = append(options, renderer.WithBoundingRectPrinter(os.Stdout))
	}

	return
}

func main() {
	w := app.NewWindow()
	w.Option(app.Maximized.Option())

	r := renderer.CreateRoot(
		w,
		flags.Options()...,
	)

	ctx := context.Background()
	theme := theming.New()
	ctx = theming.Context.Inject(ctx, theme)

	el := H(
		Preview{},

		modifier.ProvideTextStyle(modifier.TextStyle(theme.TypeScale(typescale.BodyMedium))),
	)

	r.Render(ctx, el)

	go func() {
		if err := r.Loop(); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()

	app.Main()
}
