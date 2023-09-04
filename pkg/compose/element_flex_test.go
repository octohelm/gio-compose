package compose_test

import (
	"image"
	"testing"

	"github.com/octohelm/gio-compose/pkg/layout/alignment"

	"github.com/octohelm/gio-compose/pkg/layout/arrangement"

	"github.com/octohelm/gio-compose/pkg/compose/modifier"

	. "github.com/octohelm/gio-compose/pkg/compose"
)

func TestFlex(t *testing.T) {
	t.Run("Row", func(t *testing.T) {
		t.Run("fit children", func(t *testing.T) {
			el := Row(modifier.FillMaxSize()).Children(
				Row().Children(
					Box(modifier.Size(100)),
					Box(modifier.Size(100)),
				),
			)

			ExpectLayout(t, el, image.Pt(1000, 1000), `
Row x=0,y=0,w=1000,h=1000
  Row x=0,y=0,w=200,h=100
    Box x=0,y=0,w=100,h=100
    Box x=100,y=0,w=100,h=100
`)
		})

		t.Run("flexed only", func(t *testing.T) {
			el := Row(modifier.FillMaxSize(), modifier.Arrangement(arrangement.EqualWeight)).Children(
				Box(modifier.FillMaxHeight()),
				Box(modifier.FillMaxHeight()),
			)

			ExpectLayout(t, el, image.Pt(1000, 1000), `
Row x=0,y=0,w=1000,h=1000
  Box x=0,y=0,w=500,h=1000
  Box x=500,y=0,w=500,h=1000
`)

		})

		t.Run("center text child", func(t *testing.T) {
			el := Row(modifier.FillMaxSize(), modifier.Align(alignment.Center)).Children(
				Text("Hi 你好"),
			)

			ExpectLayout(t, el, image.Pt(1000, 1000), `
Row x=0,y=0,w=1000,h=1000
  Text x=0,y=492,w=1000,h=16
`)
		})

		t.Run("with sized child", func(t *testing.T) {
			el := Row(modifier.FillMaxSize()).Children(
				Box(modifier.Width(200), modifier.FillMaxHeight()),
				Box(modifier.Weight(1), modifier.FillMaxHeight(), modifier.PaddingAll(50)).Children(
					Box(modifier.FillMaxSize()),
				),
			)

			ExpectLayout(t, el, image.Pt(1000, 1000), `
Row x=0,y=0,w=1000,h=1000
  Box x=0,y=0,w=200,h=1000
  Box x=200,y=0,w=800,h=1000
    Box x=50,y=50,w=700,h=900
`)
		})

		t.Run("with sized child and gap", func(t *testing.T) {
			el := Row(modifier.FillMaxSize(), modifier.Gap(10)).Children(
				Box(modifier.Width(200), modifier.FillMaxHeight()),
				Box(modifier.Weight(1), modifier.FillMaxHeight(), modifier.PaddingAll(50)).Children(
					Box(modifier.FillMaxSize()),
				),
			)

			ExpectLayout(t, el, image.Pt(1000, 1000), `
Row x=0,y=0,w=1000,h=1000
  Box x=0,y=0,w=200,h=1000
  Box x=210,y=0,w=790,h=1000
    Box x=50,y=50,w=690,h=900
`)
		})

		t.Run("with sized child and different weight", func(t *testing.T) {
			el := Row(modifier.FillMaxSize()).Children(
				Box(modifier.Width(200), modifier.FillMaxHeight()),
				Box(modifier.Weight(1), modifier.FillMaxHeight()),
				Box(modifier.Weight(3), modifier.FillMaxHeight()),
			)

			ExpectLayout(t, el, image.Pt(1000, 1000), `
Row x=0,y=0,w=1000,h=1000
  Box x=0,y=0,w=200,h=1000
  Box x=200,y=0,w=200,h=1000
  Box x=400,y=0,w=600,h=1000
`)
		})

		t.Run("with different weight children and sized child", func(t *testing.T) {
			el := Row(modifier.FillMaxSize()).Children(
				Box(modifier.Weight(0.5), modifier.FillMaxHeight()),
				Box(modifier.Weight(1.5), modifier.FillMaxHeight()),
				Box(modifier.Width(200), modifier.FillMaxHeight()),
			)

			ExpectLayout(t, el, image.Pt(1000, 1000), `
Row x=0,y=0,w=1000,h=1000
  Box x=0,y=0,w=200,h=1000
  Box x=200,y=0,w=600,h=1000
  Box x=800,y=0,w=200,h=1000
`)
		})
	})

	t.Run("Column", func(t *testing.T) {
		t.Run("with different weight children and sized child", func(t *testing.T) {
			el := Column(modifier.FillMaxSize()).Children(
				Box(modifier.Height(200), modifier.FillMaxWidth()),
				Box(modifier.Weight(0.5), modifier.FillMaxWidth()),
				Box(modifier.Weight(1.5), modifier.FillMaxWidth()),
			)

			ExpectLayout(t, el, image.Pt(1000, 1000), `
Column x=0,y=0,w=1000,h=1000
  Box x=0,y=0,w=1000,h=200
  Box x=0,y=200,w=1000,h=200
  Box x=0,y=400,w=1000,h=600
`)
		})
	})
}
