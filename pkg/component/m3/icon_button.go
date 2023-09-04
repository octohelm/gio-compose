package m3

import (
	"image/color"

	"github.com/octohelm/gio-compose/pkg/component/m3/colorrole"
	"github.com/octohelm/gio-compose/pkg/component/m3/theming"
	. "github.com/octohelm/gio-compose/pkg/compose"
	"github.com/octohelm/gio-compose/pkg/compose/modifier"
	"github.com/octohelm/gio-compose/pkg/paint/f32color"
)

type FilledIconButton struct {
	Icon       VNode
	Disabled   bool
	UnSelected bool
}

func (b FilledIconButton) Build(ctx BuildContext) VNode {
	t := theming.Context.Extract(ctx)

	return H(
		button{
			Disabled:   b.Disabled,
			NoShadow:   true,
			IconButton: true,

			Color: Echo(func() color.NRGBA {
				if b.UnSelected {
					return t.Color(colorrole.Primary)
				}
				return t.Color(colorrole.OnPrimary)
			}),
			BgColor: Echo(func() color.NRGBA {
				if b.UnSelected {
					return t.Color(colorrole.SurfaceContainerHighest)
				}
				return t.Color(colorrole.Primary)
			}),
		},

		modifier.Size(40),
		ctx.Modifiers(),
	).Children(
		CloneNode(b.Icon, modifier.Size(24)),
	)
}

type FilledTonalIconButton struct {
	Icon       VNode
	Disabled   bool
	UnSelected bool
}

func (b FilledTonalIconButton) Build(ctx BuildContext) VNode {
	t := theming.Context.Extract(ctx)

	return H(
		button{
			Disabled:   b.Disabled,
			IconButton: true,
			NoShadow:   true,

			Color: Echo(func() color.NRGBA {
				if b.UnSelected {
					return t.Color(colorrole.OnSurfaceVariant)
				}
				return t.Color(colorrole.OnSecondaryContainer)
			}),
			BgColor: Echo(func() color.NRGBA {
				if b.UnSelected {
					return t.Color(colorrole.SurfaceContainerHighest)
				}
				return t.Color(colorrole.SecondaryContainer)
			}),
		},

		modifier.Size(40),
		ctx.Modifiers(),
	).Children(
		CloneNode(b.Icon, modifier.Size(24)),
	)
}

type OutlinedIconButton struct {
	Icon       VNode
	Disabled   bool
	UnSelected bool
}

func (b OutlinedIconButton) Build(ctx BuildContext) VNode {
	t := theming.Context.Extract(ctx)

	return H(
		button{
			Disabled:   b.Disabled,
			IconButton: true,
			NoShadow:   true,

			Color: func() color.NRGBA {
				if b.UnSelected {
					return t.Color(colorrole.OnSurfaceVariant)
				}
				return t.Color(colorrole.InverseOnSurface)
			}(),
			BgColor: func() color.NRGBA {
				if b.UnSelected {
					return color.NRGBA{}
				}
				return t.Color(colorrole.InverseSurface)
			}(),
		},

		modifier.BorderStrokeAll(1, t.Color(colorrole.Outline)),
		modifier.When(
			b.Disabled,
			modifier.BorderStrokeAll(1, f32color.Alpha(t.Color(colorrole.OnSurface), 0.12)),
		),

		modifier.Size(40),
		ctx.Modifiers(),
	).Children(
		CloneNode(b.Icon, modifier.Size(24)),
	)
}

type StandardIconButton struct {
	Icon       VNode
	Disabled   bool
	UnSelected bool
}

func (b StandardIconButton) Build(ctx BuildContext) VNode {
	t := theming.Context.Extract(ctx)

	return H(
		button{
			Disabled:   b.Disabled,
			NoShadow:   true,
			IconButton: true,

			Color: func() color.NRGBA {
				if b.UnSelected {
					return t.Color(colorrole.OnSurfaceVariant)
				}
				return t.Color(colorrole.Primary)
			}(),
		},

		modifier.Size(40),
		ctx.Modifiers(),
	).Children(
		CloneNode(b.Icon, modifier.Size(24)),
	)
}
