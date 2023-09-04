package text

import (
	"image/color"

	"github.com/octohelm/gio-compose/pkg/text/weight"

	"github.com/octohelm/gio-compose/pkg/cmp"
	"github.com/octohelm/x/ptr"

	"gioui.org/font"
	"gioui.org/text"
	"gioui.org/unit"
)

type Alignment int

const (
	Start Alignment = iota
	End
	Middle
)

func (a Alignment) TextAlignment() text.Alignment {
	return text.Alignment(a)
}

type StyleSetter interface {
	SetStyle(s Style)
}

type FontFamilySetter interface {
	SetFontFamily(fontFamily string)
}

type FontWeightSetter interface {
	SetFontWeight(w weight.Weight)
}

type FontStyleSetter interface {
	SetFontStyle(style font.Style)
}

type FontSizeSetter interface {
	SetFontSize(size unit.Sp)
}

type LineHeightSetter interface {
	SetLineHeight(size unit.Sp)
}

type TextColorSetter interface {
	SetTextColor(color color.Color)
}

type TextAlignSetter interface {
	SetTextAlign(a Alignment)
}

type Style struct {
	FontFamily string
	FontWeight weight.Weight
	FontStyle  *font.Style

	FontSize   unit.Sp
	LineHeight unit.Sp
	// Not Implemented
	LetterSpacing unit.Sp

	Color     color.NRGBA
	TextAlign *Alignment
}

var _ FontFamilySetter = &Style{}
var _ FontWeightSetter = &Style{}
var _ FontStyleSetter = &Style{}
var _ FontSizeSetter = &Style{}
var _ LineHeightSetter = &Style{}
var _ TextColorSetter = &Style{}
var _ TextAlignSetter = &Style{}
var _ StyleSetter = &Style{}

func (s *Style) SetFontFamily(fontFamily string) {
	s.FontFamily = fontFamily
}

func (s *Style) SetFontWeight(w weight.Weight) {
	s.FontWeight = w
}

func (s *Style) SetTextAlign(textAlign Alignment) {
	s.TextAlign = ptr.Ptr(textAlign)
}

func (s *Style) SetTextColor(c color.Color) {
	if c != nil {
		s.Color = color.NRGBAModel.Convert(c).(color.NRGBA)
	}
}

func (s *Style) SetFontSize(fontSize unit.Sp) {
	s.FontSize = fontSize
}

func (s *Style) SetLineHeight(lineHeight unit.Sp) {
	s.LineHeight = lineHeight
}

func (s *Style) SetFontStyle(style font.Style) {
	s.FontStyle = ptr.Ptr(style)
}

var emptyColor = color.NRGBA{}

func (s *Style) SetStyle(style Style) {
	*s = s.Merge(style)
}

func (s Style) Merge(style Style) Style {
	if style.FontFamily != "" {
		s.FontFamily = style.FontFamily
	}
	if style.FontWeight != 0 {
		s.FontWeight = style.FontWeight
	}
	if style.FontStyle != nil {
		s.FontStyle = style.FontStyle
	}

	if style.FontSize != 0 {
		s.FontSize = style.FontSize
	}

	if style.LineHeight != 0 {
		s.LineHeight = style.LineHeight
	}

	if style.Color != emptyColor {
		s.Color = style.Color
	}

	if style.TextAlign != nil {
		s.TextAlign = style.TextAlign
	}

	return s
}

func (s *Style) Eq(style *Style) cmp.Result {
	return cmp.All(
		cmp.Eq(s.FontFamily, style.FontFamily),
		cmp.Eq(s.FontWeight, style.FontWeight),
		cmp.Eq(s.FontSize, style.FontSize),
		cmp.Eq(s.LineHeight, style.LineHeight),
		cmp.Eq(s.Color, style.Color),
		func() bool {
			if s.FontStyle == nil || style.FontStyle == nil {
				return false
			}
			return *s.FontStyle == *style.FontStyle
		},
		func() bool {
			if s.TextAlign == nil || style.TextAlign == nil {
				return false
			}
			return *s.TextAlign == *style.TextAlign
		},
	)
}

func (s Style) ToFont() font.Font {
	f := font.Font{
		Typeface: font.Typeface(s.FontFamily),
		Weight:   s.FontWeight.GioWeight(),
	}

	if style := s.FontStyle; style != nil {
		f.Style = *style
	}

	return f
}
