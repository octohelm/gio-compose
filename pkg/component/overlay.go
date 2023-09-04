package component

import (
	"context"
	"fmt"
	"strings"

	"gioui.org/io/key"
	. "github.com/octohelm/gio-compose/pkg/compose"
	"github.com/octohelm/gio-compose/pkg/compose/renderer"
	"github.com/octohelm/gio-compose/pkg/event/gesture"
	"github.com/octohelm/gio-compose/pkg/event/textinput"
	"github.com/octohelm/gio-compose/pkg/iter"
	"github.com/octohelm/gio-compose/pkg/layout"
	"github.com/octohelm/gio-compose/pkg/unit"
	"github.com/octohelm/gio-compose/pkg/util/contextutil"
)

type OverlayState interface {
	OverlayStack
	Visible() bool
	ElementRefs() []*Ref[Element]

	Topmost() bool

	IsClickInside(pt unit.Point) bool
}

type OverlayStack interface {
	Add(o OverlayState)
	Remove(o OverlayState)
}

var OverlayContext = contextutil.New[OverlayState](contextutil.Defaulter(func() OverlayState {
	return &overlayState{}
}))

type overlayState struct {
	visible     func() bool
	elementRefs func() []*Ref[Element]
	children    []OverlayState
}

func (o *overlayState) ElementRefs() []*Ref[Element] {
	return o.elementRefs()
}

func (o *overlayState) String() string {
	s := &strings.Builder{}
	_, _ = fmt.Fprintf(s, "overlay(%p", o)
	for _, c := range o.children {
		_, _ = fmt.Fprintf(s, ",%p", c)
	}
	_, _ = fmt.Fprintf(s, ")")
	return s.String()
}

func (o *overlayState) Add(child OverlayState) {
	o.children = append(o.children, child)
}

func (o *overlayState) Remove(child OverlayState) {
	o.children = iter.Filter(o.children, func(c OverlayState) bool {
		return c != child
	})
}

func (o *overlayState) Visible() bool {
	if o.visible != nil {
		return o.visible()
	}
	return false
}

func (o *overlayState) Topmost() bool {
	for _, c := range o.children {
		if c.Visible() {
			return false
		}
	}
	return true
}

func (o *overlayState) IsClickInside(pt unit.Point) bool {
	for _, child := range o.children {
		if child.IsClickInside(pt) {
			return true
		}
	}

	for _, r := range o.elementRefs() {
		rect := layout.GetBoundingClientRect(r.Current)

		if rect.Contains(pt) {
			return true
		}
	}

	return false
}

type Overlay struct {
	Visible         bool
	OnVisibleChange func(v bool)
	Refs            []*Ref[Element]
}

func (o Overlay) Build(c BuildContext) VNode {
	ref := UseRef(c, o.Refs)
	ref.Current = o.Refs

	visible := UseState(c, o.Visible)

	UseEffect(c, func() func() {
		if v := visible.Value(); v != o.Visible {
			if o.OnVisibleChange != nil {
				o.OnVisibleChange(v)
			}
		}
		return nil
	}, []any{visible.Value()})

	UseEffect(c, func() func() {
		if o.Visible != visible.Value() {
			visible.Update(o.Visible)
		}
		return nil
	}, []any{o.Visible})

	parentState := OverlayContext.Extract(c)

	state := UseMemo(c, func() *overlayState {
		return &overlayState{
			elementRefs: func() []*Ref[Element] {
				return ref.Current
			},
			visible: func() bool {
				return visible.Value()
			},
		}
	}, []any{})

	UseEffect(c, func() func() {
		parentState.Add(state)

		return func() {
			parentState.Remove(state)
		}
	}, []any{})

	// Key Escape to close topmost overlay
	UseEffect(c, func() func() {
		w := renderer.WindowContext.Extract(c)

		return textinput.WatchInputEvent(w, textinput.OnKeyUp(key.NameEscape, func() {
			if state.Topmost() {
				visible.UpdateFunc(func(prev bool) bool {
					return false
				})
			}
		}))
	}, []any{})

	// click outside
	UseEffect(c, func() func() {
		w := renderer.WindowContext.Extract(c)

		return gesture.WatchGestureEvent(w, gesture.OnTapWithEvent(func(e *gesture.PointerData) {
			if !state.IsClickInside(e.Position) {
				visible.UpdateFunc(func(prev bool) bool {
					return false
				})
			}
		}))
	}, []any{})

	if !visible.Value() {
		return nil
	}

	return RootPortal().Children(
		Provider(func(ctx context.Context) context.Context {
			return OverlayContext.Inject(ctx, state)
		}).Children(
			c.ChildVNodes()...,
		),
	)
}
