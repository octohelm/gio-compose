package testutil

import (
	"image"
	"image/png"
	"testing"

	"gioui.org/gpu/headless"
	"gioui.org/op"
)

func DrawImage(t *testing.T, sz image.Point, draw func(o *op.Ops)) (im *image.RGBA, err error) {
	w, err := headless.NewWindow(sz.X, sz.Y)
	if err != nil {
		t.Skipf("failed to create headless window, skipping: %v", err)
	}
	defer w.Release()

	ops := new(op.Ops)
	draw(ops)
	if err := w.Frame(ops); err != nil {
		return nil, err
	}
	im = image.NewRGBA(image.Rectangle{Max: sz})
	err = w.Screenshot(im)
	return im, err
}

func WritePNG(filename string, img image.Image) error {
	f, err := OpenFile(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}
