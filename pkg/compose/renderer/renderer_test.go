package renderer_test

import (
	"context"
	"fmt"
	"math/rand"
	"testing"

	"gioui.org/app"
	"github.com/octohelm/gio-compose/pkg/iter"
	"github.com/octohelm/gio-compose/pkg/text"

	"github.com/octohelm/gio-compose/pkg/component"
	. "github.com/octohelm/gio-compose/pkg/compose"
	"github.com/octohelm/gio-compose/pkg/compose/modifier"
	"github.com/octohelm/gio-compose/pkg/compose/renderer"
	. "github.com/octohelm/gio-compose/pkg/compose/testutil"
)

func TestRender(t *testing.T) {
	ctx := context.Background()

	t.Run("Direct render", func(t *testing.T) {
		r := renderer.CreateRoot(app.NewWindow())

		v := Box().Children(Fragment(
			Box(modifier.DisplayName("Child")),
		))

		r.Render(ctx, v)

		ExpectNodeRenderedEqual(t, r.RootNode(), "<Root><Box><Child/></Box></Root>")

		t.Run("Re-render should replace node", func(t *testing.T) {
			v := Box(modifier.DisplayName("ReplacedBox")).Children(Fragment(
				Box(modifier.DisplayName("Child")),
			))

			r.Render(ctx, v)

			ExpectNodeRenderedEqual(t, r.RootNode(), "<Root><ReplacedBox><Child/></ReplacedBox></Root>")
		})
	})

	t.Run("Render with Portal", func(t *testing.T) {
		r := renderer.CreateRoot(app.NewWindow())

		v := Box(modifier.DisplayName("App")).Children(
			Box(modifier.DisplayName("Normal")),
			component.RootPortal(modifier.Key("X")).Children(
				Box(modifier.DisplayName("InPortal")),
			),
		)

		r.Render(ctx, v)

		ExpectNodeRenderedEqual(t, r.WindowNode(), "<Window><Root><App><Normal/></App></Root><Portal><InPortal/></Portal></Window>")

		t.Run("Should replace node, when re-render", func(t *testing.T) {
			v := Box(modifier.DisplayName("App")).Children(
				Box(modifier.DisplayName("NormalReplaced")),
				component.RootPortal(modifier.Key("X")).Children(
					Box(modifier.DisplayName("InPortalReplaced")),
				),
			)

			r.Render(ctx, v)

			ExpectNodeRenderedEqual(t, r.WindowNode(), "<Window><Root><App><NormalReplaced/></App></Root><Portal><InPortalReplaced/></Portal></Window>")

			t.Run("Should remove node, when portal destroyed", func(t *testing.T) {
				v := Box(modifier.DisplayName("App")).Children(
					Box(modifier.DisplayName("NormalReplaced")),
				)

				r.Render(ctx, v)

				ExpectNodeRenderedEqual(t, r.WindowNode(), "<Window><Root><App><NormalReplaced/></App></Root></Window>")
			})
		})
	})

	t.Run("Re-render after diff", func(t *testing.T) {
		r := renderer.CreateRoot(app.NewWindow())

		r.Render(ctx, Box().Children(
			Box(modifier.DisplayName("Box1")),
			Box(modifier.DisplayName("Box2")),
		))
		ExpectNodeRenderedEqual(t, r.RootNode(), "<Root><Box><Box1/><Box2/></Box></Root>")

		t.Run("when move left", func(t *testing.T) {
			r.Render(ctx, Box().Children(
				Box(modifier.DisplayName("Box2")),
				Box(modifier.DisplayName("Box1")),
			))

			ExpectNodeRenderedEqual(t, r.RootNode(), "<Root><Box><Box2/><Box1/></Box></Root>")
		})

		t.Run("when move right", func(t *testing.T) {
			r.Render(ctx, Box().Children(
				Box(modifier.DisplayName("Box1")),
				Box(modifier.DisplayName("Box2")),
			))
			ExpectNodeRenderedEqual(t, r.RootNode(), "<Root><Box><Box1/><Box2/></Box></Root>")
		})

		t.Run("when with fragment", func(t *testing.T) {
			r.Render(ctx, Box().Children(
				Box(modifier.DisplayName("Box1")),
				Fragment(Box(modifier.DisplayName("Box2"))),
			))
			ExpectNodeRenderedEqual(t, r.RootNode(), "<Root><Box><Box1/><Box2/></Box></Root>")
		})

		t.Run("when insert before", func(t *testing.T) {
			r.Render(ctx, Box().Children(
				Box(modifier.DisplayName("Box1")),
				Box(modifier.DisplayName("Box3")),
				Fragment(Box(modifier.DisplayName("Box2"))),
			))
			ExpectNodeRenderedEqual(t, r.RootNode(), "<Root><Box><Box1/><Box3/><Box2/></Box></Root>")
		})

		t.Run("when insert after", func(t *testing.T) {
			r.Render(ctx, Box().Children(
				Box(modifier.DisplayName("Box1")),
				Box(modifier.DisplayName("Box3")),
				Fragment(Box(modifier.DisplayName("Box2"))),
				Box(modifier.DisplayName("Box4")),
			))
			ExpectNodeRenderedEqual(t, r.RootNode(), "<Root><Box><Box1/><Box3/><Box2/><Box4/></Box></Root>")
		})

		t.Run("when replace with", func(t *testing.T) {
			r.Render(ctx, Box().Children(
				Box(modifier.DisplayName("Box1")),
				Fragment(Box(modifier.DisplayName("Box3"))),
				Box(modifier.DisplayName("Box2")),
				Box(modifier.DisplayName("Box4")),
			))

			ExpectNodeRenderedEqual(t, r.RootNode(), "<Root><Box><Box1/><Box3/><Box2/><Box4/></Box></Root>")
		})

		t.Run("when replace by switch", func(t *testing.T) {
			r.Render(ctx, Box().Children(
				Box(modifier.DisplayName("Box1")),
				Box(modifier.DisplayName("Box3")),
				Box(modifier.DisplayName("Box2")),
				Fragment(Box(modifier.DisplayName("Box4"))),
			))

			ExpectNodeRenderedEqual(t, r.RootNode(), "<Root><Box><Box1/><Box3/><Box2/><Box4/></Box></Root>")
		})

	})

	t.Run("Render replaced with Provider Wrapper", func(t *testing.T) {
		r := renderer.CreateRoot(app.NewWindow())

		selectItem := func(selected int) {
			r.Render(ctx, Fragment(
				iter.MapIndexed(make([]int, 5), func(e int, i int) VNode {
					return Box(
						modifier.DisplayName(fmt.Sprintf("Box%d", i)),

						modifier.Size(40),
						modifier.When(i == selected, modifier.ProvideTextStyle(modifier.TextAlign(text.Middle))),
					)
				})...,
			))
		}

		t.Run("rand select", func(t *testing.T) {
			prev := -1

			selectItem(prev)
			ExpectNodeRenderedEqual(t, r.RootNode(), "<Root><Box0/><Box1/><Box2/><Box3/><Box4/></Root>")

			for idx := 0; idx < 1000; idx++ {
				selected := rand.Intn(5)
				selectItem(selected)
				ExpectNodeRenderedEqual(t, r.RootNode(), "<Root><Box0/><Box1/><Box2/><Box3/><Box4/></Root>")
				prev = selected
			}
		})
	})

	t.Run("Render wrapper", func(t *testing.T) {
		r := renderer.CreateRoot(app.NewWindow())

		r.Render(context.Background(), Column().Children(
			Box(),
			Box(
				modifier.When(true, modifier.ProvideTextStyle(modifier.TextAlign(text.Middle))),
			),
		))

		ExpectNodeRenderedEqual(t, r.RootNode(), `
<Root>
   <Column>
		<Box/>
		<Box/>
   </Column>
</Root>
`)

		r.Render(context.Background(), Column().Children(
			Box(),
			Box(
				modifier.When(true, modifier.ProvideTextStyle(modifier.TextAlign(text.Middle))),
			),
		))

		ExpectNodeRenderedEqual(t, r.RootNode(), `
<Root>
   <Column>
		<Box/>
		<Box/>
   </Column>
</Root>
`)
	})

}
