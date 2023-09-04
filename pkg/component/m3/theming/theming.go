package theming

import "github.com/octohelm/gio-compose/pkg/util/contextutil"

type Theming interface {
	Typography
	ColorScheme
}

func WithColorScheme(colorScheme ColorScheme) func(opt *theming) {
	return func(opt *theming) {
		opt.ColorScheme = colorScheme
	}
}

func WithTypography(typography Typography) func(opt *theming) {
	return func(opt *theming) {
		opt.Typography = typography
	}
}

type OptionFunc = func(opt *theming)

func New(optFns ...OptionFunc) Theming {
	t := &theming{}
	for i := range optFns {
		optFns[i](t)
	}

	if t.ColorScheme == nil {
		t.ColorScheme = DefaultColorScheme
	}

	if t.Typography == nil {
		t.Typography = DefaultTypography
	}

	return t
}

type theming struct {
	Typography
	ColorScheme
}

var Context = contextutil.New[Theming]()
