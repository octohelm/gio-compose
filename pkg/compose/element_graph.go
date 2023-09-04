package compose

import (
	"context"

	"github.com/octohelm/gio-compose/pkg/modifier"

	"github.com/octohelm/gio-compose/pkg/compose/internal"
	"github.com/octohelm/gio-compose/pkg/layout"
	"github.com/octohelm/gio-compose/pkg/node"
)

func Graph(layout func(gtx layout.Context) layout.Dimensions) Element {
	return &graphElement{
		Element: node.Element{
			Name: "Graph",
		},
		layout: layout,
	}
}

var _ Element = &graphElement{}

type graphElement struct {
	internal.ElementComponent
	node.Element
	layout func(gtx layout.Context) layout.Dimensions
}

func (w *graphElement) Update(ctx context.Context, modifiers ...modifier.Modifier[any]) bool {
	return true
}

func (w *graphElement) New(ctx context.Context) internal.Element {
	return &graphElement{
		Element: node.Element{
			Name: w.Name,
		},
	}
}

func (w *graphElement) Layout(gtx layout.Context) layout.Dimensions {
	return w.layout(gtx)
}
