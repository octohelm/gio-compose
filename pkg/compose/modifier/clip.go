package modifier

import (
	"github.com/octohelm/gio-compose/pkg/layout"
	"github.com/octohelm/gio-compose/pkg/modifier"
)

func Clip() modifier.Modifier[any] {
	return &clipModifier{}
}

type clipModifier struct{}

func (c *clipModifier) Modify(target any) {
	if s, ok := target.(layout.ClipSetter); ok {
		s.SetClip(true)
	}
}
