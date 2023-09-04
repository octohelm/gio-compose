package bezier

type EasingFunc = func(x float64) float64

func Easing(x0, y0, x1, y1 float64) EasingFunc {
	if !(0 <= x0 && x0 <= 1 && 0 <= x1 && x1 <= 1) {
		panic("bezier x values must be in [0, 1] range")
	}

	if x0 == y0 && x1 == y1 {
		return func(x float64) float64 {
			return x
		}
	}

	return func(x float64) float64 {
		if x == 0 || x == 1 {
			return x
		}

		t := x
		for i := 0; i < 5; i++ {
			t2 := t * t
			t3 := t2 * t
			d := 1 - t
			d2 := d * d

			nx := 3*d2*t*x0 + 3*d*t2*x1 + t3
			dxdt := 3*d2*x0 + 6*d*t*(x1-x0) + 3*t2*(1-x1)
			if dxdt == 0 {
				break
			}

			t -= (nx - x) / dxdt
			if t <= 0 || t >= 1 {
				break
			}
		}
		if t < 0 {
			t = 0
		}
		if t > 1 {
			t = 1
		}

		// Solve for y using t.
		t2 := t * t
		t3 := t2 * t
		d := 1 - t
		d2 := d * d
		y := 3*d2*t*y0 + 3*d*t2*y1 + t3

		return y
	}
}
