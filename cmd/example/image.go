package main

import (
	"bytes"
	_ "embed"
	"image"
	"image/png"

	"github.com/disintegration/imaging"
	"golang.org/x/image/colornames"

	. "github.com/octohelm/gio-compose/pkg/compose"
	"github.com/octohelm/gio-compose/pkg/compose/modifier"
	"github.com/octohelm/gio-compose/pkg/layout/alignment"
	"github.com/octohelm/gio-compose/pkg/layout/arrangement"
	"github.com/octohelm/gio-compose/pkg/paint/contentscale"
)

func init() {
	AddShowCase("paint/images", Images{})
}

//go:embed assets/hiking.png
var imageData []byte

var imageH image.Image
var imageV image.Image

func init() {
	raw, _ := png.Decode(bytes.NewBuffer(imageData))

	imageH = imaging.Fit(raw, 160, 80, imaging.Lanczos)
	imageV = imaging.Rotate90(imageH)

}

type Images struct{}

func (f Images) Build(context BuildContext) VNode {
	return Column(
		modifier.FillMaxSize(),
		modifier.Gap(20),
		modifier.Align(alignment.Center),
	).Children(
		MapIndexed([]contentscale.ContentScale{
			contentscale.None,
			contentscale.Inside,
			contentscale.FillWidth,
			contentscale.FillHeight,
			contentscale.FillBounds,
			contentscale.Crop,
			contentscale.Fit,
		}, func(scale contentscale.ContentScale, i int) VNode {
			return Row(
				modifier.FillMaxWidth(),
				modifier.Gap(20),
				modifier.Align(alignment.Center),
				modifier.Arrangement(arrangement.SpaceAround),
			).Children(
				Image(imageH, modifier.ContentScale(scale), modifier.Size(100), modifier.RoundedAll(20), modifier.BackgroundColor(colornames.Lightcyan)),
				Image(imageH, modifier.ContentScale(scale), modifier.Size(50), modifier.RoundedAll(10), modifier.BackgroundColor(colornames.Lightgrey)),
				Image(imageV, modifier.ContentScale(scale), modifier.Size(100), modifier.RoundedAll(20), modifier.BackgroundColor(colornames.Lightcyan)),
				Image(imageV, modifier.ContentScale(scale), modifier.Size(50), modifier.RoundedAll(10), modifier.BackgroundColor(colornames.Lightgrey)),
			)
		})...,
	)
}
