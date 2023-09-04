package modifier

import (
	"fmt"

	"github.com/octohelm/gio-compose/pkg/modifier"

	"github.com/octohelm/gio-compose/pkg/layout"
	"github.com/octohelm/gio-compose/pkg/layout/position"
	"github.com/octohelm/gio-compose/pkg/unit"
	"github.com/octohelm/x/slices"
)

func PaddingAll(dp unit.Dp) modifier.Modifier[any] {
	return &edgeInsetModifier{
		dp: dp,
		positions: []position.Position{
			position.Left, position.Right, position.Top, position.Bottom,
		},
	}
}

func PaddingLeft(dp unit.Dp) modifier.Modifier[any] {
	return &edgeInsetModifier{
		dp: dp,
		positions: []position.Position{
			position.Left,
		},
	}
}

func PaddingRight(dp unit.Dp) modifier.Modifier[any] {
	return &edgeInsetModifier{
		dp: dp,
		positions: []position.Position{
			position.Right,
		},
	}
}

func PaddingBottom(dp unit.Dp) modifier.Modifier[any] {
	return &edgeInsetModifier{
		dp: dp,
		positions: []position.Position{
			position.Bottom,
		},
	}
}

func PaddingTop(dp unit.Dp) modifier.Modifier[any] {
	return &edgeInsetModifier{
		dp: dp,
		positions: []position.Position{
			position.Top,
		},
	}
}

func PaddingVertical(dp unit.Dp) modifier.Modifier[any] {
	return &edgeInsetModifier{
		dp: dp,
		positions: []position.Position{
			position.Top,
			position.Bottom,
		},
	}
}

func PaddingHorizontal(dp unit.Dp) modifier.Modifier[any] {
	return &edgeInsetModifier{
		dp: dp,
		positions: []position.Position{
			position.Left,
			position.Right,
		},
	}
}

type edgeInsetModifier struct {
	dp        unit.Dp
	positions []position.Position
}

func (m *edgeInsetModifier) String() string {
	return fmt.Sprintf("%s = %v", slices.Map(m.positions, func(e position.Position) string {
		return "Padding" + e.String()
	}), m.dp)
}

func (e *edgeInsetModifier) Modify(w any) {
	if setter, ok := w.(layout.EdgeInsetSetter); ok {
		setter.SetEdgeInset(e.dp, e.positions...)
	}
}
