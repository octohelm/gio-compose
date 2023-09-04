package renderer_test

import (
	"context"
	"fmt"
	"image/color"
	"testing"

	modifierapi "github.com/octohelm/gio-compose/pkg/modifier"

	"gioui.org/app"

	. "github.com/octohelm/gio-compose/pkg/compose"
	"github.com/octohelm/gio-compose/pkg/compose/modifier"
	"github.com/octohelm/gio-compose/pkg/compose/renderer"
	. "github.com/octohelm/gio-compose/pkg/compose/testutil"
)

func TestComponent(t *testing.T) {
	t.Run("Component", func(t *testing.T) {
		ctx := context.Background()

		t.Run("render", func(t *testing.T) {
			r := renderer.CreateRoot(app.NewWindow())

			r.Render(ctx, H(App{Name: "Created"}))
			ExpectNodeRenderedEqual(t, r.WindowNode(), "<Window><Root><AppRoot><Box><Created0/></Box></AppRoot></Root></Window>")

			t.Run("when prop changed, should render correct value", func(t *testing.T) {
				r.Render(ctx, H(App{Name: "Changed"}))
				r.Act(func() {})

				ExpectNodeRenderedEqual(t, r.WindowNode(), "<Window><Root><AppRoot><Box><Changed3/></Box></AppRoot></Root></Window>")
			})
		})
	})

	t.Run("Modifier Component", func(t *testing.T) {
		t.Run("direct render", func(t *testing.T) {
			r := renderer.CreateRoot(app.NewWindow())

			r.Render(context.Background(), Box().Children(
				Box(
					modifier.DisplayName("Child"),

					Wrap("W0"),
					Wrap("W1"),
				),
			))

			ExpectNodeRenderedEqual(t, r.RootNode(), "<Root><Box><W0><W1><Child/></W1></W0></Box></Root>")
		})
	})

	t.Run("Context", func(t *testing.T) {
		r := renderer.CreateRoot(app.NewWindow())

		r.Render(context.Background(), Box(
			modifier.ProvideTextStyle(modifier.FontSize(16)),
		).Children(
			H(txt{}),
		))

		ExpectNodeRenderedEqual(t, r.RootNode(), "<Root><Box><T.16/></Box></Root>")

		r.Render(context.Background(), Box(
			modifier.ProvideTextStyle(modifier.When(true, modifier.FontSize(10))),
		).Children(
			Row().Children(
				H(txt{}),
			),
		))

		ExpectNodeRenderedEqual(t, r.RootNode(), "<Root><Box><Row><T.10/></Row></Box></Root>")
	})
}

type SubComp struct {
	Name string
}

func (a SubComp) Build(ctx BuildContext) VNode {
	state := UseState(ctx, 0)

	UseEffect(ctx, func() func() {
		for i := 0; i < 3; i++ {
			state.UpdateFunc(func(prev int) int {
				return prev + 1
			})
		}
		return nil
	}, []any{})

	return Box().Children(
		Box(
			modifier.DisplayName(fmt.Sprintf("%s%d", a.Name, state.Value())),
			modifier.When(state.Value() > 0, modifier.ProvideTextStyle(modifier.TextColor(color.NRGBA{}))),
		),
	)
}

type App struct {
	Name string
}

func (a App) Build(ctx BuildContext) VNode {
	return Box(modifier.DisplayName("AppRoot")).
		Children(
			H(SubComp{Name: a.Name}),
		)
}

func Wrap(name string) interface {
	modifierapi.Modifier[any]
	Component
} {
	return wrapper{name: name}
}

type wrapper struct {
	modifierapi.Discord

	name   string
	target VNode
}

func (w wrapper) Wrap(n VNode) Component {
	w.target = n
	return w
}

func (w wrapper) Build(ctx BuildContext) VNode {
	return Box(modifier.DisplayName(w.name)).Children(
		w.target.Children(ctx.ChildVNodes()...),
	)
}

type txt struct {
}

func (w txt) Build(ctx BuildContext) VNode {
	return Box(modifier.DisplayName(fmt.Sprintf("T.%v", TextStyleContext.Extract(ctx).FontSize)))
}
