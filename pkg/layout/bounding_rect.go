package layout

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/octohelm/gio-compose/pkg/node"
	"github.com/octohelm/gio-compose/pkg/unit"
)

type BoundingRectGetter interface {
	BoundingRect() BoundingRect
}

type BoundingRect struct {
	X, Y          unit.Dp
	Width, Height unit.Dp
}

func (r BoundingRect) Contains(pt unit.Point) bool {
	return pt.X >= r.X && pt.Y >= r.Y && pt.X <= (r.X+r.Width) && pt.Y <= (r.Y+r.Height)
}

func (r BoundingRect) String() string {
	return fmt.Sprintf("x=%v,y=%v,w=%v,h=%v", r.X, r.Y, r.Width, r.Height)
}

func GetBoundingRect(n node.Node) (boundingRect BoundingRect) {
	if n == nil {
		return
	}

	if getter, ok := n.(BoundingRectGetter); ok {
		boundingRect = getter.BoundingRect()
	}

	return
}

func GetBoundingClientRect(n node.Node) (boundingRect BoundingRect) {
	if n == nil {
		return
	}

	if getter, ok := n.(BoundingRectGetter); ok {
		boundingRect = getter.BoundingRect()
	}

	for p := n.ParentNode(); p != nil; p = p.ParentNode() {
		if node.IsFragment(p) {
			continue
		}

		if getter, ok := p.(BoundingRectGetter); ok {
			parentRect := getter.BoundingRect()
			boundingRect.X += parentRect.X
			boundingRect.Y += parentRect.Y
		}
	}

	return
}

func PrintBoundingRectTo(w io.Writer, n node.Node) {
	printBoundingRectTo(w, n, 0)
}

func printBoundingRectTo(w io.Writer, n node.Node, depth int) {
	if !(node.IsFragment(n)) {
		_, _ = fmt.Fprintf(w, "%s%s ", strings.Repeat("  ", depth), n.DisplayName())

		if g, ok := n.(BoundingRectGetter); ok {
			_, _ = fmt.Fprintf(w, "%s", g.BoundingRect())
		}

		_, _ = io.WriteString(w, "\n")

		for c := range node.IterChildElement(context.Background(), n) {
			printBoundingRectTo(w, c.(node.Node), depth+1)
		}
	} else {
		for c := range node.IterChildElement(context.Background(), n) {
			printBoundingRectTo(w, c.(node.Node), depth)
		}
	}

}
