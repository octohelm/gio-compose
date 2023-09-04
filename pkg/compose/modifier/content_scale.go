package modifier

import (
	"github.com/octohelm/gio-compose/pkg/modifier"
	"github.com/octohelm/gio-compose/pkg/paint"
	"github.com/octohelm/gio-compose/pkg/paint/contentscale"
)

func ContentScale(scale contentscale.ContentScale) modifier.Modifier[any] {
	return &contentScaleModifier{
		scale: scale,
	}
}

type contentScaleModifier struct {
	scale contentscale.ContentScale
}

func (c *contentScaleModifier) Modify(target any) {
	if s, ok := target.(paint.ContentScaleSetter); ok {
		s.SetContentScale(c.scale)
	}
}
