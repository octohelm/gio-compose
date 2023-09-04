package cmp

type Result = func() bool

func Cast[T any](v any, do func(v *T) Result) Result {
	switch x := v.(type) {
	case T:
		return do(&x)
	case *T:
		return do(x)
	}
	return func() bool {
		return false
	}
}

func UpdateWhen[T any](when Result, target *T, next *T) bool {
	if !when() {
		return false
	}
	*target = *next
	return true
}

func All(results ...Result) Result {
	return func() bool {
		for _, r := range results {
			if r == nil {
				continue
			}
			if !r() {
				return false
			}
		}
		return true
	}
}

func Any(results ...Result) Result {
	return func() bool {
		for _, r := range results {
			if r == nil {
				continue
			}
			if r() {
				return true
			}
		}
		return false
	}
}

func Not(result Result) Result {
	return func() bool {
		return !result()
	}
}

func Eq[T comparable](a T, b T) Result {
	return func() bool {
		return a == b
	}
}
