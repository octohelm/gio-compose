package modifier

import (
	"github.com/octohelm/gio-compose/pkg/layout"
	"github.com/octohelm/gio-compose/pkg/layout/direction"
	"github.com/octohelm/gio-compose/pkg/modifier"
	"github.com/octohelm/gio-compose/pkg/unit"
)

func Offset(dp unit.Dp) modifier.Modifier[any] {
	return modifier.Modifiers{
		OffsetX(dp),
		OffsetY(dp),
	}
}

func OffsetXY(x unit.Dp, y unit.Dp) modifier.Modifier[any] {
	return modifier.Modifiers{
		OffsetX(x),
		OffsetY(y),
	}
}

func OffsetX(dp unit.Dp) modifier.Modifier[any] {
	return &offsetModifier{
		direction: direction.X,
		dp:        dp,
	}
}

func OffsetY(dp unit.Dp) modifier.Modifier[any] {
	return &offsetModifier{
		direction: direction.Y,
		dp:        dp,
	}
}

type offsetModifier struct {
	direction direction.Direction
	dp        unit.Dp
}

func (e *offsetModifier) Modify(w any) {
	if setter, ok := w.(layout.OffsetSetter); ok {
		setter.SetOffset(e.dp, e.direction)
	}
}
