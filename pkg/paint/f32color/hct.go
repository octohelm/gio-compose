package f32color

import (
	"image/color"

	"goki.dev/cam/hct"
)

func HCTFromColor(c color.Color) HCT {
	cc := hct.FromColor(c)
	return HCT{Hue: cc.Hue, Chroma: cc.Chroma, Tone: cc.Tone}
}

func Toned(c color.Color, tone float32) HCT {
	cc := HCTFromColor(c)
	cc.Tone = tone
	return cc
}

type HCT struct {
	Hue, Chroma, Tone float32
}

func (h HCT) RGBA() (uint32, uint32, uint32, uint32) {
	r, g, b := hct.SolveToRGB(h.Hue, h.Chroma, h.Tone)
	return uint32(r*65535.0 + 0.5), uint32(g*65535.0 + 0.5), uint32(b*65535.0 + 0.5), uint32(65535)
}
