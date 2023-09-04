package modifier

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/octohelm/gio-compose/pkg/modifier"
	"github.com/octohelm/gio-compose/pkg/paint"
	"github.com/octohelm/gio-compose/pkg/paint/size"
	"github.com/octohelm/gio-compose/pkg/unit"
)

func Height(dp unit.Dp) modifier.Modifier[any] {
	return &sizeModifier{v: dp, sizes: []size.Type{
		size.Height,
	}}
}

func Width(dp unit.Dp) modifier.Modifier[any] {
	return &sizeModifier{v: dp, sizes: []size.Type{
		size.Width,
	}}
}

func Size(dp unit.Dp) modifier.Modifier[any] {
	return &sizeModifier{v: dp, sizes: []size.Type{
		size.Width,
		size.Height,
	}}
}

type sizeModifier struct {
	v     unit.Dp
	sizes []size.Type
}

func (m *sizeModifier) String() string {
	return fmt.Sprintf("%s = %v", m.sizes, m.v)
}

func FillMaxHeight(fractions ...float32) modifier.Modifier[any] {
	return &sizeModifier{v: -getFraction(fractions), sizes: []size.Type{
		size.Height,
	}}
}

func FillMaxWidth(fractions ...float32) modifier.Modifier[any] {
	return &sizeModifier{v: -getFraction(fractions), sizes: []size.Type{
		size.Width,
	}}
}

func FillMaxSize(fractions ...float32) modifier.Modifier[any] {
	return &sizeModifier{v: -getFraction(fractions), sizes: []size.Type{
		size.Width,
		size.Height,
	}}
}

func (m *sizeModifier) Modify(w any) {
	if setter, ok := w.(paint.SizeSetter); ok {
		setter.SetSize(m.v, m.sizes...)
	}
}

func getFraction(fractions []float32) unit.Dp {
	fraction := float32(1.0)
	if len(fractions) > 0 {
		if fractions[0] > 0.0 && fractions[0] <= 1.0 {
			fraction = fractions[0]
		} else {
			panic(errors.Errorf("invalid fraction, expect (0,1], but got %v", fractions[0]))
		}
	}
	return unit.Dp(fraction)
}
