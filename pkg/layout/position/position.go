package position

type Position int

const (
	Top Position = iota + 1
	Right
	Bottom
	Left
	TopLeft
	BottomLeft
	TopRight
	BottomRight
)

func (a Position) String() string {
	switch a {
	case Top:
		return "Top"
	case Right:
		return "Right"
	case Bottom:
		return "Bottom"
	case Left:
		return "Left"
	case TopLeft:
		return "TopLeft"
	case BottomLeft:
		return "BottomLeft"
	case TopRight:
		return "TopRight"
	case BottomRight:
		return "BottomRight"
	}
	return ""
}

var AllFour = []Position{
	Top,
	Right,
	Bottom,
	Left,
}
