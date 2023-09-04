package compose

import (
	"github.com/octohelm/gio-compose/pkg/layout"

	"github.com/octohelm/gio-compose/pkg/compose/internal"
)

type VNode = internal.VNode
type BuildContext = internal.BuildContext
type Component = internal.Component
type ComponentWrapper = internal.ComponentWrapper
type Element = internal.Element
type ElementPainter = internal.ElementPainter

func ElementPainterFunc(layout func(gtx layout.Context) layout.Dimensions) ElementPainter {
	return &elementPainterFunc{
		layout: layout,
	}
}

type elementPainterFunc struct {
	layout func(gtx layout.Context) layout.Dimensions
}

func (w *elementPainterFunc) Layout(gtx layout.Context) layout.Dimensions {
	return w.layout(gtx)
}
