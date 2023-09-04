package component

import (
	. "github.com/octohelm/gio-compose/pkg/compose"
	"github.com/octohelm/gio-compose/pkg/compose/renderer"
	modifierapi "github.com/octohelm/gio-compose/pkg/modifier"
	"github.com/octohelm/gio-compose/pkg/node"
)

func RootPortal(modifiers ...modifierapi.Modifier[any]) VNode {
	return H(&rootPortal{}, modifiers...)
}

type rootPortal struct {
}

func (rootPortal) Build(c BuildContext) VNode {
	portalRoot := UseMemo(c, func() node.Node {
		return renderer.CreateRootWidget(c, "Portal")
	}, []any{})

	UseEffect(c, func() func() {
		w := renderer.WindowContext.Extract(c)
		node.AppendChild(w, portalRoot)
		return func() {
			node.RemoveChild(w, portalRoot)
		}
	}, []any{})

	return Portal(portalRoot).Children(c.ChildVNodes()...)
}
