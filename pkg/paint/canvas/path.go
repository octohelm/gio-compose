package canvas

import (
	"gioui.org/op"
	"gioui.org/op/clip"
	"github.com/octohelm/gio-compose/pkg/f32"
	"github.com/tdewolff/canvas"
)

var (
	RoundCap  = canvas.RoundCap
	ButtCap   = canvas.ButtCap
	SquareCap = canvas.SquareCap
)

var (
	RoundJoin = canvas.RoundJoin
	MiterJoin = canvas.MiterJoin
	ArcsJoin  = canvas.ArcsJoin
	BevelJoin = canvas.BevelJoin
)

type Path = canvas.Path

type Point = canvas.Point

func Pt(x, y float64) Point {
	return Point{
		X: x,
		Y: y,
	}
}

func ToF32(p canvas.Point) f32.Point {
	return f32.Point{
		X: float32(p.X),
		Y: float32(p.Y),
	}
}

func PathSpec(ops *op.Ops, path *Path) clip.PathSpec {
	p := &clip.Path{}
	p.Begin(ops)

	var copyTo func(p *clip.Path, path *Path, skipMove bool)

	copyTo = func(p *clip.Path, path *Path, skipMove bool) {
		for scanner := path.Scanner(); scanner.Scan(); {
			switch scanner.Cmd() {
			case canvas.CloseCmd:
				p.Close()
			case canvas.MoveToCmd:
				if skipMove {
					continue
				}
				p.MoveTo(ToF32(scanner.End()))
			case canvas.LineToCmd:
				p.LineTo(ToF32(scanner.End()))
			case canvas.QuadToCmd:
				p.QuadTo(ToF32(scanner.CP1()), ToF32(scanner.End()))
			case canvas.CubeToCmd:
				p.CubeTo(ToF32(scanner.CP1()), ToF32(scanner.CP2()), ToF32(scanner.End()))
			case canvas.ArcToCmd:
				// GIO use the same way to process arcs
				p2 := &Path{}

				start := scanner.Start()
				end := scanner.End()

				rx, ry, rot, large, sweep := scanner.Arc()
				p2.MoveTo(start.X, start.Y)
				p2.ArcTo(rx, ry, rot, large, sweep, end.X, end.Y)

				copyTo(p, p2.ReplaceArcs(), true)
			}
		}
	}

	copyTo(p, path, false)

	return p.End()

}
