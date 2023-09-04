package canvas

import (
	"fmt"
	"io"
	"math"
	"testing"

	"github.com/octohelm/gio-compose/internal/testutil"
)

func TestEllipse(t *testing.T) {
	t.Run("ellipseFromArc", func(t *testing.T) {
		file, _ := testutil.OpenFile("results/ellipse2.svg")
		defer file.Close()

		sz := Pt(400, 400)

		writeSvgTo(file, sz, func(w io.Writer) {
			_, _ = fmt.Fprintf(w, `
<path d="M 270.7107 220.71068 A 100 50 45 0 1 217.83092 222.31444" fill="none" stroke="red" />
`)

			e := ellipseFromArc(
				270.7107, 220.71068,
				100, 50, 45,
				0,
				1,
				217.83092, 222.31444,
			)

			writeEllipseTo(w,
				e.cx, e.cy, e.rx, e.ry,
				"rgba(255,0,0,0.15)",
				fmt.Sprintf("rotate(%v %v %v)", e.rotate/deg, e.cx, e.cy),
			)

			_, _ = fmt.Fprintf(w, `
<path d="M 270.7107 220.71068 A 100 50 45 1 1 217.83092 222.31444" fill="none" stroke="blue" />
`)

			e2 := ellipseFromArc(
				270.7107, 220.71068,
				100, 50, 45,
				1,
				1,
				217.83092, 222.31444,
			)

			writeEllipseTo(w,
				e2.cx, e2.cy, e2.rx, e2.ry,
				"rgba(0,0,255,0.15)",
				fmt.Sprintf("rotate(%v %v %v)", e2.rotate/deg, e2.cx, e2.cy),
			)
		})
	})

	t.Run("results", func(t *testing.T) {
		e := ellipse{
			cx:     200,
			cy:     150,
			rx:     100,
			ry:     50,
			rotate: 45 * deg,
		}

		file, _ := testutil.OpenFile("results/ellipse.svg")
		defer file.Close()

		sz := Pt(400, 300)

		writeSvgTo(file, sz, func(w io.Writer) {
			fmt.Fprintf(w, `
<path d="M 270.7107 220.71068 A 100 50 45 0 1 217.83092 222.31444" fill="none" stroke="red" />
`)

			f1, f2 := e.Foci()

			writeEllipseTo(w,
				e.cx, e.cy, e.rx, e.ry,
				"rgba(0,0,0,0.15)",
				fmt.Sprintf("rotate(%v %v %v)", e.rotate/deg, e.cx, e.cy),
			)

			writeCircle(w, f1.X, f1.Y, 2, "green")
			writeCircle(w, f2.X, f2.Y, 2, "orange")

			writeCircle(w, e.cx, e.cy, 2, "blue")

			for i := 0; i < 50; i++ {
				angle := float64(i) * math.Pi * 2 / 50.0

				p := e.PointAt(angle)

				writeCircle(w, p.X, p.Y, 2, "red")

				p2 := e.PointAt(e.AngleOf(p))
				writeCircle(w, p2.X, p2.Y, 2, "cyan")
			}
		})
	})
}

func writeSvgTo(w io.Writer, sz Point, do func(p io.Writer)) {
	_, _ = fmt.Fprintf(w, `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 %v %v">`, sz.X, sz.Y)
	do(w)
	_, _ = fmt.Fprintf(w, `</svg>`)
}

func writeEllipseTo(w io.Writer, cx, cy, rx, ry float64, fill string, transform string) {
	_, _ = fmt.Fprintf(w, `<ellipse cx="%v" cy="%v" rx="%v" ry="%v" fill=%q transform=%q/>
`, cx, cy, rx, ry, fill, transform)
}

func writeCircle(w io.Writer, cx, cy, r float64, fill string) {
	_, _ = fmt.Fprintf(w, `<circle cx="%v" cy="%v" r="%v"  fill=%q/>
`, cx, cy, r, fill)
}
