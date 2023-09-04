package main

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/octohelm/gio-compose/pkg/component/m3/colorrole"

	"gioui.org/app"

	"github.com/octohelm/gio-compose/pkg/component/m3/theming"
	"github.com/octohelm/gio-compose/pkg/component/m3/typescale"
	. "github.com/octohelm/gio-compose/pkg/compose"
	"github.com/octohelm/gio-compose/pkg/compose/modifier"
	"github.com/octohelm/gio-compose/pkg/compose/renderer"
	"github.com/octohelm/gio-compose/pkg/event/gesture"
	"github.com/octohelm/gio-compose/pkg/layout/alignment"
)

var showcases = map[string]Component{}

func AddShowCase(name string, c Component) {
	showcases[name] = c
}

type Showcase struct {
}

func (Showcase) Build(b BuildContext) VNode {
	tm := theming.Context.Extract(b)

	selected := UseState(b, "")
	hovered := UseState(b, false)

	return Row(
		modifier.FillMaxSize(),
	).Children(
		Column(
			modifier.DetectGesture(
				gesture.OnMouseEnter(func() {
					hovered.Update(true)
				}),
				gesture.OnMouseLeave(func() {
					hovered.Update(false)
				}),
			),

			modifier.BackgroundColor(tm.Color(colorrole.SurfaceContainer)),
			modifier.When(hovered.Value(),
				modifier.Shadow(8),
			),

			modifier.Width(240),
			modifier.PaddingVertical(24),
			modifier.RoundedRight(10),
			modifier.FillMaxHeight(),
		).Children(
			Fragment(SortedMap(showcases, func(e Component, name string) VNode {
				return Row(
					modifier.DisplayName(name),
					modifier.DetectGesture(gesture.OnTap(func() {
						selected.Update(name)
					})),
					modifier.When[any](name == selected.Value(),
						modifier.BackgroundColor(tm.Color(colorrole.PrimaryContainer)),
						modifier.ProvideTextStyle(modifier.TextColor(tm.Color(colorrole.OnPrimaryContainer))),
					),
					modifier.FillMaxWidth(),
					modifier.PaddingHorizontal(24),
					modifier.Height(40),
					modifier.Align(alignment.Center),
				).Children(
					Text(strings.ToUpper(name)),
				)
			})...),
		),
		Box(
			modifier.Weight(1),
			modifier.FillMaxHeight(),
			modifier.PaddingAll(20),
		).Children(
			SafeBuild(showcases[selected.Value()], func(c Component) VNode {
				return H(c)
			}),
		),
	)
}

func main() {
	w := app.NewWindow()
	w.Option(app.Maximized.Option())

	r := renderer.CreateRoot(
		w,
		//renderer.WithGuideline(true),
		//renderer.WithBoundingRectPrinter(os.Stdout),
	)

	ctx := context.Background()
	theme := theming.New()
	ctx = theming.Context.Inject(ctx, theme)

	el := H(
		Showcase{},

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
