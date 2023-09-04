package m3

import (
	"image/color"

	"github.com/octohelm/x/reflect"

	"github.com/octohelm/x/encoding"

	"github.com/octohelm/gio-compose/pkg/component/m3/colorrole"
	"github.com/octohelm/gio-compose/pkg/component/m3/theming"
	"github.com/octohelm/gio-compose/pkg/component/m3/typescale"
	. "github.com/octohelm/gio-compose/pkg/compose"
	"github.com/octohelm/gio-compose/pkg/compose/modifier"
	"github.com/octohelm/gio-compose/pkg/event/gesture"
	"github.com/octohelm/gio-compose/pkg/layout/alignment"
	"github.com/octohelm/gio-compose/pkg/paint/f32color"
	"github.com/octohelm/gio-compose/pkg/unit"
)

// FilledTextField
//
// support encoding.MarshalText/UnmarshalText
type FilledTextField[T any] struct {
	LeadingIcon  VNode
	Label        VNode
	TrailingIcon VNode
	Supporting   VNode

	Disabled bool
	HasError bool

	Value         T
	OnValueChange func(v T)
}

func (t FilledTextField[T]) Build(c BuildContext) VNode {
	tm := theming.Context.Extract(c)

	hasValue := !reflect.IsEmptyValue(t.Value)

	state := UseState(c, inputState{})

	hovered := state.Value().Hovered
	populated := state.Value().Focused || hasValue

	return Column(
		modifier.DisplayName("Container"),
		modifier.FillMaxWidth(),
		modifier.Gap(4),

		modifier.ProvideTextStyle(
			modifier.TextStyle(tm.TypeScale(typescale.BodyLarge)),
			modifier.TextColor(Echo(func() color.Color {
				if t.Disabled {
					return f32color.Alpha(tm.Color(colorrole.OnSurface), 0.38)
				}

				if t.HasError {
					if hovered {
						return tm.Color(colorrole.OnErrorContainer)
					}
					return tm.Color(colorrole.Error)
				}

				return tm.Color(colorrole.OnSurfaceVariant)
			})),
		),
	).Children(
		H(StateLayerContainer{
			Color: Echo(func() color.NRGBA {
				if hovered {
					return f32color.Alpha(tm.Color(colorrole.OnSurface), 0.08)
				}
				return color.NRGBA{}
			}),
		},
			modifier.RoundedTop(4),
			modifier.FillMaxWidth(),

			modifier.BorderStrokeBottom(Echo2(func() (borderWidth unit.Dp, c color.Color) {
				borderWidth = 1
				c = tm.Color(colorrole.OnSurfaceVariant)

				if hovered {
					c = tm.Color(colorrole.OnSurface)
				} else if populated {
					c = tm.Color(colorrole.Primary)
				}

				if t.HasError {
					c = tm.Color(colorrole.Error)

					if hovered {
						c = tm.Color(colorrole.OnErrorContainer)
					}
				} else if populated {
					borderWidth = 2
				}

				if hovered {
					borderWidth = 1
				} else if populated {
					borderWidth = 2
				}

				if t.Disabled {
					borderWidth = 1
					c = f32color.Alpha(tm.Color(colorrole.OnSurface), 0.38)
				}

				return
			})),

			modifier.BackgroundColor(Echo(func() color.Color {
				if t.Disabled {
					return f32color.Alpha(tm.Color(colorrole.OnSurface), 0.04)
				}
				return tm.Color(colorrole.SurfaceContainerHighest)
			})),
		).Children(
			Row(
				modifier.FillMaxWidth(),
				modifier.Height(56),
				modifier.Gap(16),
				modifier.Align(alignment.Center),

				modifier.PaddingVertical(8),
				modifier.PaddingHorizontal(16),
				modifier.When(t.LeadingIcon != nil, modifier.PaddingLeft(12)),
				modifier.When(t.TrailingIcon != nil, modifier.PaddingRight(12)),

				modifier.DetectGesture(
					gesture.OnMouseEnter(func() {
						state.UpdateFunc(func(prev inputState) inputState {
							prev.Hovered = true
							return prev
						})
					}),
					gesture.OnMouseLeave(func() {
						state.UpdateFunc(func(prev inputState) inputState {
							prev.Hovered = false
							return prev
						})
					}),
				),

				modifier.When(t.Disabled, modifier.GestureDisabled()),
			).Children(
				CloneNode(
					t.LeadingIcon,

					modifier.DisplayName("LeadingIcon"),
					modifier.Size(24),
					modifier.ProvideTextStyle(
						modifier.TextColor(Echo(func() color.Color {
							if t.Disabled {
								return f32color.Alpha(tm.Color(colorrole.OnSurface), 0.38)
							}
							return tm.Color(colorrole.OnSurfaceVariant)
						})),
					),
				),
				Column(
					modifier.Weight(1),
				).Children(
					Echo(func() VNode {
						if t.Label == nil {
							return nil
						}

						return Box(
							modifier.Key("Label"),

							modifier.DisplayName("Label"),

							modifier.FillMaxWidth(),

							modifier.When(!populated,
								modifier.OffsetY(unit.Dp(tm.TypeScale(typescale.BodySmall).FontSize)),
							),

							modifier.When[any](populated,
								modifier.ProvideTextStyle(
									modifier.TextStyle(tm.TypeScale(typescale.BodySmall)),
									modifier.TextColor(Echo(func() color.Color {
										if !t.HasError && !t.Disabled {
											return tm.Color(colorrole.Primary)
										}
										// inherit parent
										return nil
									})),
								),
							),
						).Children(
							t.Label,
						)
					}),
					Input(
						modifier.Key("Input"),

						modifier.ProvideTextStyle(
							modifier.TextColor(tm.Color(colorrole.OnSurface)),
						),

						modifier.Weight(1),
						modifier.FillMaxWidth(),

						modifier.Value(t.Value, func(v T) string {
							d, _ := encoding.MarshalText(v)
							str := string(d)
							return str
						}),

						modifier.When(t.OnValueChange != nil, modifier.OnValueChange(func(newValue string) {
							v := new(T)
							// FIXME handle error
							_ = encoding.UnmarshalText(v, []byte(newValue))
							t.OnValueChange(*v)
						})),

						modifier.DetectGesture(
							gesture.OnFocus(func() {
								state.UpdateFunc(func(prev inputState) inputState {
									prev.Focused = true
									return prev
								})
							}),
							gesture.OnBlur(func() {
								state.UpdateFunc(func(prev inputState) inputState {
									prev.Focused = false
									return prev
								})
							}),
						),

						modifier.When(t.Disabled, modifier.GestureDisabled()),

						c.Modifiers(),
					),
				),
				CloneNode(
					t.TrailingIcon,
					modifier.Size(24),
					modifier.DisplayName("TrailingIcon"),
				),
			),
		),
		Echo(func() VNode {
			if t.Supporting == nil {
				return nil
			}

			return Box(
				modifier.DisplayName("Supporting"),
				modifier.FillMaxWidth(),
				modifier.PaddingHorizontal(16),
				modifier.PaddingVertical(2),
				modifier.When(t.LeadingIcon != nil, modifier.PaddingLeft(12)),
				modifier.When(t.TrailingIcon != nil, modifier.PaddingRight(12)),

				modifier.ProvideTextStyle(
					modifier.TextStyle(tm.TypeScale(typescale.BodySmall)),
					modifier.When(t.HasError && hovered, modifier.TextColor(tm.Color(colorrole.Error))),
				),
			).Children(
				t.Supporting,
			)
		}),
	)
}

// OutlinedTextField
//
// support encoding.MarshalText/UnmarshalText
type OutlinedTextField[T any] struct {
	LeadingIcon  VNode
	Label        VNode
	TrailingIcon VNode
	Supporting   VNode

	ContainerColor color.NRGBA
	Disabled       bool
	HasError       bool
	Value          T
	OnValueChange  func(v T)
}

func (t OutlinedTextField[T]) Build(c BuildContext) VNode {
	tm := theming.Context.Extract(c)

	containerColor := t.ContainerColor
	if containerColor.A == 0 {
		containerColor = tm.Color(colorrole.SurfaceContainerLowest)
	}

	hasValue := !reflect.IsEmptyValue(t.Value)

	state := UseState(c, inputState{})

	hovered := state.Value().Hovered
	populated := state.Value().Focused || hasValue

	return Column(
		modifier.FillMaxWidth(),
		modifier.Gap(4),

		modifier.ProvideTextStyle(
			modifier.TextStyle(tm.TypeScale(typescale.BodyLarge)),
			modifier.TextColor(Echo(func() color.Color {
				if t.Disabled {
					return f32color.Alpha(tm.Color(colorrole.OnSurface), 0.38)
				}

				if t.HasError {
					if hovered {
						return tm.Color(colorrole.OnErrorContainer)
					}
					return tm.Color(colorrole.Error)
				}

				return tm.Color(colorrole.OnSurfaceVariant)
			})),
		),
	).Children(
		Box(
			modifier.FillMaxWidth(),
		).Children(
			Box(
				modifier.DisplayName("Container"),

				modifier.RoundedAll(4),
				modifier.FillMaxWidth(),

				modifier.BorderStrokeAll(Echo2(func() (borderWidth unit.Dp, c color.Color) {
					borderWidth = 1
					c = tm.Color(colorrole.Outline)

					if hovered {
						c = tm.Color(colorrole.OnSurface)
					} else if populated {
						c = tm.Color(colorrole.Primary)
					}

					if t.HasError {
						c = tm.Color(colorrole.Error)

						if hovered {
							c = tm.Color(colorrole.OnErrorContainer)
						}
					} else if populated {
						borderWidth = 2
					}

					if hovered {
						borderWidth = 1
					} else if populated {
						borderWidth = 2
					}

					if t.Disabled {
						borderWidth = 1
						c = f32color.Alpha(tm.Color(colorrole.OnSurface), 0.38)
					}

					return
				})),
			).Children(
				Row(
					modifier.FillMaxWidth(),
					modifier.Height(56),
					modifier.Gap(16),
					modifier.Align(alignment.Center),

					modifier.PaddingVertical(8),
					modifier.PaddingHorizontal(16),
					modifier.When(t.LeadingIcon != nil, modifier.PaddingLeft(12)),
					modifier.When(t.TrailingIcon != nil, modifier.PaddingRight(12)),

					modifier.DetectGesture(
						gesture.OnMouseEnter(func() {
							state.UpdateFunc(func(prev inputState) inputState {
								prev.Hovered = true
								return prev
							})
						}),
						gesture.OnMouseLeave(func() {
							state.UpdateFunc(func(prev inputState) inputState {
								prev.Hovered = false
								return prev
							})
						}),
					),

					modifier.When(t.Disabled, modifier.GestureDisabled()),
				).Children(
					CloneNode(
						t.LeadingIcon,

						modifier.Size(24),

						modifier.ProvideTextStyle(
							modifier.TextColor(Echo(func() color.Color {
								if t.Disabled {
									return f32color.Alpha(tm.Color(colorrole.OnSurface), 0.38)
								}
								return tm.Color(colorrole.OnSurfaceVariant)
							})),
						),
					),
					Input(
						modifier.Weight(1),
						modifier.Key("Input"),

						modifier.ProvideTextStyle(
							modifier.TextColor(tm.Color(colorrole.OnSurface)),
						),

						modifier.Weight(1),
						modifier.FillMaxWidth(),

						modifier.Value(t.Value, func(v T) string {
							d, _ := encoding.MarshalText(v)
							str := string(d)
							return str
						}),

						modifier.When(t.OnValueChange != nil, modifier.OnValueChange(func(newValue string) {
							v := new(T)
							// FIXME handle error
							_ = encoding.UnmarshalText(v, []byte(newValue))
							t.OnValueChange(*v)
						})),

						modifier.DetectGesture(
							gesture.OnFocus(func() {
								state.UpdateFunc(func(prev inputState) inputState {
									prev.Focused = true
									return prev
								})
							}),
							gesture.OnBlur(func() {
								state.UpdateFunc(func(prev inputState) inputState {
									prev.Focused = false
									return prev
								})
							}),
						),

						modifier.When(t.Disabled, modifier.GestureDisabled()),

						c.Modifiers(),
					),
					CloneNode(
						t.TrailingIcon,
						modifier.Size(24),
					),
				),
			),

			Echo(func() VNode {
				if t.Label == nil {
					return nil
				}

				return Box(
					modifier.DisplayName("LabelContainer"),

					modifier.Align(alignment.TopStart),
					modifier.When(populated,
						modifier.OffsetY(unit.Dp(-tm.TypeScale(typescale.BodySmall).FontSize/2)),
						modifier.OffsetX(Echo(func() unit.Dp {
							if t.LeadingIcon != nil {
								return 12
							}
							return 16
						})),
					),

					modifier.When(!populated,
						modifier.Align(alignment.Start),
						modifier.OffsetX(Echo(func() unit.Dp {
							if t.LeadingIcon != nil {
								// padding and icon size
								return 12 + 16 + 24 - 4
							}
							return 16
						})),
					),
				).Children(
					Box(
						modifier.PaddingHorizontal(4),
						// FIXME support inherit
						modifier.BackgroundColor(containerColor),

						modifier.When[any](populated,
							modifier.ProvideTextStyle(
								modifier.TextStyle(tm.TypeScale(typescale.BodySmall)),
								modifier.TextColor(Echo(func() color.Color {
									if !t.HasError && !t.Disabled {
										return tm.Color(colorrole.Primary)
									}
									// inherit parent
									return nil
								})),
							),
						),
					).Children(
						t.Label,
					),
				)
			}),
		),
		Echo(func() VNode {
			if t.Supporting == nil {
				return nil
			}

			return Box(
				modifier.DisplayName("Supporting"),
				modifier.FillMaxWidth(),
				modifier.PaddingHorizontal(16),
				modifier.PaddingVertical(2),
				modifier.When(t.LeadingIcon != nil, modifier.PaddingLeft(12)),
				modifier.When(t.TrailingIcon != nil, modifier.PaddingRight(12)),

				modifier.ProvideTextStyle(
					modifier.TextStyle(tm.TypeScale(typescale.BodySmall)),
					modifier.When(t.HasError && hovered, modifier.TextColor(tm.Color(colorrole.Error))),
				),
			).Children(
				t.Supporting,
			)
		}),
	)
}

type inputState struct {
	Focused bool
	Hovered bool
}
