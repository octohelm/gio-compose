package m3

import (
	"image/color"

	"github.com/octohelm/gio-compose/pkg/layout/arrangement"

	"github.com/octohelm/gio-compose/pkg/component/m3/colorrole"
	"github.com/octohelm/gio-compose/pkg/component/m3/theming"
	"github.com/octohelm/gio-compose/pkg/component/m3/typescale"
	. "github.com/octohelm/gio-compose/pkg/compose"
	"github.com/octohelm/gio-compose/pkg/compose/modifier"
	"github.com/octohelm/gio-compose/pkg/event/gesture"
	"github.com/octohelm/gio-compose/pkg/layout/alignment"
	"github.com/octohelm/gio-compose/pkg/paint/f32color"
)

type ElevatedButton struct {
	LeadingIcon  VNode
	Label        VNode
	TrailingIcon VNode

	Disabled bool
}

func (b ElevatedButton) Build(ctx BuildContext) VNode {
	t := theming.Context.Extract(ctx)

	return H(
		button{
			Disabled:        b.Disabled,
			HasLeadingIcon:  b.LeadingIcon != nil,
			HasTrailingIcon: b.TrailingIcon != nil,
			Color:           t.Color(colorrole.Primary),
			BgColor:         t.Color(colorrole.SurfaceContainerLow),
		},

		ctx.Modifiers(),
	).Children(
		CloneNode(b.LeadingIcon, modifier.Size(18)),
		b.Label,
		CloneNode(b.TrailingIcon, modifier.Size(18)),
	)
}

type FilledButton struct {
	LeadingIcon  VNode
	Label        VNode
	TrailingIcon VNode
	Disabled     bool
}

func (b FilledButton) Build(ctx BuildContext) VNode {
	t := theming.Context.Extract(ctx)

	return H(
		button{
			Disabled:        b.Disabled,
			HasLeadingIcon:  b.LeadingIcon != nil,
			HasTrailingIcon: b.TrailingIcon != nil,
			Color:           t.Color(colorrole.OnPrimary),
			BgColor:         t.Color(colorrole.Primary),
		},

		ctx.Modifiers(),
	).Children(
		CloneNode(b.LeadingIcon, modifier.Size(18)),
		b.Label,
		CloneNode(b.TrailingIcon, modifier.Size(18)),
	)
}

type FilledTonalButton struct {
	LeadingIcon  VNode
	Label        VNode
	TrailingIcon VNode
	Disabled     bool
}

func (b FilledTonalButton) Build(ctx BuildContext) VNode {
	t := theming.Context.Extract(ctx)

	return H(
		button{
			Disabled:        b.Disabled,
			HasLeadingIcon:  b.LeadingIcon != nil,
			HasTrailingIcon: b.TrailingIcon != nil,
			Color:           t.Color(colorrole.OnSecondaryContainer),
			BgColor:         t.Color(colorrole.SecondaryContainer),
		},

		ctx.Modifiers(),
	).Children(
		CloneNode(b.LeadingIcon, modifier.Size(18)),
		b.Label,
		CloneNode(b.TrailingIcon, modifier.Size(18)),
	)
}

type OutlinedButton struct {
	LeadingIcon  VNode
	Label        VNode
	TrailingIcon VNode
	Disabled     bool
}

func (b OutlinedButton) Build(ctx BuildContext) VNode {
	t := theming.Context.Extract(ctx)

	return H(
		button{
			Disabled:        b.Disabled,
			HasLeadingIcon:  b.LeadingIcon != nil,
			HasTrailingIcon: b.TrailingIcon != nil,
			NoShadow:        true,
			Color:           t.Color(colorrole.Primary),
		},

		modifier.BorderStrokeAll(1, t.Color(colorrole.Outline)),
		modifier.When(
			b.Disabled,
			modifier.BorderStrokeAll(1, f32color.Alpha(t.Color(colorrole.OnSurface), 0.12)),
		),

		ctx.Modifiers(),
	).Children(
		CloneNode(b.LeadingIcon, modifier.Size(18)),
		b.Label,
		CloneNode(b.TrailingIcon, modifier.Size(18)),
	)
}

type TextButton struct {
	LeadingIcon  VNode
	Label        VNode
	TrailingIcon VNode
	Disabled     bool
}

func (b TextButton) Build(ctx BuildContext) VNode {
	t := theming.Context.Extract(ctx)

	return H(
		button{
			Disabled:        b.Disabled,
			HasLeadingIcon:  b.LeadingIcon != nil,
			HasTrailingIcon: b.TrailingIcon != nil,
			Color:           t.Color(colorrole.Primary),
			SmallPadding:    true,
			NoShadow:        true,
		},

		ctx.Modifiers(),
	).Children(
		CloneNode(b.LeadingIcon, modifier.Size(18)),
		b.Label,
		CloneNode(b.TrailingIcon, modifier.Size(18)),
	)
}

type button struct {
	Disabled        bool
	HasLeadingIcon  bool
	HasTrailingIcon bool
	NoShadow        bool
	SmallPadding    bool
	IconButton      bool

	Color   color.NRGBA
	BgColor color.NRGBA
}

func (b button) Build(ctx BuildContext) VNode {
	t := theming.Context.Extract(ctx)

	btnState := UseState(ctx, buttonState{})

	return H(StateLayerContainer{
		Color: func() color.NRGBA {
			if btnState.Value().Pressed {
				return f32color.Alpha(b.Color, 0.12)
			}
			if btnState.Value().Focused {
				return f32color.Alpha(b.Color, 0.12)
			}
			if btnState.Value().Hovered {
				return f32color.Alpha(b.Color, 0.08)
			}
			return color.NRGBA{}
		}(),
	},
		modifier.DisplayName("ButtonContainer"),

		modifier.ProvideTextStyle(
			modifier.TextStyle(t.TypeScale(typescale.LabelLarge)),
			modifier.TextColor(b.Color),
		),
		modifier.BackgroundColor(b.BgColor),

		modifier.RoundedAll(20),
		modifier.Shadow(0),
		modifier.When(
			!b.NoShadow && (btnState.Value().Pressed || btnState.Value().Hovered || btnState.Value().Focused),
			modifier.Shadow(1),
		),

		modifier.When(
			b.Disabled,

			modifier.GestureDisabled(),
			modifier.Shadow(0),

			modifier.When(
				!f32color.IsTransparent(b.BgColor),

				modifier.BackgroundColor(f32color.Alpha(t.Color(colorrole.OnSurface), 0.12)),
			),

			modifier.ProvideTextStyle(modifier.TextColor(f32color.Alpha(t.Color(colorrole.OnSurface), 0.38))),
		),

		modifier.DetectGesture(
			gesture.OnPress(func() {
				btnState.UpdateFunc(func(prev buttonState) buttonState {
					prev.Pressed = true
					return prev
				})
			}),
			gesture.OnRelease(func() {
				btnState.UpdateFunc(func(prev buttonState) buttonState {
					prev.Pressed = false
					return prev
				})
			}),
			gesture.OnMouseEnter(func() {
				btnState.UpdateFunc(func(prev buttonState) buttonState {
					prev.Hovered = true
					return prev
				})
			}),
			gesture.OnMouseLeave(func() {
				btnState.UpdateFunc(func(prev buttonState) buttonState {
					prev.Hovered = false
					return prev
				})
			}),
			gesture.OnFocus(func() {
				btnState.UpdateFunc(func(prev buttonState) buttonState {
					prev.Focused = true
					return prev
				})
			}),
			gesture.OnBlur(func() {
				btnState.UpdateFunc(func(prev buttonState) buttonState {
					prev.Focused = false
					return prev
				})
			}),
		),

		ctx.Modifiers(),
	).Children(
		Row(
			modifier.DisplayName("Button"),
			modifier.Align(alignment.Center),
			modifier.Gap(8),
			modifier.Height(40),

			modifier.PaddingHorizontal(24),
			modifier.When(b.HasLeadingIcon, modifier.PaddingLeft(16)),
			modifier.When(b.HasTrailingIcon, modifier.PaddingRight(16)),
			modifier.When(b.SmallPadding, modifier.PaddingHorizontal(12)),

			modifier.When(b.IconButton,
				modifier.FillMaxWidth(),
				modifier.PaddingAll(0),
				modifier.Arrangement(arrangement.SpaceAround),
			),
		).Children(
			ctx.ChildVNodes()...,
		),
	)
}

type buttonState struct {
	Pressed bool
	Focused bool
	Hovered bool
}
