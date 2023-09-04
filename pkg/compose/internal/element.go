package internal

import (
	"context"

	"github.com/octohelm/gio-compose/pkg/modifier"

	"github.com/octohelm/gio-compose/pkg/layout"
	"github.com/octohelm/gio-compose/pkg/node"
)

type Element interface {
	node.Node

	ElementCreator
	ElementPatcher
	ElementPainter
}

type ElementCreator interface {
	New(ctx context.Context) Element
}

type ElementPatcher interface {
	Update(ctx context.Context, modifiers ...modifier.Modifier[any]) bool
}

type ElementPainter interface {
	Layout(gtx layout.Context) layout.Dimensions
}
