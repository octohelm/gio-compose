package modifier

import (
	"github.com/octohelm/gio-compose/pkg/layout"
	"github.com/octohelm/gio-compose/pkg/modifier"
)

func Weight(weight float32) modifier.Modifier[any] {
	return &weightModifier{weight: weight}
}

type weightModifier struct {
	weight float32
}

func (m weightModifier) Modify(w any) {
	if s, ok := w.(layout.WeightSetter); ok {
		s.SetWeight(m.weight)
	}
}
