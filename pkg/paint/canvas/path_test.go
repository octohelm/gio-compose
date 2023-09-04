package canvas

import (
	"image"
	"image/color"
	"testing"

	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"

	"github.com/octohelm/gio-compose/internal/testutil"
)

func TestCanvas(t *testing.T) {
	p := &Path{}

	p.MoveTo(300, 300)
	p.LineTo(270.7107, 220.71068)
	p.ArcTo(100, 50, 45, false, true, 217.83092, 222.31444)
	p.Close()

	img, _ := testutil.DrawImage(t, image.Pt(400, 400), func(ops *op.Ops) {
		func() {
			// draw fill
			defer clip.Outline{Path: PathSpec(ops, p)}.Op().Push(ops).Pop()
			paint.Fill(ops, color.NRGBA{R: 0xff, A: 0xff})
		}()

		func() {
			stroke := p.Stroke(
				8,
				RoundCap,
				BevelJoin,
				10,
			)
			// draw stroke
			defer clip.Outline{Path: PathSpec(ops, stroke)}.Op().Push(ops).Pop()
			paint.Fill(ops, color.NRGBA{A: 0xff})
		}()

	})

	_ = testutil.WritePNG("results/demo.png", img)
}
