package modifier

import (
	"github.com/octohelm/gio-compose/pkg/compose/internal"
	"github.com/octohelm/gio-compose/pkg/modifier"
)

func Key(key any) modifier.Modifier[any] {
	return &keyModifier{key: key}
}

type keyModifier struct {
	key any
}

func (k *keyModifier) VNodeOnly() {
}

func (k *keyModifier) Modify(t any) {
	if s, ok := t.(internal.VNodeAccessor); ok {
		s.SetKey(k.key)
	}
}
