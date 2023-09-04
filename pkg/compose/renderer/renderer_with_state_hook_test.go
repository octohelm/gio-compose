package renderer_test

import (
	"context"
	"fmt"
	"testing"

	"gioui.org/app"
	. "github.com/octohelm/gio-compose/pkg/compose"
	"github.com/octohelm/gio-compose/pkg/compose/modifier"
	"github.com/octohelm/gio-compose/pkg/compose/renderer"
	. "github.com/octohelm/gio-compose/pkg/compose/testutil"
)

func TestRenderWithStateHook(t *testing.T) {

	t.Run("should re render when stage changed", func(t *testing.T) {
		r := renderer.CreateRoot(app.NewWindow())

		var updateValue *func(v string)

		r.Render(context.Background(), H(StatedComponent[string]{
			Value:       "StateHookInited",
			UpdateValue: &updateValue,
		}).Children(
			Box(modifier.DisplayName("Appended")),
		))

		ExpectNodeRenderedEqual(t, r.RootNode(), "<Root><Label><StateHookInited/><Appended/></Label><Input/></Root>")

		r.Act(func() {
			(*updateValue)("StateHookUpdated")
		})

		ExpectNodeRenderedEqual(t, r.RootNode(), "<Root><Label><StateHookUpdated/><Appended/></Label><Input/></Root>")

		for i := 0; i < 2; i++ {
			r.Act(func() {
				(*updateValue)(fmt.Sprintf("StateHookUpdatedAgent%d", i))
			})

			ExpectNodeRenderedEqual(t, r.RootNode(), fmt.Sprintf("<Root><Label><StateHookUpdatedAgent%d/><Appended/></Label><Input/></Root>", i))
		}
	})

	t.Run("nested", func(t *testing.T) {
		r := renderer.CreateRoot(app.NewWindow())

		var updateParent *func(v string)
		var updateChild *func(v string)

		r.Render(context.Background(), H(StatedComponent[string]{
			Value:       "Parent",
			UpdateValue: &updateParent,
		}).Children(
			H(StatedComponent[string]{
				Value:       "Child",
				UpdateValue: &updateChild,
			}),
		))

		ExpectNodeRenderedEqual(t, r.RootNode(), `
<Root>
	<Label>
		<Parent/>
		<Label>
			<Child/>
		</Label>
		<Input/>
	</Label>
	<Input/>
</Root>
`)

		r.Act(func() {
			(*updateChild)("ChildUpdated")
		})

		ExpectNodeRenderedEqual(t, r.RootNode(), `
<Root>
	<Label>
		<Parent/>
		<Label>
			<ChildUpdated/>
		</Label>
		<Input/>
	</Label>
	<Input/>
</Root>
`)

		r.Act(func() {
			(*updateParent)("ParentUpdated")
		})

		ExpectNodeRenderedEqual(t, r.RootNode(), `
<Root>
	<Label>
		<ParentUpdated/>
		<Label>
			<ChildUpdated/>
		</Label>
		<Input/>
	</Label>
	<Input/>
</Root>
`)
	})

}

type StatedComponent[T comparable] struct {
	Value       T
	UpdateValue **func(v T)
}

func (a StatedComponent[T]) Build(ctx BuildContext) VNode {
	state := UseState[T](ctx, a.Value)

	updateValue := func(v T) {
		state.Update(v)
	}

	if a.UpdateValue != nil {
		*a.UpdateValue = &updateValue
	}

	return Fragment().Children(
		Box(modifier.DisplayName("Label")).Children(
			Box(modifier.DisplayName(fmt.Sprint(state.Value()))),
			Fragment(ctx.ChildVNodes()...),
		),
		Box(modifier.DisplayName("Input")),
	)
}
