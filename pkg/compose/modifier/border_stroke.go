package modifier

import (
	"image/color"

	"github.com/octohelm/gio-compose/pkg/layout/position"
	"github.com/octohelm/gio-compose/pkg/modifier"
	"github.com/octohelm/gio-compose/pkg/paint"
	"github.com/octohelm/gio-compose/pkg/unit"
)

func BorderStrokeAll(width unit.Dp, c color.Color) modifier.Modifier[any] {
	ss := paint.StrokeStyle{
		Width: width,
		Color: color.NRGBAModel.Convert(c).(color.NRGBA),
	}

	return &borderStrokeModifier{
		strokeStyles: []paint.StrokeStyleWithPosition{
			{Position: position.Top, StrokeStyle: ss},
			{Position: position.Right, StrokeStyle: ss},
			{Position: position.Bottom, StrokeStyle: ss},
			{Position: position.Left, StrokeStyle: ss},
		},
	}
}

func BorderStrokeTop(width unit.Dp, c color.Color) modifier.Modifier[any] {
	ss := paint.StrokeStyle{
		Width: width,
		Color: color.NRGBAModel.Convert(c).(color.NRGBA),
	}

	return &borderStrokeModifier{
		strokeStyles: []paint.StrokeStyleWithPosition{
			{Position: position.Top, StrokeStyle: ss},
		},
	}
}

func BorderStrokeRight(width unit.Dp, c color.Color) modifier.Modifier[any] {
	ss := paint.StrokeStyle{
		Width: width,
		Color: color.NRGBAModel.Convert(c).(color.NRGBA),
	}

	return &borderStrokeModifier{
		strokeStyles: []paint.StrokeStyleWithPosition{
			{Position: position.Right, StrokeStyle: ss},
		},
	}
}

func BorderStrokeBottom(width unit.Dp, c color.Color) modifier.Modifier[any] {
	ss := paint.StrokeStyle{
		Width: width,
		Color: color.NRGBAModel.Convert(c).(color.NRGBA),
	}

	return &borderStrokeModifier{
		strokeStyles: []paint.StrokeStyleWithPosition{
			{Position: position.Bottom, StrokeStyle: ss},
		},
	}
}

func BorderStrokeLeft(width unit.Dp, c color.Color) modifier.Modifier[any] {
	ss := paint.StrokeStyle{
		Width: width,
		Color: color.NRGBAModel.Convert(c).(color.NRGBA),
	}

	return &borderStrokeModifier{
		strokeStyles: []paint.StrokeStyleWithPosition{
			{Position: position.Left, StrokeStyle: ss},
		},
	}
}

type borderStrokeModifier struct {
	strokeStyles []paint.StrokeStyleWithPosition
}

func (f *borderStrokeModifier) Modify(w any) {
	if setter, ok := w.(paint.BorderStrokeSetter); ok {
		setter.SetBorderStroke(f.strokeStyles...)
	}
}
