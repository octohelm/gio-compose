package layout

import (
	"github.com/octohelm/gio-compose/pkg/cmp"
	"github.com/octohelm/gio-compose/pkg/layout/direction"
)

type ScrollableSetter interface {
	SetScrollable(axis direction.Axis, enabled bool)
}

var _ ScrollableSetter = &Scrollable{}

type Scrollable struct {
	Enabled bool
	Axis    direction.Axis
}

func (s *Scrollable) Eq(v *Scrollable) cmp.Result {
	return func() bool {
		return *s == *v
	}
}

func (s *Scrollable) SetScrollable(axis direction.Axis, enabled bool) {
	s.Axis = axis
	s.Enabled = enabled
}
