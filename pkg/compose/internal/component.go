package internal

import (
	"context"

	"github.com/octohelm/gio-compose/pkg/modifier"
)

type Component interface {
	Build(BuildContext) VNode
}

type BuildContext interface {
	context.Context

	ChildVNodes() []VNode
	Modifiers() modifier.Modifiers
}

type BuildContextAccessor interface {
	VNode() VNodeAccessor
	RawContext() context.Context
}

type ElementComponent struct{}

func (ElementComponent) Build(ctx BuildContext) VNode {
	return nil
}

type ComponentWrapper interface {
	Component
	Wrap(v VNode) Component
}
