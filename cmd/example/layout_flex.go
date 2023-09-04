package main

import (
	. "github.com/octohelm/gio-compose/pkg/compose"

	"github.com/octohelm/gio-compose/pkg/compose/modifier"
	"github.com/octohelm/gio-compose/pkg/layout/alignment"
	"github.com/octohelm/gio-compose/pkg/layout/arrangement"
	"golang.org/x/image/colornames"
)

func init() {
	AddShowCase("layout/flex", LayoutFlex{})
}

type LayoutFlex struct {
}

func (LayoutFlex) Build(context BuildContext) VNode {
	return Row(modifier.FillMaxSize()).Children(
		Column(
			modifier.Width(200),
			modifier.FillMaxHeight(),
		).Children(
			Box(
				modifier.Weight(1),
				modifier.FillMaxWidth(),
				modifier.BackgroundColor(colornames.Lightblue),
			),
			Box(
				modifier.Height(40),
				modifier.FillMaxWidth(),
				modifier.BackgroundColor(colornames.Lightpink),
			),
		),
		Column(
			modifier.Weight(1),
			modifier.FillMaxHeight(),
			modifier.BackgroundColor(colornames.Lightgreen),
			modifier.PaddingAll(10),
		).Children(
			Box(modifier.FillMaxWidth(), modifier.BackgroundColor(colornames.Black), modifier.Height(40), modifier.FillMaxWidth()),
			Row(modifier.FillMaxWidth(), modifier.Weight(0.5), modifier.PaddingAll(10), modifier.BackgroundColor(colornames.Royalblue)),
			Column(modifier.Arrangement(arrangement.EqualWeight), modifier.FillMaxWidth(), modifier.Weight(3), modifier.Gap(10), modifier.PaddingAll(10), modifier.BackgroundColor(colornames.Rosybrown)).Children(
				H(LayoutRow{Arrangement: arrangement.EqualWeight, Alignment: alignment.Start}),
				H(LayoutRow{Arrangement: arrangement.SpaceBetween, Alignment: alignment.Middle}),
				H(LayoutRow{Arrangement: arrangement.SpaceAround, Alignment: alignment.End}),
				H(LayoutRow{Arrangement: arrangement.SpaceEvenly, Alignment: alignment.Baseline}),
				H(LayoutRow{Arrangement: arrangement.End, Alignment: alignment.Middle}),
				H(LayoutRow{Arrangement: arrangement.Center, Alignment: alignment.Middle}),
				H(LayoutRow{Arrangement: arrangement.Start, Alignment: alignment.Middle}),
			),
		),
	)
}

type LayoutRow struct {
	Arrangement arrangement.Arrangement
	Alignment   alignment.Alignment
}

func (f LayoutRow) Build(context BuildContext) VNode {
	return Row(
		modifier.Align(f.Alignment),
		modifier.Arrangement(f.Arrangement),
		modifier.Gap(10),
		modifier.PaddingAll(10),
		modifier.BackgroundColor(colornames.Aliceblue),
		modifier.FillMaxWidth(),
		modifier.FillMaxHeight(),
	).Children(
		Box(
			modifier.BackgroundColor(colornames.Lightpink),
			modifier.When(f.Arrangement != arrangement.EqualWeight, modifier.Width(80)),
			modifier.Height(10),
		),
		Box(
			modifier.BackgroundColor(colornames.Lightcoral),
			modifier.When(f.Arrangement != arrangement.EqualWeight, modifier.Width(80)),
			modifier.Height(10),
		),
		Box(
			modifier.BackgroundColor(colornames.Lightblue),
			modifier.When(f.Arrangement != arrangement.EqualWeight, modifier.Width(80)),
			modifier.Height(10),
		),
	)
}
