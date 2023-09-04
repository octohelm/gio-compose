package component

import (
	. "github.com/octohelm/gio-compose/pkg/compose"
	"github.com/octohelm/gio-compose/pkg/compose/modifier"
	"github.com/octohelm/gio-compose/pkg/event"
	"github.com/octohelm/gio-compose/pkg/layout"
	"github.com/octohelm/gio-compose/pkg/layout/alignment"
	"github.com/octohelm/gio-compose/pkg/layout/position"
	"github.com/octohelm/gio-compose/pkg/layout/visible"
	modifierapi "github.com/octohelm/gio-compose/pkg/modifier"
	"github.com/octohelm/gio-compose/pkg/paint"
	"github.com/octohelm/gio-compose/pkg/unit"
)

func Positioned(contents VNode, modifiers ...modifierapi.Modifier[any]) interface {
	modifierapi.Modifier[any]
	ComponentWrapper
} {
	return &positioned{
		contents:  contents,
		modifiers: modifiers,
	}
}

var _ ComponentWrapper = &positioned{}

type positioned struct {
	modifierapi.Discord

	target    VNode
	contents  VNode
	modifiers modifierapi.Modifiers
}

func (p positioned) Wrap(node VNode) Component {
	p.target = node
	return &p
}

func (p *positioned) Build(c BuildContext) VNode {
	if p.target == nil || p.contents == nil {
		return nil
	}

	triggerWidget := UseRef[Element](c, nil)
	contentWidget := UseRef[Element](c, nil)

	a := &positionState{}
	a.Position = position.Bottom
	a.Alignment = alignment.Center

	modifierapi.Modify[any](a, p.modifiers)

	positionUpdate := func() {
		modifier.OffsetXY(
			a.Calc(triggerWidget.Current, contentWidget.Current),
		).Modify(contentWidget.Current)
	}

	return Fragment(
		CloneNode(
			p.target,
			modifier.Ref(func(w Element) {
				triggerWidget.Current = w
			}),
		).Children(c.ChildVNodes()...),
		H(Overlay{
			Visible: a.Visible,
			OnVisibleChange: func(v bool) {
				a.VisibleEvents().Trigger(visible.Change, &visible.EventData{
					Visible: v,
				})
			},
			Refs: []*Ref[Element]{
				triggerWidget,
				contentWidget,
			},
		}).Children(
			Box(
				modifier.Ref(func(w Element) {
					contentWidget.Current = w
				}),
				modifier.DetectLayout(
					layout.OnDidSize(positionUpdate),
				),
			).Children(
				CloneNode(p.contents),
			),
		),
	)
}

type positionState struct {
	paint.Visibility
	layout.Aligner
	layout.Positioner

	events event.Events[visible.EventType]
}

var _ visible.EventsAccessor = &positionState{}

func (a *positionState) VisibleEvents() *event.Events[visible.EventType] {
	return &a.events
}

func (a *positionState) Calc(trigger Element, content Element) (x unit.Dp, y unit.Dp) {
	triggerRect := layout.GetBoundingClientRect(trigger)
	rect := layout.GetBoundingClientRect(content)

	x = triggerRect.X
	y = triggerRect.Y

	vertical := false

	switch a.Position {
	case position.Top:
		y -= rect.Height
		vertical = true
	case position.Right:
		x += triggerRect.Width
	case position.Left:
		x -= rect.Width
	case position.Bottom:
		y += triggerRect.Height
		vertical = true
	}

	switch a.Alignment {
	case alignment.Center:
		if !vertical {
			y += triggerRect.Height/2 - rect.Height/2
		} else {
			x += triggerRect.Width/2 - rect.Width/2
		}
	case alignment.End:
		if !vertical {
			y += triggerRect.Height - rect.Height
		} else {
			x += triggerRect.Width - rect.Width
		}
	}

	return
}
