package modifier

import (
	"fmt"

	"github.com/octohelm/gio-compose/pkg/layout/position"
	"github.com/octohelm/gio-compose/pkg/modifier"
	"github.com/octohelm/gio-compose/pkg/paint"
	"github.com/octohelm/gio-compose/pkg/unit"
	"github.com/octohelm/x/slices"
)

func RoundedAll(dp unit.Dp) modifier.Modifier[any] {
	return &cornerRadiusModifier{
		dp: dp,
		positions: []position.Position{
			position.TopLeft,
			position.TopRight,
			position.BottomLeft,
			position.BottomRight,
		},
	}
}

func RoundedLeft(dp unit.Dp) modifier.Modifier[any] {
	return &cornerRadiusModifier{
		dp: dp,
		positions: []position.Position{
			position.TopLeft,
			position.BottomLeft,
		},
	}
}

func RoundedRight(dp unit.Dp) modifier.Modifier[any] {
	return &cornerRadiusModifier{
		dp: dp,
		positions: []position.Position{
			position.TopRight,
			position.BottomRight,
		},
	}
}

func RoundedBottom(dp unit.Dp) modifier.Modifier[any] {
	return &cornerRadiusModifier{
		dp: dp,
		positions: []position.Position{
			position.BottomLeft,
			position.BottomRight,
		},
	}
}

func RoundedTop(dp unit.Dp) modifier.Modifier[any] {
	return &cornerRadiusModifier{
		dp: dp,
		positions: []position.Position{
			position.TopLeft,
			position.TopRight,
		},
	}
}

func RoundedTopLeft(dp unit.Dp) modifier.Modifier[any] {
	return &cornerRadiusModifier{
		dp: dp,
		positions: []position.Position{
			position.TopLeft,
		},
	}
}

func RoundedBottomLeft(dp unit.Dp) modifier.Modifier[any] {
	return &cornerRadiusModifier{
		dp: dp,
		positions: []position.Position{
			position.BottomLeft,
		},
	}
}

func RoundedTopRight(dp unit.Dp) modifier.Modifier[any] {
	return &cornerRadiusModifier{
		dp: dp,
		positions: []position.Position{
			position.TopRight,
		},
	}
}

func RoundedBottomRight(dp unit.Dp) modifier.Modifier[any] {
	return &cornerRadiusModifier{
		dp: dp,
		positions: []position.Position{
			position.BottomRight,
		},
	}
}

type cornerRadiusModifier struct {
	dp        unit.Dp
	positions []position.Position
}

func (m *cornerRadiusModifier) String() string {
	return fmt.Sprintf("%s = %v", slices.Map(m.positions, func(e position.Position) string {
		return "CornerRadius" + e.String()
	}), m.dp)
}

func (e *cornerRadiusModifier) Modify(widget any) {
	if setter, ok := widget.(paint.CornerRadiusSetter); ok {
		setter.SetCornerRadius(e.dp, e.positions...)
	}
}
