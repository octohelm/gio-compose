package layout

import (
	"github.com/octohelm/gio-compose/pkg/cmp"
	"github.com/octohelm/gio-compose/pkg/layout/alignment"
)

type Alignment = alignment.Alignment

type AlignSetter interface {
	SetAlign(align Alignment)
}

type AlignGetter interface {
	Align() Alignment
}

type Aligner struct {
	Alignment Alignment
}

func (r *Aligner) Eq(v *Aligner) cmp.Result {
	return func() bool {
		return r.Alignment == v.Alignment
	}
}

var _ AlignSetter = &Aligner{}

func (a *Aligner) SetAlign(alignment Alignment) {
	a.Alignment = alignment
}

func (a *Aligner) Align() Alignment {
	return a.Alignment
}
