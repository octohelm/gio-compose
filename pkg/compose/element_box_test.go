package compose_test

import (
	"bytes"
	"context"
	"image"
	"strings"
	"testing"

	"gioui.org/app"
	"gioui.org/io/system"
	giolayout "gioui.org/layout"
	"gioui.org/op"
	"github.com/octohelm/gio-compose/pkg/compose/modifier"
	"github.com/octohelm/gio-compose/pkg/layout"
	"github.com/octohelm/gio-compose/pkg/layout/alignment"
	testingx "github.com/octohelm/x/testing"

	. "github.com/octohelm/gio-compose/pkg/compose"
	"github.com/octohelm/gio-compose/pkg/compose/renderer"
)

func TestBox(t *testing.T) {
	t.Run("Box", func(t *testing.T) {
		t.Run("Centered Box", func(t *testing.T) {
			el := Box(modifier.FillMaxSize()).Children(
				Box(modifier.Align(alignment.Center), modifier.Size(200)),
			)

			ExpectLayout(t, el, image.Pt(1000, 1000), `
Box x=0,y=0,w=1000,h=1000
  Box x=400,y=400,w=200,h=200
`)
		})

		t.Run("Centered Padding Box", func(t *testing.T) {
			el := Box(modifier.FillMaxSize()).Children(
				Box(modifier.Align(alignment.Center), modifier.PaddingAll(50)).Children(
					Box(modifier.Size(100)),
					Box(modifier.FillMaxSize()),
				),
			)

			ExpectLayout(t, el, image.Pt(1000, 1000), `
Box x=0,y=0,w=1000,h=1000
  Box x=400,y=400,w=200,h=200
    Box x=50,y=50,w=100,h=100
    Box x=50,y=50,w=100,h=100
`)
		})
	})

	t.Run("Box in Scrollable", func(t *testing.T) {
		t.Run("Centered Box", func(t *testing.T) {
			el := Column(
				modifier.FillMaxSize(),
				modifier.VerticalScroll(),
			).Children(
				Box(modifier.FillMaxWidth()).Children(
					Box(modifier.Align(alignment.Center), modifier.Size(200)),
				),
			)

			ExpectLayout(t, el, image.Pt(1000, 1000), `
Column x=0,y=0,w=1000,h=1000
  Box x=0,y=0,w=1000,h=200
    Box x=400,y=0,w=200,h=200
`)
		})
	})
}

func ExpectLayout(tb testing.TB, vnode VNode, size image.Point, expect string) {
	tb.Helper()

	r := renderer.CreateRoot(app.NewWindow())
	r.Render(context.Background(), vnode)
	gtx := giolayout.NewContext(&op.Ops{}, system.FrameEvent{Size: size})
	r.WindowNode().(Element).Layout(gtx)

	buf := bytes.NewBuffer(nil)
	layout.PrintBoundingRectTo(buf, r.RootNode().FirstChild())

	testingx.Expect(tb, strings.TrimSpace(buf.String()), testingx.Be(strings.TrimSpace(expect)))
}
