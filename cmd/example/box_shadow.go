package main

import (
	"fmt"
	"image/color"

	"github.com/octohelm/gio-compose/pkg/text"

	. "github.com/octohelm/gio-compose/pkg/compose"
	"github.com/octohelm/gio-compose/pkg/compose/modifier"
	"github.com/octohelm/gio-compose/pkg/layout/alignment"
	"github.com/octohelm/gio-compose/pkg/layout/arrangement"
	"github.com/octohelm/gio-compose/pkg/unit"
)

func init() {
	AddShowCase("paint/shadows", BoxShadow{})
}

type BoxShadow struct{}

func (f BoxShadow) Build(context BuildContext) VNode {
	return Column(
		modifier.FillMaxSize(),
		modifier.Arrangement(arrangement.EqualWeight),
	).Children(
		Row(
			modifier.Align(alignment.Center),
			modifier.Arrangement(arrangement.SpaceAround),
			modifier.FillMaxWidth(),
		).Children(
			MapIndexed([]unit.Dp{0, 1, 2, 3, 5, 6, 7, 8, 12, 18, 22}, func(e unit.Dp, i int) VNode {
				return Box(
					modifier.Size(40),
					modifier.RoundedTopLeft(10),
					modifier.RoundedBottomRight(10),
					modifier.BackgroundColor(color.White),
					modifier.Shadow(e),
				).Children(
					Text(fmt.Sprint(e), modifier.TextAlign(text.Middle)),
				)
			})...,
		),
		Row(
			modifier.Align(alignment.Center),
			modifier.Arrangement(arrangement.SpaceAround),
			modifier.FillMaxWidth(),
		).Children(
			MapIndexed([]unit.Dp{0, 1, 3, 6, 12, 18, 22}, func(e unit.Dp, i int) VNode {
				return Box(
					modifier.Size(40),
					modifier.RoundedRight(10),
					modifier.Shadow(e),
				).Children(
					Text(fmt.Sprint(e), modifier.TextAlign(text.Middle)),
				)
			})...,
		),
	)
}
