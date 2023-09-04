package unit

import (
	"gioui.org/unit"
)

type Metric = unit.Metric
type Dp = unit.Dp
type Sp = unit.Sp

func Pt(x, y Dp) Point {
	return Point{X: x, Y: y}
}

func PtFromPx(metric Metric, x, y int) Point {
	return Point{
		X: metric.PxToDp(x),
		Y: metric.PxToDp(y),
	}
}

type Point struct {
	X Dp
	Y Dp
}
