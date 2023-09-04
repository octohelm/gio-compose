package layout

import "github.com/octohelm/gio-compose/pkg/cmp"

type WeightSetter interface {
	SetWeight(weight float32)
}

type WeightGetter interface {
	Weight() (float32, bool)
}

type FlexWeight struct {
	weight float32
}

var _ WeightGetter = &FlexWeight{}

func (r *FlexWeight) Eq(v *FlexWeight) cmp.Result {
	return func() bool {
		return r.weight == v.weight
	}
}

func (w *FlexWeight) Weight() (float32, bool) {
	return w.weight, w.weight > 0
}

func (w *FlexWeight) SetWeight(weight float32) {
	w.weight = weight
}
