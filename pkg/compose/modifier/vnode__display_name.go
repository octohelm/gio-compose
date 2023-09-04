package modifier

import (
	"github.com/octohelm/gio-compose/pkg/compose/internal"
	"github.com/octohelm/gio-compose/pkg/modifier"
	"github.com/octohelm/gio-compose/pkg/node"
)

var _ modifier.Modifier[any] = DisplayName("")

type DisplayName string

func (a DisplayName) VNodeOnly() {
}

func (a DisplayName) Modify(w any) {
	switch x := w.(type) {
	case node.DisplayNameSetter:
		x.SetDisplayName(string(a))
	case internal.VNodeAccessor:
		a.Modify(x.Type())
	}
}
