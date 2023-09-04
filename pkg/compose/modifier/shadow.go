package modifier

import (
	"github.com/octohelm/gio-compose/pkg/modifier"
	"github.com/octohelm/gio-compose/pkg/paint"
	"github.com/octohelm/gio-compose/pkg/unit"
)

func Shadow(elevation unit.Dp) modifier.Modifier[any] {
	return &shadowModifier{
		elevation: elevation,
	}
}

type shadowModifier struct {
	elevation unit.Dp
}

func (m *shadowModifier) Modify(w any) {
	if setter, ok := w.(paint.ShadowSetter); ok {
		setter.SetShadow(m.elevation)
	}
}
