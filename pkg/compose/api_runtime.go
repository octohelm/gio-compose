package compose

import (
	"context"

	"github.com/octohelm/gio-compose/pkg/modifier"

	"github.com/octohelm/gio-compose/pkg/node"

	"github.com/octohelm/gio-compose/pkg/compose/internal"
	"github.com/pkg/errors"
)

// H Create VNode from Component
func H(component Component, modifiers ...modifier.Modifier[any]) VNode {
	return internal.H(component, modifiers...)
}

// Fragment children without Element wrapper
func Fragment(vnodes ...VNode) VNode {
	return H(internal.Fragment{}).Children(vnodes...)
}

// Provider inject something into context for child nodes
func Provider(c func(ctx context.Context) context.Context, modifiers ...modifier.Modifier[any]) VNode {
	return H(internal.Provider(c), modifiers...)
}

func Portal(mountPoint node.Node) VNode {
	n := internal.H(internal.Portal{})
	if mountPoint == nil {
		panic(errors.New("mount point required for Portal"))
	}
	n.(internal.VNodeAccessor).BindNode(mountPoint)
	return n
}

func CloneNode(n VNode, modifiers ...modifier.Modifier[any]) VNode {
	if n == nil {
		return nil
	}
	na := n.(internal.VNodeAccessor)

	vn := H(na.Type(), append(na.Modifiers(), modifiers...)...)

	if children := na.ChildVNodes(); len(children) > 0 {
		return vn.Children(children...)
	}

	return vn
}
