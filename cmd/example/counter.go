package main

import (
	"fmt"

	"github.com/octohelm/gio-compose/pkg/component/m3"
	. "github.com/octohelm/gio-compose/pkg/compose"
	"github.com/octohelm/gio-compose/pkg/compose/modifier"
	"github.com/octohelm/gio-compose/pkg/event/gesture"
	"github.com/octohelm/gio-compose/pkg/layout/alignment"
	"github.com/octohelm/gio-compose/pkg/layout/arrangement"
	"github.com/octohelm/gio-compose/pkg/text"
)

func init() {
	AddShowCase("demo/counter", Counter{})
}

type Counter struct {
}

func (Counter) Build(b BuildContext) VNode {
	state := UseState(b, 0)

	v :=
		Row(
			modifier.FillMaxSize(),
			modifier.Align(alignment.Center),
			modifier.Arrangement(arrangement.SpaceAround),
		).Children(
			H(
				m3.ElevatedButton{
					Label: Text(fmt.Sprintf("Counts %d", state.Value()), modifier.TextAlign(text.Middle)),
				},

				modifier.DetectGesture(gesture.OnTap(func() {
					state.UpdateFunc(func(prev int) int {
						return prev + 1
					})
				})),
			),
		)
	return v
}
