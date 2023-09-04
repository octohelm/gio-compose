package main

import (
	"image/color"

	"github.com/octohelm/gio-compose/pkg/layout/arrangement"

	"github.com/octohelm/gio-compose/pkg/component"
	. "github.com/octohelm/gio-compose/pkg/compose"
	"github.com/octohelm/gio-compose/pkg/compose/modifier"
	"github.com/octohelm/gio-compose/pkg/layout/alignment"
	"github.com/octohelm/gio-compose/pkg/unit"
	"golang.org/x/image/colornames"
)

func init() {
	AddShowCase("overlay/positioned", OverlayPositioned{})
}

type OverlayPositioned struct {
}

func (OverlayPositioned) Build(b BuildContext) VNode {
	return Row(
		modifier.FillMaxSize(),
		modifier.Arrangement(arrangement.EqualWeight),
	).Children(
		Box(
			modifier.DisplayName("Container"),
			modifier.FillMaxHeight(),
		).Children(
			Box(
				modifier.DisplayName("TriggerGloballyPositioned"),
				modifier.Align(alignment.Center),
				modifier.Width(80),
				modifier.Height(40),
				modifier.BackgroundColor(colornames.Lightgrey),

				component.Positioned(
					Box(
						modifier.DisplayName("TriggerContent"),
						modifier.Size(100),
						modifier.BackgroundColor(color.RGBA{A: 0x20}),
					),
				),
			),
		),
		Column(
			modifier.VerticalScroll(),
			modifier.DisplayName("ContainerWithScrolling"),
			modifier.FillMaxHeight(),
			modifier.Align(alignment.Center),
		).Children(
			MapIndexed(make([]VNode, 100), func(e VNode, idx int) VNode {
				return Box(
					modifier.FillMaxWidth(),
					modifier.Height(func() unit.Dp {
						if idx%2 == 0 {
							return 40
						}
						return 80
					}()),
					modifier.Align(alignment.Center),
					modifier.BackgroundColor(
						color.RGBA{
							R: func() uint8 {
								if idx%2 == 0 {
									return 0x00
								}
								return 0xff
							}(),
							A: 0x40,
						}),
				).Children(
					Echo(func() VNode {

						if idx != 9 {
							return nil
						}

						return Box(
							modifier.DisplayName("TriggerGloballyPositioned"),
							modifier.Align(alignment.Center),
							modifier.Width(80),
							modifier.FillMaxHeight(),
							modifier.BackgroundColor(colornames.Lightgrey),

							component.Positioned(
								// Content
								Box(
									modifier.DisplayName("TriggerContent"),
									modifier.Size(100),
									modifier.BackgroundColor(color.RGBA{A: 0x20}),
								),
							),
						)
					}),
				)
			})...,
		),
	)
}
