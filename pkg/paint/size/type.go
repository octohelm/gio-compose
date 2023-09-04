package size

type Type int

const (
	Width Type = iota
	Height
)

func (t Type) String() string {
	if t == Width {
		return "Width"
	}
	return "Height"
}
