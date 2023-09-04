package internal

import "context"

type ContextProvider interface {
	GetChildContext(ctx context.Context) context.Context
}

type Provider func(ctx context.Context) context.Context

func (p Provider) GetChildContext(ctx context.Context) context.Context {
	return p(ctx)
}

func (Provider) Build(c BuildContext) VNode {
	return H(Fragment{}).Children(c.ChildVNodes()...)
}
