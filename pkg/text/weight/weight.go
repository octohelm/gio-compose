package weight

import "gioui.org/font"

type Weight int

const (
	Thin       Weight = 100
	ExtraLight Weight = 200
	Light      Weight = 300
	Normal     Weight = 400
	Medium     Weight = 500
	SemiBold   Weight = 600
	Bold       Weight = 700
	ExtraBold  Weight = 800
	Black      Weight = 900
)

func (w Weight) GioWeight() font.Weight {
	return font.Weight(w - 400)
}
