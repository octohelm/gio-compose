//go:build preview

package component

import (
	. "github.com/octohelm/gio-compose/pkg/compose"
	"github.com/octohelm/gio-compose/pkg/compose/modifier"
	"github.com/octohelm/gio-compose/pkg/event/gesture"
	"github.com/octohelm/gio-compose/pkg/iter"
	"github.com/octohelm/gio-compose/pkg/layout/alignment"
	"github.com/octohelm/gio-compose/pkg/layout/position"
	modifierapi "github.com/octohelm/gio-compose/pkg/modifier"
	"golang.org/x/image/colornames"
)

//go:generate giocompose preview . PositionPreview
type PositionPreview struct {
}

func (PositionPreview) Build(c BuildContext) VNode {
	visible := UseState(c, false)

	contents := Box(
		modifier.BackgroundColor(colornames.Beige),
		modifier.PaddingHorizontal(24),
		modifier.PaddingVertical(12),
	).Children(
		Text("Tooltip"),
	)

	return Box(
		modifier.FillMaxSize(),
	).Children(
		Box(
			modifier.Align(alignment.Center),
			modifier.Size(400),
			modifier.BackgroundColor(colornames.Lightcyan),

			modifier.DetectGesture(
				gesture.OnMouseEnter(func() {
					visible.UpdateFunc(func(prev bool) bool {
						return true
					})
				}),
				gesture.OnMouseLeave(func() {
					visible.UpdateFunc(func(prev bool) bool {
						return false
					})
				}),
			),
			modifierapi.Modifiers(iter.Map(position.AllFour, func(p position.Position) modifierapi.Modifier[any] {
				return modifierapi.Modifiers(
					iter.Map([]alignment.Alignment{alignment.Start, alignment.Middle, alignment.End}, func(a alignment.Alignment) modifierapi.Modifier[any] {
						return Positioned(
							contents,
							modifier.Position(p),
							modifier.Align(a),
							modifier.Visible(visible.Value()),
						)
					}))
			})),
		).Children(
			Box(
				modifier.Align(alignment.Center),
			).Children(
				Text("Hover me"),
			),
		),
	)
}
