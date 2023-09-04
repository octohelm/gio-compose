package compose

import (
	"image"
	"testing"

	"github.com/octohelm/gio-compose/pkg/layout/position"
	"github.com/octohelm/gio-compose/pkg/paint/size"

	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	testingx "github.com/octohelm/x/testing"
)

func TestContainer(t *testing.T) {
	t.Run("Sized Container should return dims with self size", func(t *testing.T) {
		t.Run("when zero padding", func(t *testing.T) {
			p := newContainer()
			p.SetSize(200, size.Width, size.Height)

			gtx := layout.NewContext(&op.Ops{}, system.FrameEvent{})

			dims := p.Layout(gtx, nil, ElementPainterFunc(func(gtx layout.Context) layout.Dimensions {
				testingx.Expect(t, gtx.Constraints.Min, testingx.Equal(image.Pt(200, 200)))
				testingx.Expect(t, gtx.Constraints.Max, testingx.Equal(image.Pt(200, 200)))

				return layout.Dimensions{Size: image.Pt(100, 80)}
			}))

			testingx.Expect(t, dims.Size, testingx.Equal(image.Pt(200, 200)))
		})

		t.Run("when with padding", func(t *testing.T) {
			p := newContainer()
			p.SetSize(200, size.Width, size.Height)
			p.SetEdgeInset(10, position.AllFour...)

			gtx := layout.NewContext(&op.Ops{}, system.FrameEvent{})

			dims := p.Layout(gtx, nil, ElementPainterFunc(func(gtx layout.Context) layout.Dimensions {
				testingx.Expect(t, gtx.Constraints.Min, testingx.Equal(image.Pt(200-20, 200-20)))
				testingx.Expect(t, gtx.Constraints.Max, testingx.Equal(image.Pt(200-20, 200-20)))

				return layout.Dimensions{Size: image.Pt(100, 80)}
			}))

			testingx.Expect(t, dims.Size, testingx.Equal(image.Pt(200, 200)))
		})
	})

	t.Run("Un-sized Container should return dims fit child sizes", func(t *testing.T) {
		t.Run("when zero padding", func(t *testing.T) {
			p := newContainer()
			gtx := layout.NewContext(&op.Ops{}, system.FrameEvent{})

			dims := p.Layout(gtx, nil, ElementPainterFunc(func(gtx layout.Context) layout.Dimensions {
				testingx.Expect(t, gtx.Constraints.Min, testingx.Equal(image.Pt(0, 0)))

				return layout.Dimensions{Size: image.Pt(100, 80)}
			}))

			testingx.Expect(t, dims.Size, testingx.Equal(image.Pt(100, 80)))
		})

		t.Run("when padding", func(t *testing.T) {
			p := newContainer()
			p.SetEdgeInset(10, position.AllFour...)

			gtx := layout.NewContext(&op.Ops{}, system.FrameEvent{})

			dims := p.Layout(gtx, nil, ElementPainterFunc(func(gtx layout.Context) layout.Dimensions {
				return layout.Dimensions{Size: image.Pt(100, 80)}
			}))

			testingx.Expect(t, dims.Size, testingx.Equal(image.Pt(100+20, 80+20)))
		})

		t.Run("when padding nested", func(t *testing.T) {
			p := newContainer()
			p.SetSize(200, size.Width, size.Height)
			p.SetEdgeInset(10, position.AllFour...)

			gtx := layout.NewContext(&op.Ops{}, system.FrameEvent{})

			dims := p.Layout(gtx, nil, ElementPainterFunc(func(gtx layout.Context) layout.Dimensions {
				p2 := newContainer()
				p2.SetEdgeInset(10, position.AllFour...)

				childDims := p2.Layout(gtx, nil, ElementPainterFunc(func(gtx layout.Context) layout.Dimensions {
					return layout.Dimensions{Size: image.Pt(100, 80)}
				}))

				testingx.Expect(t, childDims.Size, testingx.Equal(image.Pt(100+20, 80+20)))

				return childDims
			}))

			testingx.Expect(t, dims.Size, testingx.Equal(image.Pt(200, 200)))
		})
	})
}
