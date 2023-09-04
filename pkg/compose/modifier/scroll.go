package modifier

import (
	"github.com/octohelm/gio-compose/pkg/layout"
	"github.com/octohelm/gio-compose/pkg/layout/direction"
	"github.com/octohelm/gio-compose/pkg/modifier"
)

func VerticalScroll() modifier.Modifier[any] {
	return &scrollModifier{
		axis: direction.Vertical,
	}
}

func HorizontalScroll() modifier.Modifier[any] {
	return &scrollModifier{
		axis: direction.Horizontal,
	}
}

type scrollModifier struct {
	axis direction.Axis
}

func (m *scrollModifier) Modify(w any) {
	if setter, ok := w.(layout.ScrollableSetter); ok {
		setter.SetScrollable(m.axis, true)
	}
}
