package paint

import (
	"image"

	"gioui.org/f32"
	"gioui.org/op"
	"github.com/octohelm/gio-compose/pkg/paint/contentscale"
)

type ContentScaleSetter interface {
	SetContentScale(scale contentscale.ContentScale)
}

var _ ContentScaleSetter = &ContentScale{}

type ContentScale struct {
	scale contentscale.ContentScale
}

func (s *ContentScale) SetContentScale(scale contentscale.ContentScale) {
	s.scale = scale
}

func (s *ContentScale) Fit(ops *op.Ops, imageSize image.Point, containerSize image.Point) {
	scale := s.transform(imageSize, containerSize)

	trans := f32.Affine2D{}.
		Scale(f32.Pt(
			float32(imageSize.X)/2,
			float32(imageSize.Y)/2,
		), scale).
		Offset(f32.Pt(
			(float32(containerSize.X)-float32(imageSize.X))/2,
			(float32(containerSize.Y)-float32(imageSize.Y))/2),
		)

	op.Affine(trans).Add(ops)
}

func (s *ContentScale) transform(imageSize image.Point, containerSize image.Point) (scale f32.Point) {
	// https://developer.android.com/jetpack/compose/graphics/images/customize

	scaleBaseHeight := func() f32.Point {
		x := float32(containerSize.Y) / float32(imageSize.Y)
		return f32.Pt(x, x)
	}

	scaleBaseWidth := func() f32.Point {
		x := float32(containerSize.X) / float32(imageSize.X)
		return f32.Pt(x, x)
	}

	scaleBaseMax := func() f32.Point {
		if imageSize.X >= imageSize.Y {
			return scaleBaseWidth()
		}
		return scaleBaseHeight()
	}

	scaleBaseMin := func() f32.Point {
		if imageSize.X <= imageSize.Y {
			return scaleBaseWidth()
		}
		return scaleBaseHeight()
	}

	switch s.scale {
	case contentscale.Inside:
		if imageSize.X > containerSize.X || imageSize.Y > containerSize.Y {
			return scaleBaseMax()
		}
	case contentscale.FillBounds:
		return f32.Pt(float32(containerSize.X)/float32(imageSize.X), float32(containerSize.Y)/float32(imageSize.Y))
	case contentscale.FillWidth:
		return scaleBaseWidth()
	case contentscale.FillHeight:
		return scaleBaseHeight()
	case contentscale.Crop:
		return scaleBaseMin()
	case contentscale.Fit:
		return scaleBaseMax()
	}

	return f32.Pt(1, 1)
}
