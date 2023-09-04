package layout

import (
	"github.com/octohelm/gio-compose/pkg/cmp"
	"github.com/octohelm/gio-compose/pkg/layout/arrangement"
)

type Arrangement = arrangement.Arrangement

type ArrangementSetter interface {
	SetArrangement(Arrangement Arrangement)
}

type Arrangementer struct {
	Arrangement Arrangement
}

func (r *Arrangementer) Eq(v *Arrangementer) cmp.Result {
	return func() bool {
		return r.Arrangement == v.Arrangement
	}
}

var _ ArrangementSetter = &Arrangementer{}

func (a *Arrangementer) SetArrangement(arrangement Arrangement) {
	a.Arrangement = arrangement
}
