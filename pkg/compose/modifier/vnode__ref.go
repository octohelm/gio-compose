package modifier

import (
	"context"

	"github.com/octohelm/gio-compose/pkg/compose"
	"github.com/octohelm/gio-compose/pkg/compose/internal"
	"github.com/octohelm/gio-compose/pkg/modifier"
	"github.com/octohelm/gio-compose/pkg/node"
)

func emitRef(n node.Node, ref func(w compose.Element)) {
	if ref == nil {
		return
	}
	switch x := n.(type) {
	case internal.Element:
		ref(x)
	case *node.Fragment:
		// FIXME may support forwardRef
		emitRef(n.FirstChild(), ref)
	default:
	}
}

func Ref(ref func(w compose.Element)) modifier.Modifier[any] {
	return &refModifier{
		onMount: func(ctx context.Context, w internal.VNodeAccessor) {
			emitRef(w.Node(), ref)
		}}
}

type refModifier struct {
	onMount func(ctx context.Context, w internal.VNodeAccessor)
}

func (k *refModifier) VNodeOnly() {
}

func (k *refModifier) Modify(t any) {
	if s, ok := t.(internal.VNodeAccessor); ok {
		s.OnMount(k.onMount)
	}
}
