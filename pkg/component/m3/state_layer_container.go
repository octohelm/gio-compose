package m3

import (
	"image/color"

	. "github.com/octohelm/gio-compose/pkg/compose"
	"github.com/octohelm/gio-compose/pkg/compose/modifier"
)

type StateLayerContainer struct {
	Color color.NRGBA
}

func (stateLayer StateLayerContainer) Build(ctx BuildContext) VNode {
	return Box(
		ctx.Modifiers(),
		modifier.Clip(),
	).Children(
		Fragment(ctx.ChildVNodes()...),
		Box(
			modifier.DisplayName("StateLayer"),
			modifier.BackgroundColor(stateLayer.Color),
			modifier.FillMaxSize(),
		),
	)
}
