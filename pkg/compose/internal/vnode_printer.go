package internal

import (
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

func RenderToString(w io.Writer, n VNodeAccessor, optFns ...RenderOptionFunc) {
	p := &printer{
		Writer: w,
		ptr:    true,
		pretty: true,
	}

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

func (p *printer) PrintTo(w io.Writer, n VNodeAccessor, depth int) {
	if n == nil {
		return
	}

	p.OpenTag(n, depth)
	defer p.CloseTag(n, depth)

	for _, c := range n.ChildVNodeAccessors() {
		p.PrintTo(w, c, depth+1)
	}
}

func (p *printer) TagName(n VNodeAccessor) string {
	if p.ptr {
		return fmt.Sprintf("%s.%p", n.String(), n)
	}
	return n.String()
}

func (p *printer) OpenTag(n VNodeAccessor, depth int) {
	if p.pretty {
		_, _ = fmt.Fprintf(p, "\n")

		if depth > 0 {
			_, _ = fmt.Fprintf(p, strings.Repeat("  ", depth))
		}
	}

	if len(n.ChildVNodeAccessors()) == 0 {
		_, _ = fmt.Fprintf(p, "<%s", p.TagName(n))
		return
	}
	_, _ = fmt.Fprintf(p, "<%s>", p.TagName(n))
}

func (p *printer) CloseTag(n VNodeAccessor, depth int) {
	if len(n.ChildVNodeAccessors()) == 0 {
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
