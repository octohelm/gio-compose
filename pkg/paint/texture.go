package paint

import (
	"image/color"

	"gioui.org/op"
	"gioui.org/op/paint"
)

func ColorTexture(ops *op.Ops, c color.NRGBA) op.CallOp {
	macro := op.Record(ops)
	paint.ColorOp{Color: c}.Add(ops)
	return macro.Stop()
}
