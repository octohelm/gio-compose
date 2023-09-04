package cmp

import (
	"image/color"
	"reflect"
	"testing"

	testingx "github.com/octohelm/x/testing"
)

func TestUpdateWhen(t *testing.T) {
	t.Run("Do nothing when not changed", func(t *testing.T) {
		f1 := &Fill{Color: color.NRGBA{A: 0x01}}
		f2 := &Fill{Color: color.NRGBA{A: 0x01}}

		updated := UpdateWhen(Not(Eq(f1.Color, f2.Color)), f1, f2)

		testingx.Expect(t, updated, testingx.Be(false))
	})

	t.Run("Do nothing when not changed", func(t *testing.T) {
		f1 := &Fill{Color: color.NRGBA{A: 0x01}}
		f2 := &Fill{Color: color.NRGBA{A: 0x02}}

		updated := UpdateWhen(Not(Eq(f1.Color, f2.Color)), f1, f2)

		testingx.Expect(t, updated, testingx.Be(true))
		testingx.Expect(t, f1, testingx.Equal(f2))
	})
}

func BenchmarkEqual(b *testing.B) {
	f1 := Border{Fill: Fill{Color: color.NRGBA{A: 0x01}}}
	f2 := Border{Fill: Fill{Color: color.NRGBA{A: 0x01}}}

	b.Run("direct equal", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = f1 == f2
		}
	})

	b.Run("Result interface", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = f1.Eq(&f2)()
		}
	})

	b.Run("reflect.DeepEqual", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = reflect.DeepEqual(f1, f2)
		}
	})
}

type Fill struct {
	Color color.NRGBA
}

func (f *Fill) Eq(v *Fill) Result {
	return Eq(f.Color, v.Color)
}

type Border struct {
	Width float32
	Fill
}

func (f *Border) Eq(v *Border) Result {
	return All(
		Eq(f.Width, v.Width),
		f.Fill.Eq(&v.Fill),
	)
}
