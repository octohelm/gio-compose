package layout

import (
	"github.com/octohelm/gio-compose/pkg/cmp"
	"github.com/octohelm/gio-compose/pkg/layout/position"
)

type PositionSetter interface {
	SetPosition(position position.Position)
}

type Positioner struct {
	Position position.Position
}

func (p *Positioner) Eq(v *Positioner) cmp.Result {
	return func() bool {
		return p.Position == v.Position
	}
}

var _ PositionSetter = &Positioner{}

func (p *Positioner) SetPosition(position position.Position) {
	p.Position = position
}
