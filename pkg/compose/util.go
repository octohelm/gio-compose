package compose

func Echo[O any](fn func() O) O {
	return fn()
}

func Echo2[O any, O1 any](fn func() (O, O1)) (O, O1) {
	return fn()
}

func Echo3[O any, O1 any, O2 any](fn func() (O, O1, O2)) (O, O1, O2) {
	return fn()
}
