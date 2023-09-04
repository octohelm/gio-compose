package f32

import (
	"math"

	"gioui.org/f32"
)

type Point = f32.Point

func Pt(x, y float32) Point {
	return Point{X: x, Y: y}
}

func Rotate(p f32.Point, radian float32) f32.Point {
	sin, cos := math.Sincos(float64(radian))

	x := p.X*float32(cos) - p.Y*float32(sin)
	y := p.X*float32(sin) + p.Y*float32(cos)

	p.X = x
	p.Y = y

	return p
}
