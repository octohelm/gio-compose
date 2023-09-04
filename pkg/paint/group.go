package paint

import (
	"gioui.org/op"
)

func Group(ops *op.Ops, do func()) op.CallOp {
	macro := op.Record(ops)
	do()
	return macro.Stop()
}
