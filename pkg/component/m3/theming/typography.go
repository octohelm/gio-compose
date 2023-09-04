package theming

import (
	"strings"

	"github.com/octohelm/gio-compose/pkg/text/weight"

	"github.com/octohelm/gio-compose/pkg/component/m3/font"
	"github.com/octohelm/gio-compose/pkg/component/m3/typescale"
	"github.com/octohelm/gio-compose/pkg/text"
)

var DefaultTypography = NewTypography()

type FontProvider interface {
	Font(font font.Font) string
}

type Typography interface {
	FontProvider
	TypeScale(typeScale typescale.TypeScale) text.Style
}

func WithFonts(fonts Fonts) TypographyOptionFunc {
	return func(t *typography) {
		if t.fonts == nil {
			t.fonts = map[font.Font]string{}
		}

		for f := range fonts {
			t.fonts[f] = strings.Join(fonts[f], ", ")
		}
	}
}

func WithTypeScales(create func(fp FontProvider) TypeScales) TypographyOptionFunc {
	return func(t *typography) {
		t.typeScales = create(t)
	}
}

func NewTypography(opts ...TypographyOptionFunc) Typography {
	t := &typography{}

	WithFonts(defaultFonts)(t)

	for i := range opts {
		opts[i](t)
	}

	if len(t.typeScales) == 0 {
		WithTypeScales(CreateDefaultTypeScales)(t)
	}

	return t
}

type TypographyOptionFunc func(t *typography)

type typography struct {
	fonts      map[font.Font]string
	typeScales TypeScales
}

func (t *typography) Font(font font.Font) string {
	return t.fonts[font]
}

func (t *typography) TypeScale(typeScale typescale.TypeScale) text.Style {
	return t.typeScales[typeScale]
}

type Fonts map[font.Font][]string
type TypeScales map[typescale.TypeScale]text.Style

func CreateDefaultTypeScales(fp FontProvider) TypeScales {
	return TypeScales{
		typescale.DisplayLarge: {
			LineHeight:    64,
			FontSize:      57,
			LetterSpacing: -0.25,
			FontWeight:    weight.Normal,
			FontFamily:    fp.Font(font.Brand),
		},
		typescale.DisplayMedium: {
			LineHeight:    52,
			FontSize:      45,
			LetterSpacing: 0,
			FontWeight:    weight.Normal,
			FontFamily:    fp.Font(font.Brand),
		},
		typescale.DisplaySmall: {
			LineHeight:    44,
			FontSize:      36,
			LetterSpacing: 0,
			FontWeight:    weight.Normal,
			FontFamily:    fp.Font(font.Brand),
		},

		typescale.HeadlineLarge: {
			LineHeight:    40,
			FontSize:      32,
			LetterSpacing: 0,
			FontWeight:    weight.Normal,
			FontFamily:    fp.Font(font.Brand),
		},
		typescale.HeadlineMedium: {
			LineHeight:    36,
			FontSize:      28,
			LetterSpacing: 0,
			FontWeight:    weight.Normal,
			FontFamily:    fp.Font(font.Brand),
		},
		typescale.HeadlineSmall: {
			LineHeight:    28,
			FontSize:      22,
			LetterSpacing: 0,
			FontWeight:    weight.Normal,
			FontFamily:    fp.Font(font.Plain),
		},

		typescale.TitleLarge: {
			LineHeight:    28,
			FontSize:      22,
			LetterSpacing: 0,
			FontWeight:    weight.Normal,
			FontFamily:    fp.Font(font.Brand),
		},
		typescale.TitleMedium: {
			LineHeight:    24,
			FontSize:      16,
			LetterSpacing: 0.15,
			FontWeight:    weight.Medium,
			FontFamily:    fp.Font(font.Plain),
		},
		typescale.TitleSmall: {
			LineHeight:    20,
			FontSize:      14,
			LetterSpacing: 0.15,
			FontWeight:    weight.Medium,
			FontFamily:    fp.Font(font.Plain),
		},

		typescale.LabelLarge: {
			LineHeight:    20,
			FontSize:      14,
			LetterSpacing: 0.1,
			FontWeight:    weight.Medium,
			FontFamily:    fp.Font(font.Plain),
		},
		typescale.LabelMedium: {
			LineHeight:    16,
			FontSize:      12,
			LetterSpacing: 0.5,
			FontWeight:    weight.Medium,
			FontFamily:    fp.Font(font.Plain),
		},

		typescale.LabelSmall: {
			LineHeight:    16,
			FontSize:      11,
			LetterSpacing: 0.5,
			FontWeight:    weight.Medium,
			FontFamily:    fp.Font(font.Plain),
		},

		typescale.BodyLarge: {
			LineHeight:    24,
			FontSize:      16,
			LetterSpacing: 0.5,
			FontWeight:    weight.Normal,
			FontFamily:    fp.Font(font.Plain),
		},
		typescale.BodyMedium: {
			LineHeight:    20,
			FontSize:      14,
			LetterSpacing: 0.25,
			FontWeight:    weight.Normal,
			FontFamily:    fp.Font(font.Plain),
		},
		typescale.BodySmall: {
			LineHeight:    16,
			FontSize:      12,
			LetterSpacing: 0.25,
			FontWeight:    weight.Normal,
			FontFamily:    fp.Font(font.Plain),
		},
	}
}

var defaultFonts = Fonts{
	font.Brand: []string{
		"Roboto",
		`"Noto Sans"`,
		`"Helvetica Neue"`,
		"sans-serif",
		`emoji`,
	},
	font.Plain: []string{
		"Roboto",
		`"Noto Sans"`,
		`"Helvetica Neue"`,
		"sans-serif",
		`emoji`,
	},
	font.Code: []string{
		`"Lucida Console"`,
		"Consolas",
		"Monaco",
		`"Andale Mono"`,
		`"Ubuntu Mono"`,
		`monospace`,
	},
}
