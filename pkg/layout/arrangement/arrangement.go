package arrangement

import "gioui.org/layout"

type Arrangement int

const (
	Start Arrangement = iota
	Center
	End
	SpaceBetween
	SpaceAround
	SpaceEvenly
	EqualWeight
)

func (a Arrangement) LayoutSpacing() layout.Spacing {
	switch a {
	case Start:
		return layout.SpaceEnd
	case End:
		return layout.SpaceStart
	case Center:
		return layout.SpaceSides
	case SpaceAround:
		return layout.SpaceAround
	case SpaceBetween:
		return layout.SpaceBetween
	case SpaceEvenly:
		return layout.SpaceEvenly
	}
	return layout.SpaceSides
}
