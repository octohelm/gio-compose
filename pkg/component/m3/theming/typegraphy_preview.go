//go:build preview

package theming

import (
	"strings"

	"github.com/octohelm/gio-compose/pkg/component/m3/colorrole"
	"github.com/octohelm/gio-compose/pkg/component/m3/typescale"
	. "github.com/octohelm/gio-compose/pkg/compose"
	"github.com/octohelm/gio-compose/pkg/compose/modifier"
	"github.com/octohelm/gio-compose/pkg/iter"
)

//go:generate giocompose preview . TypographyPreview
type TypographyPreview struct {
}

func (TypographyPreview) Build(ctx BuildContext) VNode {
	t := Context.Extract(ctx)

	return Column(
		modifier.FillMaxSize(),
		modifier.PaddingAll(40),
		modifier.VerticalScroll(),
	).Children(
		Fragment(iter.Map([]typescale.TypeScale{
			typescale.DisplayLarge,
			typescale.DisplayMedium,
			typescale.DisplaySmall,
			typescale.HeadlineLarge,
			typescale.HeadlineMedium,
			typescale.HeadlineSmall,
			typescale.TitleLarge,
			typescale.TitleMedium,
			typescale.TitleSmall,
			typescale.LabelLarge,
			typescale.LabelMedium,
			typescale.LabelSmall,
			typescale.BodyLarge,
			typescale.BodyMedium,
			typescale.BodySmall,
		}, func(ts typescale.TypeScale) VNode {
			textStyle := t.TypeScale(ts)

			return ProvideTextStyle(modifier.TextStyle(textStyle)).Children(
				Box(
					modifier.BackgroundColor(t.Color(colorrole.SurfaceContainer)),
					modifier.PaddingHorizontal(20),
					modifier.PaddingVertical(8),
				).Children(
					Text(
						displayTypeScale(ts),
					),
				),
				Box(
					modifier.PaddingAll(20),
				).Children(
					Text("SayHi 中文测试"),
				),
			)
		})...),
	)
}

func displayTypeScale(scale typescale.TypeScale) string {
	s := string(scale)

	return strings.ToUpper(s[0:1]) + s[1:]
}
