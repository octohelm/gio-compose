package node

import (
	"context"
	"fmt"
	"io"
	"strings"
)

func WithPtr(ptr bool) RenderOptionFunc {
	return func(w *printer) {
		w.ptr = ptr
	}
}

func WithPretty(pretty bool) RenderOptionFunc {
	return func(w *printer) {
		w.pretty = pretty
	}
}

type RenderOptionFunc = func(w *printer)

func RenderToString(w io.Writer, n Node, optFns ...RenderOptionFunc) {
	p := &printer{Writer: w}

	for _, optFn := range optFns {
		optFn(p)
	}

	p.PrintTo(w, n, 0)
}

type printer struct {
	io.Writer

	ptr    bool
	pretty bool
}

func (p *printer) PrintTo(w io.Writer, n Node, depth int) {
	if n == nil {
		return
	}

	p.OpenTag(n, depth)
	defer p.CloseTag(n, depth)

	for c := range IterChildElement(context.Background(), n) {
		p.PrintTo(w, c, depth+1)
	}
}

func (p *printer) TagName(n Node) string {
	if p.ptr {
		return fmt.Sprintf("%s.%p", n.DisplayName(), n)
	}
	return n.DisplayName()
}

func (p *printer) OpenTag(n Node, depth int) {
	isFrag := IsFragment(n)
	if isFrag {
		return
	}

	noChild := n.FirstChild() == nil

	if p.pretty && depth > 0 {
		_, _ = fmt.Fprintf(p, "\n")
		_, _ = fmt.Fprintf(p, strings.Repeat("  ", depth))
	}

	if noChild {
		_, _ = fmt.Fprintf(p, "<%s", p.TagName(n))
		return
	}

	_, _ = fmt.Fprintf(p, "<%s>", p.TagName(n))
}

func (p *printer) CloseTag(n Node, depth int) {
	isFrag := IsFragment(n)
	if isFrag {
		return
	}

	noChild := n.FirstChild() == nil

	if noChild {
		_, _ = io.WriteString(p, "/>")
		return
	}

	if p.pretty {
		_, _ = fmt.Fprintf(p, "\n")
		if depth > 0 {
			_, _ = fmt.Fprintf(p, strings.Repeat("  ", depth))
		}
	}

	_, _ = fmt.Fprintf(p, "</%s>", p.TagName(n))
}
