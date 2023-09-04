package modifier

import (
	"fmt"

	"github.com/octohelm/gio-compose/pkg/layout"
	"github.com/octohelm/gio-compose/pkg/modifier"
	"github.com/octohelm/gio-compose/pkg/unit"
)

func Gap(spacing unit.Dp) modifier.Modifier[any] {
	return &gapModifier{spacing: spacing}
}

type gapModifier struct {
	spacing unit.Dp
}

func (m *gapModifier) String() string {
	return fmt.Sprintf("[Spacing] = %v", m.spacing)
}

func (m *gapModifier) Modify(w any) {
	if s, ok := w.(layout.SpacingSetter); ok {
		s.SetSpacing(m.spacing)
	}
}
