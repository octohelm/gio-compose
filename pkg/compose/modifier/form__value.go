package modifier

import (
	"github.com/octohelm/gio-compose/pkg/compose"
	"github.com/octohelm/gio-compose/pkg/modifier"
)

func Value[T any](v T, stringify func(v T) string) modifier.Modifier[any] {
	return &valueModifier{
		txt: stringify(v),
	}
}

type valueModifier struct {
	txt string
}

func (v *valueModifier) Modify(target any) {
	if setter, ok := target.(compose.TextSetter); ok {
		setter.SetText(v.txt)
	}
}
