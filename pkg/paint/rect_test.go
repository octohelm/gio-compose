package paint

import (
	"testing"

	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	testingx "github.com/octohelm/x/testing"

	"github.com/octohelm/gio-compose/pkg/layout/position"
	"github.com/octohelm/gio-compose/pkg/paint/size"
)

func TestRect(t *testing.T) {
	r := NewRect()
	r.SetSize(10, size.Width, size.Height)

	gtx := layout.NewContext(new(op.Ops), system.FrameEvent{})

	p := r.Path(gtx)

	testingx.Expect(t, p.String(), testingx.Be("M0 0L10 0L10 10L0 10z"))
}

func TestRoundedRect(t *testing.T) {
	rr := NewRoundedRect()
	rr.SetSize(10, size.Width, size.Height)
	rr.SetCornerRadius(2, position.AllFour...)

	gtx := layout.NewContext(new(op.Ops), system.FrameEvent{})

	p := rr.Path(gtx,
		position.Top,
		position.Right,
		position.Bottom,
		position.Left,
	)

	testingx.Expect(t, p.String(), testingx.Be("M0.5857864376269051 0.5857864376269051A2 2 0.7853981633974483 0 1 2 0L8 0A2 2 0.7853981633974483 0 1 9.414213562373096 0.5857864376269051A2 2 0.7853981633974483 0 1 10 2L10 8A2 2 0.7853981633974483 0 1 9.414213562373096 9.414213562373096A2 2 0.7853981633974483 0 1 8 10L2 10A2 2 0.7853981633974483 0 1 0.5857864376269051 9.414213562373096A2 2 0.7853981633974483 0 1 0 8L0 2A2 2 0.7853981633974483 0 1 0.5857864376269051 0.5857864376269051z"))
}
