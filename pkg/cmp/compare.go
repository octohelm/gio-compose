package cmp

import (
	"reflect"
)

func ShallowEqual(a any, b any) bool {
	switch xa := a.(type) {
	case Comparer:
		return xa.Compare(b) == 0
	case map[string]any:
		switch xb := b.(type) {
		case map[string]any:
			if len(xa) != len(xb) {
				return false
			}

			for k := range xa {
				if xbv, ok := xb[k]; ok {
					if xa[k] != xbv {
						return false
					}
				} else {
					return false
				}
			}

			return true
		}
		return false
	case []any:
		switch xb := b.(type) {
		case []any:
			if len(xa) != len(xb) {
				return false
			}

			for i := range xa {
				if !ShallowEqual(xa[i], xb[i]) {
					return false
				}
			}
			return true
		}
		return false
	}

	if a == nil || b == nil {
		return false
	}

	if reflect.TypeOf(a).Comparable() && reflect.TypeOf(b).Comparable() {
		return a == b
	}

	return false
}

type Comparer interface {
	Compare(v any) int
}
