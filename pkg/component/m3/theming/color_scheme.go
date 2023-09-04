package theming

import (
	"fmt"
	"image/color"

	"goki.dev/cam/hct"

	"github.com/octohelm/gio-compose/pkg/component/m3/colorrole"
	"github.com/octohelm/gio-compose/pkg/paint/f32color"
)

var DefaultColorScheme = ColorSchemeFrom(f32color.RGB(0x6750A4))

func AsDarkMode() ColorSchemeOptionFunc {
	return func(p *palette) {
		p.dark = true
	}
}

func WithSecondaryKeyColor(keyColor color.Color) ColorSchemeOptionFunc {
	return func(p *palette) {
		cc := hct.FromColor(keyColor)
		p.a2 = TonalPaletteFromHueAndChroma(cc.Hue, cc.Chroma)
	}
}

func WithTertiaryKeyColor(keyColor color.Color) ColorSchemeOptionFunc {
	return func(p *palette) {
		cc := hct.FromColor(keyColor)
		p.a3 = TonalPaletteFromHueAndChroma(cc.Hue, cc.Chroma)
	}
}

func WithNeutralKeyColor(keyColor color.Color) ColorSchemeOptionFunc {
	return func(p *palette) {
		cc := hct.FromColor(keyColor)
		p.n1 = TonalPaletteFromHueAndChroma(cc.Hue, cc.Chroma)
	}
}

func WithNeutralVariantKeyColor(keyColor color.Color) ColorSchemeOptionFunc {
	return func(p *palette) {
		cc := hct.FromColor(keyColor)
		p.n2 = TonalPaletteFromHueAndChroma(cc.Hue, cc.Chroma)
	}
}

func WithCustomKeyColor(role colorrole.ColorRole, keyColor color.Color) ColorSchemeOptionFunc {
	return func(p *palette) {
		cc := hct.FromColor(keyColor)

		if role == "error" {
			p.error = TonalPaletteFromHueAndChroma(cc.Hue, cc.Chroma)
		} else {
			if p.custom == nil {
				p.custom = map[colorrole.ColorRole]TonalPalette{}
			}
			p.custom[role] = TonalPaletteFromHueAndChroma(cc.Hue, cc.Chroma)
		}
	}
}

func ColorSchemeFrom(primaryKeyColor color.Color, optionFns ...ColorSchemeOptionFunc) ColorScheme {
	p := &palette{}
	for i := range optionFns {
		optionFns[i](p)
	}

	cc := hct.FromColor(primaryKeyColor)

	if p.a1 == nil {
		p.a1 = TonalPaletteFromHueAndChroma(cc.Hue, max(48, cc.Chroma))
	}
	if p.a2 == nil {
		p.a2 = TonalPaletteFromHueAndChroma(cc.Hue, 16)
	}
	if p.a3 == nil {
		p.a3 = TonalPaletteFromHueAndChroma(cc.Hue+60, 24)
	}
	if p.n1 == nil {
		p.n1 = TonalPaletteFromHueAndChroma(cc.Hue, 4)
	}
	if p.n2 == nil {
		p.n2 = TonalPaletteFromHueAndChroma(cc.Hue, 8)
	}
	if p.error == nil {
		p.error = TonalPaletteFromHueAndChroma(25, 84)
	}

	rules := normalizeColorRules(p)

	m := colorScheme{}

	if p.dark {
		for k := range rules {
			m[k] = rules[k].colorForDark(p)
		}
		return m
	}

	for k := range rules {
		m[k] = rules[k].colorForLight(p)
	}
	return m
}

func normalizeColorRules(p *palette) map[colorrole.ColorRole]colorRule {
	rules := map[colorrole.ColorRole]colorRule{
		"outline":         {neutralVariant, 60, 50},
		"outline-variant": {neutralVariant, 30, 80},

		"surface":    {neutral, 10, 99},
		"on-surface": {neutral, 90, 10},

		"surface-variant":    {neutralVariant, 30, 90},
		"on-surface-variant": {neutralVariant, 80, 30},

		"surface-dim":    {neutral, 6, 87},
		"surface-bright": {neutral, 24, 98},

		"surface-container-lowest":  {neutral, 4, 100},
		"surface-container-low":     {neutral, 10, 96},
		"surface-container":         {neutral, 12, 94},
		"surface-container-high":    {neutral, 17, 92},
		"surface-container-highest": {neutral, 22, 90},

		"inverse-surface":    {neutral, 90, 20},
		"inverse-on-surface": {neutral, 20, 95},
		"inverse-primary":    {primary, 40, 80},
	}

	seedColors := map[colorrole.ColorRole]TonalPalette{
		colorrole.Primary:   p.a1,
		colorrole.Secondary: p.a2,
		colorrole.Tertiary:  p.a3,
		colorrole.Error:     p.error,
	}

	for k := range p.custom {
		seedColors[k] = p.custom[k]
	}

	// https://m3.material.io/styles/color/the-color-system/custom-colors
	for k := range seedColors {
		p2 := seedColors[k]
		keyColor := func(p *palette) TonalPalette { return p2 }

		rules[k] = colorRule{keyColor, 80, 40}
		rules[colorrole.ColorRole(fmt.Sprintf("on-%s", k))] = colorRule{keyColor, 20, 100}
		rules[colorrole.ColorRole(fmt.Sprintf("%s-container", k))] = colorRule{keyColor, 30, 90}
		rules[colorrole.ColorRole(fmt.Sprintf("on-%s-container", k))] = colorRule{keyColor, 90, 10}
	}

	return rules
}

type colorRule struct {
	keyColor     func(p *palette) TonalPalette
	toneForDark  float32
	toneForLight float32
}

func (c colorRule) colorForLight(p *palette) color.NRGBA {
	cc := c.keyColor(p).Toned(c.toneForLight)
	return color.NRGBAModel.Convert(cc).(color.NRGBA)
}

func (c colorRule) colorForDark(p *palette) color.NRGBA {
	cc := c.keyColor(p).Toned(c.toneForDark)
	return color.NRGBAModel.Convert(cc).(color.NRGBA)
}

type ColorSchemeOptionFunc func(p *palette)

func primary(p *palette) TonalPalette {
	return p.a1
}
func neutralVariant(p *palette) TonalPalette {
	return p.n2
}

func neutral(p *palette) TonalPalette {
	return p.n1
}

type palette struct {
	a1     TonalPalette
	a2     TonalPalette
	a3     TonalPalette
	n1     TonalPalette
	n2     TonalPalette
	error  TonalPalette
	custom map[colorrole.ColorRole]TonalPalette
	dark   bool
}

type colorScheme map[colorrole.ColorRole]color.NRGBA

func (m colorScheme) Color(role colorrole.ColorRole) color.NRGBA {
	return m[role]
}

type ColorScheme interface {
	Color(role colorrole.ColorRole) color.NRGBA
}

func TonalPaletteFromHueAndChroma(hue, chroma float32) TonalPalette {
	return &tonalPalette{
		Hue:    hue,
		Chroma: chroma,
	}
}

type TonalPalette interface {
	Toned(tone float32) color.Color
}

type tonalPalette struct {
	Hue    float32
	Chroma float32
}

func (p *tonalPalette) Toned(tone float32) color.Color {
	return f32color.HCT{
		Hue:    p.Hue,
		Chroma: p.Chroma,
		Tone:   tone,
	}
}
