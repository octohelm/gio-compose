package text

import (
	"gioui.org/font/gofont"
	"gioui.org/text"

	"github.com/octohelm/gio-compose/pkg/util/contextutil"
)

var DefaultShaper = text.NewShaper(text.WithCollection(gofont.Collection()))

var ShaperContext = contextutil.New[*text.Shaper](contextutil.Defaulter(func() *text.Shaper {
	return DefaultShaper
}))
