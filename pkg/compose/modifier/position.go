package modifier

import (
	"github.com/octohelm/gio-compose/pkg/layout/position"

	"github.com/octohelm/gio-compose/pkg/layout"
	"github.com/octohelm/gio-compose/pkg/modifier"
)

func Position(pos position.Position) modifier.Modifier[any] {
	return &positionModifier{
		position: pos,
	}
}

type positionModifier struct {
	position position.Position
}

func (m *positionModifier) Modify(w any) {
	if setter, ok := w.(layout.PositionSetter); ok {
		setter.SetPosition(m.position)
	}
}
