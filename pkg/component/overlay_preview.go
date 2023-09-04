//go:build preview

package component

import (
	. "github.com/octohelm/gio-compose/pkg/compose"
	"github.com/octohelm/gio-compose/pkg/compose/modifier"
	"github.com/octohelm/gio-compose/pkg/event/gesture"
	"github.com/octohelm/gio-compose/pkg/layout/alignment"
	"github.com/octohelm/gio-compose/pkg/layout/position"
	"golang.org/x/image/colornames"
)

//go:generate giocompose preview . OverlayPreview
type OverlayPreview struct {
}

func (OverlayPreview) Build(c BuildContext) VNode {
	visible := UseState(c, false)

	return Column().Children(
		Box(
			modifier.Align(alignment.Center),
			modifier.Width(100),
			modifier.PaddingVertical(20),
			modifier.PaddingHorizontal(28),
			modifier.BackgroundColor(colornames.Lightcyan),

			modifier.DetectGesture(
				gesture.OnTap(func() {
					visible.UpdateFunc(func(visible bool) bool {
						return !visible
					})
				}),
			),

			Positioned(
				H(OverlayPreview{}),
				modifier.Position(position.Right),
				modifier.Visible(visible.Value()),
				modifier.OnVisibleChange(func(v bool) {
					visible.UpdateFunc(func(visible bool) bool {
						return v
					})
				}),
			),
		).Children(
			Box(
				modifier.Align(alignment.Center),
			).Children(
				Text("tap me"),
			),
		),
	)
}
