package renderer

import (
	"context"

	"github.com/octohelm/gio-compose/pkg/util/contextutil"

	"github.com/octohelm/gio-compose/pkg/compose"
	"github.com/octohelm/gio-compose/pkg/compose/internal"
	"github.com/octohelm/gio-compose/pkg/compose/modifier"
	"github.com/octohelm/gio-compose/pkg/node"
)

var WindowContext = contextutil.New[node.Node]()

func CreateRootWidget(ctx context.Context, name string) node.Node {
	rootVNode := compose.Box(
		modifier.DisplayName(name),
		modifier.FillMaxSize(),
	).(internal.VNodeAccessor)

	if widget, ok := rootVNode.Type().(compose.Element); ok {
		n := widget.New(ctx)
		_ = n.Update(ctx, rootVNode.Modifiers()...)
		rootVNode.BindNode(n)
	}

	return rootVNode.Node()
}
