package modifier

import (
	"fmt"
	"image/color"

	"github.com/octohelm/gio-compose/pkg/modifier"
	"github.com/octohelm/gio-compose/pkg/paint"
	"github.com/octohelm/gio-compose/pkg/paint/f32color"
)

func BackgroundColor(c color.Color) modifier.Modifier[any] {
	return &backgroundColorModifier{c: c}
}

type backgroundColorModifier struct {
	c color.Color
}

func (m *backgroundColorModifier) String() string {
	return fmt.Sprintf("[BackgroudColor] = %s", f32color.HexString(m.c))
}

func (f *backgroundColorModifier) Modify(w any) {
	if setter, ok := w.(paint.FillSetter); ok {
		setter.SetFill(f.c)
	}
}
