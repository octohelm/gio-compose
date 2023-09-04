//go:build preview

package theming

import (
	"fmt"
	"image/color"
	"strings"

	"github.com/octohelm/gio-compose/pkg/layout/arrangement"

	"github.com/octohelm/gio-compose/pkg/component/m3/colorrole"
	"github.com/octohelm/gio-compose/pkg/component/m3/font"
	. "github.com/octohelm/gio-compose/pkg/compose"
	"github.com/octohelm/gio-compose/pkg/compose/modifier"
	"github.com/octohelm/gio-compose/pkg/iter"
	"github.com/octohelm/gio-compose/pkg/paint/f32color"
)

//go:generate giocompose preview . ColorSchemePreview
type ColorSchemePreview struct {
}

func (ColorSchemePreview) Build(ctx BuildContext) VNode {
	t := Context.Extract(ctx)

	return Column(
		modifier.FillMaxSize(),
		modifier.PaddingAll(20),
		modifier.VerticalScroll(),
	).Children(
		Fragment(iter.Map([]string{"primary", "secondary", "tertiary", "error"}, func(key string) VNode {
			return Column(
				modifier.FillMaxWidth(),
				modifier.Height(320),
			).Children(
				iter.Map([]colorrole.ColorRole{
					colorrole.ColorRole(fmt.Sprintf(key)),
					colorrole.ColorRole(fmt.Sprintf("on-%s", key)),
					colorrole.ColorRole(fmt.Sprintf("%s-container", key)),
					colorrole.ColorRole(fmt.Sprintf("on-%s-container", key)),
				}, func(colorRole colorrole.ColorRole) VNode {
					isOnColorRole := strings.HasPrefix(string(colorRole), "on-")

					return Row(
						modifier.Weight(func() float32 {
							if isOnColorRole {
								return 2
							}
							return 3
						}()),
						modifier.FillMaxWidth(),
						modifier.Arrangement(arrangement.EqualWeight),
						modifier.PaddingHorizontal(20),
						modifier.PaddingVertical(16),
						modifier.BackgroundColor(t.Color(colorRole)),
					).Children(
						Text(
							string(colorRole),
							modifier.TextColor(func() color.Color {
								if isOnColorRole {
									return t.Color(colorRole[3:])
								}
								return t.Color(colorrole.ColorRole(fmt.Sprintf("on-%s", colorRole)))
							}()),
						),
						Text(
							f32color.HexString(t.Color(colorRole)),
							modifier.FontFamily(t.Font(font.Code)),
							modifier.TextColor(func() color.Color {
								if isOnColorRole {
									return t.Color(colorRole[3:])
								}
								return t.Color(colorrole.ColorRole(fmt.Sprintf("on-%s", colorRole)))
							}()),
						),
					)
				})...,
			)
		})...),

		Fragment(iter.Map([]colorrole.ColorRole{
			colorrole.SurfaceDim,
			colorrole.Surface,
			colorrole.SurfaceBright,

			colorrole.SurfaceContainerLowest,
			colorrole.SurfaceContainerLow,
			colorrole.SurfaceContainer,
			colorrole.SurfaceContainerHigh,
			colorrole.SurfaceContainerHighest,

			colorrole.OnSurface,
			colorrole.OnSurfaceVariant,
			colorrole.Outline,
			colorrole.OutlineVariant,
		}, func(colorRole colorrole.ColorRole) VNode {
			isOnColorRole := strings.HasPrefix(string(colorRole), "on-")

			return Row(
				modifier.FillMaxWidth(),
				modifier.Arrangement(arrangement.EqualWeight),
				modifier.Height(100),
				modifier.PaddingHorizontal(20),
				modifier.PaddingVertical(16),
				modifier.BackgroundColor(t.Color(colorRole)),
			).Children(
				Text(
					string(colorRole),
					modifier.TextColor(func() color.Color {
						if isOnColorRole {
							return t.Color(colorrole.InverseOnSurface)
						}
						return t.Color(colorrole.OnSurface)
					}()),
				),
				Text(
					f32color.HexString(t.Color(colorRole)),
					modifier.FontFamily(t.Font(font.Code)),
					modifier.TextColor(func() color.Color {
						if isOnColorRole {
							return t.Color(colorrole.InverseOnSurface)
						}
						return t.Color(colorrole.OnSurface)
					}()),
				),
			)
		})...),
	)
}
