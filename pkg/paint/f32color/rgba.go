package f32color

import (
	"image/color"
	"math"

	_ "unsafe"
)

//go:linkname  MulAlpha gioui.org/internal/f32color.MulAlpha
func MulAlpha(c color.NRGBA, alpha uint8) color.NRGBA

//go:linkname Disabled gioui.org/internal/f32color.Disabled
func Disabled(c color.NRGBA) (d color.NRGBA)

//go:linkname Hovered gioui.org/internal/f32color.Hovered
func Hovered(c color.NRGBA) (d color.NRGBA)

//go:linkname NRGBAToLinearRGBA gioui.org/internal/f32color.NRGBAToLinearRGBA
func NRGBAToLinearRGBA(c color.NRGBA) color.RGBA

func RGB(c uint32) color.NRGBA {
	return ARGB((0xff << 24) | c)
}

func ARGB(c uint32) color.NRGBA {
	return color.NRGBA{A: uint8(c >> 24), R: uint8(c >> 16), G: uint8(c >> 8), B: uint8(c)}
}

func Alpha(c color.NRGBA, alpha float32) color.NRGBA {
	c.A = uint8(math.Round(float64(alpha*256) + 0.5))
	return c
}

func IsTransparent(c color.NRGBA) bool {
	return c.A == 0
}
