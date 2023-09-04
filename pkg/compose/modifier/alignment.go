package modifier

import (
	"fmt"

	"github.com/octohelm/gio-compose/pkg/layout"
	"github.com/octohelm/gio-compose/pkg/layout/alignment"
	"github.com/octohelm/gio-compose/pkg/modifier"
)

func Align(alignment alignment.Alignment) modifier.Modifier[any] {
	return &alignmentModifier{
		alignment: alignment,
	}
}

type alignmentModifier struct {
	alignment alignment.Alignment
}

func (m *alignmentModifier) String() string {
	return fmt.Sprintf("[Alignment] = %v", m.alignment)
}

func (m *alignmentModifier) Modify(w any) {
	if setter, ok := w.(layout.AlignSetter); ok {
		setter.SetAlign(m.alignment)
	}
}
