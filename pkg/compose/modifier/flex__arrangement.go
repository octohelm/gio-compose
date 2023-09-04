package modifier

import (
	"fmt"

	"github.com/octohelm/gio-compose/pkg/layout"
	"github.com/octohelm/gio-compose/pkg/layout/arrangement"
	"github.com/octohelm/gio-compose/pkg/modifier"
)

func Arrangement(arrangement arrangement.Arrangement) modifier.Modifier[any] {
	return &arrangementModifier{
		arrangement: arrangement,
	}
}

type arrangementModifier struct {
	arrangement arrangement.Arrangement
}

func (m *arrangementModifier) String() string {
	return fmt.Sprintf("[arrangement] = %v", m.arrangement)
}

func (m *arrangementModifier) Modify(w any) {
	if setter, ok := w.(layout.ArrangementSetter); ok {
		setter.SetArrangement(m.arrangement)
	}
}
