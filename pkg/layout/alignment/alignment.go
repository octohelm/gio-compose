package alignment

import "gioui.org/layout"

type Alignment int

const (
	TopStart    Alignment = iota // layout.NW
	Top                          // layout.N
	TopEnd                       // layout.NE
	End                          // layout.E
	BottomStart                  // layout.SE
	Bottom                       // layout.S
	BottomEnd                    // layout.SW
	Start                        // layout.W
	Center                       // layout.Center

	Baseline // layout.Center
)

const Middle = Center

func (a Alignment) String() string {
	switch a {
	case TopStart:
		return "TopStart"
	case Top:
		return "Top"
	case TopEnd:
		return "TopEnd"
	case BottomStart:
		return "BottomStart"
	case Bottom:
		return "Bottom"
	case BottomEnd:
		return "BottomEnd"
	case Start:
		return "Start"
	case Center:
		return "Center"
	case End:
		return "End"
	case Baseline:
		return "Baseline"
	}
	return ""
}

func (a Alignment) LayoutDirection() layout.Direction {
	if a == Baseline {
		return layout.Center
	}
	return layout.Direction(a)
}

func (a Alignment) LayoutAlignment() layout.Alignment {
	switch a {
	case Middle:
		return layout.Middle
	case Baseline:
		return layout.Baseline
	case End:
		return layout.End
	}
	return layout.Start
}
