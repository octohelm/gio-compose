package main

import (
	"golang.org/x/image/colornames"

	. "github.com/octohelm/gio-compose/pkg/compose"
	"github.com/octohelm/gio-compose/pkg/compose/modifier"
	"github.com/octohelm/gio-compose/pkg/layout/alignment"
)

func init() {
	AddShowCase("layout/absolute", LayoutAbsolute{})
}

type LayoutAbsolute struct {
}

func (LayoutAbsolute) Build(context BuildContext) VNode {
	return Box(
		modifier.FillMaxSize(),
		modifier.BackgroundColor(colornames.Yellowgreen),
		modifier.PaddingAll(20),
	).Children(
		Box(
			modifier.BackgroundColor(colornames.Aliceblue),
			modifier.Align(alignment.TopStart),
		).Children(
			Box(
				modifier.BackgroundColor(colornames.Lightblue),
				modifier.PaddingAll(10),
			).Children(
				Box(
					modifier.BackgroundColor(colornames.Lightpink),
					modifier.Width(80),
					modifier.Height(60),
				),
			),
		),
		Box(
			modifier.Size(100),
			modifier.BackgroundColor(colornames.Aliceblue),
			modifier.Align(alignment.TopEnd),
		),
		Box(
			modifier.Size(100),
			modifier.BackgroundColor(colornames.Aliceblue),
			modifier.Align(alignment.BottomStart)),

		Box(
			modifier.Size(100),
			modifier.BackgroundColor(colornames.Aliceblue),
			modifier.Align(alignment.BottomEnd)),

		Box(
			modifier.Size(100),
			modifier.BackgroundColor(colornames.Beige),
			modifier.Align(alignment.Center),
		),

		Box(
			modifier.Size(100),
			modifier.BackgroundColor(colornames.Beige),
			modifier.Align(alignment.Start),
		),
		Box(
			modifier.Size(100),
			modifier.BackgroundColor(colornames.Beige),
			modifier.Align(alignment.End),
		),

		Box(
			modifier.Size(100),
			modifier.BackgroundColor(colornames.Beige),
			modifier.Align(alignment.Top),
		),

		Box(
			modifier.Size(100),
			modifier.BackgroundColor(colornames.Beige),
			modifier.Align(alignment.Bottom),
		),

		Box(
			modifier.Size(100),
			modifier.BackgroundColor(colornames.Lightblue),
			modifier.Align(alignment.Center),
			modifier.Offset(50),
			modifier.BorderStrokeAll(2, colornames.Black),
			modifier.PaddingAll(10),
			modifier.RoundedAll(50),
		).Children(
			Box(
				modifier.FillMaxSize(),
				modifier.BackgroundColor(colornames.Royalblue),
				modifier.RoundedAll(20),
			),
		),
	)
}
