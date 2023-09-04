package paint

import (
	"image/color"

	"github.com/octohelm/gio-compose/pkg/cmp"

	"gioui.org/op"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

type FillSetter interface {
	SetFill(c color.Color)
}

type Fill struct {
	Color color.NRGBA
}

func (f *Fill) Eq(v *Fill) cmp.Result {
	return func() bool {
		return v.Color == f.Color
	}
}

func (f *Fill) Transparent() bool {
	return f.Color.A == 0
}

var _ FillSetter = &Fill{}

func (f *Fill) SetFill(c color.Color) {
	if c != nil {
		f.Color = color.NRGBAModel.Convert(c).(color.NRGBA)
	}
}

func (f *Fill) PaintClip(gtx layout.Context, clip Clip) {
	if !f.Transparent() {
		paint.FillShape(gtx.Ops, f.Color, clip.Op(gtx.Ops))
	}
}

func (f *Fill) PaintOps(gtx layout.Context) {
	if !f.Transparent() {
		paint.ColorOp{Color: f.Color}.Add(gtx.Ops)
	}
}

type Clip interface {
	Op(ops *op.Ops) clip.Op
}

func ClipOutlineOf(path clip.PathSpec) Clip {
	return ClipFunc(func(ops *op.Ops) clip.Op {
		return clip.Outline{Path: path}.Op()
	})
}

func ClipFunc(fn func(ops *op.Ops) clip.Op) Clip {
	return &clipFunc{
		fn: fn,
	}
}

type clipFunc struct {
	fn func(ops *op.Ops) clip.Op
}

func (c *clipFunc) Op(ops *op.Ops) clip.Op {
	return c.fn(ops)
}
