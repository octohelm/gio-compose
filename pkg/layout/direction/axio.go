package direction

import "gioui.org/layout"

type Axis int

const (
	Horizontal Axis = iota
	Vertical
)

func (a Axis) LayoutAxis() layout.Axis {
	return layout.Axis(a)
}
