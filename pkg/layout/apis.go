package layout

import (
	"gioui.org/layout"
	"gioui.org/op"
)

type Context = layout.Context
type Constraints = layout.Constraints
type Dimensions = layout.Dimensions

func Layout(ops *op.Ops, do func(ops *op.Ops)) {
	do(ops)
}

func PostLayout(ops *op.Ops, do func(ops *op.Ops)) {
	m := op.Record(ops)
	do(ops)
	op.Defer(ops, m.Stop())
}
