package canvas

import (
	"math"
)

const (
	pi2 = 2 * math.Pi
	deg = math.Pi / 180
)

func ellipseFromArc(startX, startY float64, rx, ry float64, xAxisRotation float64, largeArcFlag, sweepFlag int, endX, endY float64) *ellipse {
	largeArc := largeArcFlag != 0
	if sweepFlag != 0 {
		largeArc = !largeArc
	}

	rot := xAxisRotation * deg

	sin, cos := math.Sincos(rot)

	// Move origin to start point
	nx, ny := endX-startX, endY-startY

	// rot ellipse x-axis to coordinate x-axis
	nx, ny = nx*cos+ny*sin, -nx*sin+ny*cos

	// Scale X dimension so that rx = ry
	// Now the ellipse is a circle radius ry; therefore foci and center coincide
	nx *= ry / rx

	midX, midY := nx/2, ny/2
	midlenSq := midX*midX + midY*midY

	var hr float64
	if ry*ry < midlenSq {
		// Requested ellipse does not exist; scale rx, ry to fit. Length of
		// span is greater than max width of ellipse, must scale *rx, *ry
		nrb := math.Sqrt(midlenSq)
		if rx == ry {
			rx = nrb // prevents roundoff
		} else {
			rx = rx * nrb / ry
		}
		ry = nrb
	} else {
		hr = math.Sqrt(ry*ry-midlenSq) / math.Sqrt(midlenSq)
	}

	var cx, cy float64

	if largeArc {
		cx = midX - midY*hr
		cy = midY + midX*hr
	} else {
		cx = midX + midY*hr
		cy = midY - midX*hr
	}

	// Reverse scale
	cx *= rx / ry

	e := &ellipse{
		rotate:   xAxisRotation * deg,
		rx:       rx,
		ry:       ry,
		largeArc: largeArc,
	}

	// reverse rot and translate back to original coordinates
	e.cx = cx*cos - cy*sin + (startX)
	e.cy = cx*sin + cy*cos + (startY)

	return e
}

type ellipse struct {
	cx, cy   float64 // translate
	rx, ry   float64 // scale
	rotate   float64 // rotate
	largeArc bool
}

func (eps *ellipse) Foci() (Point, Point) {
	// sqrt(a^2 - b^2)
	c := math.Sqrt(eps.rx*eps.rx - eps.ry*eps.ry)
	sin, cos := math.Sincos(eps.rotate)
	d := Pt(c*(cos), c*(sin))

	center := Pt(eps.cx, eps.cy)
	return center.Add(d), center.Sub(d)
}

func (eps *ellipse) PointAt(angle float64) Point {
	sin, cos := math.Sincos(angle)
	pt := Pt(eps.rx*(cos), eps.ry*(sin))

	if eps.rotate != 0 {
		pt = rotate(pt, eps.rotate)
	}

	return pt.Add(Pt(eps.cx, eps.cy))
}

func (eps *ellipse) AngleOf(p Point) float64 {
	// de translate
	p = p.Sub(Pt(eps.cx, eps.cy))

	// de rotate
	if eps.rotate != 0 {
		p = rotate(p, -eps.rotate)
	}

	angle := math.Atan2(
		p.Y/eps.ry,
		p.X/eps.rx,
	)

	if angle < 0 {
		angle += pi2
	}

	if angle > pi2 {
		angle -= pi2
	}

	return angle
}

func rotate(p Point, radian float64) Point {
	sin, cos := math.Sincos(radian)

	x := p.X*(cos) - p.Y*(sin)
	y := p.X*(sin) + p.Y*(cos)

	p.X = x
	p.Y = y

	return p
}

func boolInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
