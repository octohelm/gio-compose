package internal

import (
	"context"
	"fmt"
	"reflect"
	"slices"

	"github.com/octohelm/gio-compose/pkg/modifier"
	"github.com/octohelm/gio-compose/pkg/node"
)

func H(component Component, modifiers ...modifier.Modifier[any]) VNode {
	n := &vnode{}
	n.typ = component
	n.modifiers = make([]modifier.Modifier[any], 0, len(modifiers))

	if len(modifiers) > 0 {
		var wrappers []ComponentWrapper

		for m := range modifier.Iter(context.Background(), modifiers...) {
			switch x := modifier.Unwrap(m).(type) {
			case ComponentWrapper:
				wrappers = append(wrappers, x)
			case VNodeModifier:
				m.Modify(n)
			default:
				n.modifiers = append(n.modifiers, m)
			}
		}

		if len(wrappers) > 0 {
			var vn VNode = n
			for i := len(wrappers) - 1; i >= 0; i-- {
				vn = H(wrappers[i].Wrap(vn))
			}
			if n.key != nil {
				vn.(VNodeAccessor).SetKey(n.key)
			}
			return vn
		}
	}

	return n
}

type VNodeModifier interface {
	VNodeOnly()
}

type VNode interface {
	Children(vnodes ...VNode) VNode
}

type VNodeAccessor interface {
	String() string
	Key() any
	SetKey(key any)

	Type() Component
	Modifiers() modifier.Modifiers
	ChildVNodes() []VNode

	PutChildVNodes(childVNodes ...VNode)
	ReplaceChildVNodeAccessor(old VNodeAccessor, new VNodeAccessor)
	ChildVNodeAccessors() []VNodeAccessor

	Parent() VNodeAccessor
	IsRoot() bool
	Node() node.Node

	BindNode(n node.Node)
	BindParent(va VNodeAccessor)

	OnMount(fn func(ctx context.Context, va VNodeAccessor))
	OnUpdate(fn func(ctx context.Context, va VNodeAccessor, old VNodeAccessor))

	WillRender(oldVNode VNodeAccessor)
	Update(ctx context.Context)
	DidMount(ctx context.Context)

	Use(hook Hook) Hook
	Hooks() Hooks
	Destroy() error
	CommitAsync()
}

var _ VNodeAccessor = &vnode{}

type vnodeAttrs struct {
	key         any
	typ         Component
	modifiers   modifier.Modifiers
	childVNodes []VNode
}

type vnode struct {
	vnodeAttrs

	parent              VNodeAccessor
	childVNodeAccessors []VNodeAccessor
	node                node.Node
	hooks               Hooks

	eventUpdateConsumers []func(ctx context.Context, vnode VNodeAccessor, old VNodeAccessor)
	eventMountConsumers  []func(ctx context.Context, vnode VNodeAccessor)
}

func (v *vnode) String() string {
	if _, ok := v.typ.(Element); ok {
		if v.node != nil {
			return node.Debug(v.node)
		}
	}
	t := reflect.TypeOf(v.typ)
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return fmt.Sprintf("%s.%p", t.String(), v)
}

func (v *vnode) Hooks() Hooks {
	return v.hooks
}

func (v vnode) Children(children ...VNode) VNode {
	v.childVNodes = children
	return &v
}

func (v *vnode) ChildVNodes() []VNode {
	return v.childVNodes
}

func (v *vnode) ChildVNodeAccessors() []VNodeAccessor {
	return v.childVNodeAccessors
}

func (v *vnode) BindParent(va VNodeAccessor) {
	v.parent = va
}

func (v *vnode) Parent() VNodeAccessor {
	return v.parent
}

func (v *vnode) Modifiers() modifier.Modifiers {
	return v.modifiers
}

func (v *vnode) PutChildVNodes(childVNodes ...VNode) {
	v.childVNodeAccessors = make([]VNodeAccessor, 0, len(childVNodes))

	for i := range childVNodes {
		if va, ok := childVNodes[i].(VNodeAccessor); ok {
			va.BindParent(v)
			v.childVNodeAccessors = append(v.childVNodeAccessors, va)
		}
	}
}

func (v *vnode) ReplaceChildVNodeAccessor(old VNodeAccessor, new VNodeAccessor) {
	i := slices.Index(v.childVNodeAccessors, old)
	if i > -1 {
		new.BindParent(v)
		v.childVNodeAccessors[i] = new
	}
}

func (v *vnode) Type() Component {
	return v.typ
}

func (v *vnode) Key() any {
	return v.key
}

func (v *vnode) OnMount(fn func(ctx context.Context, va VNodeAccessor)) {
	v.eventMountConsumers = append(v.eventMountConsumers, fn)
}

func (v *vnode) BindNode(n node.Node) {
	v.node = n
}

func (v *vnode) Node() node.Node {
	return v.node
}

func (v *vnode) IsRoot() bool {
	if r, ok := v.typ.(interface{ IsRoot() bool }); ok {
		return r.IsRoot()
	}
	return false
}

func (v *vnode) OnUpdate(fn func(ctx context.Context, vn VNodeAccessor, old VNodeAccessor)) {
	v.eventUpdateConsumers = append(v.eventUpdateConsumers, fn)
}

func (v *vnode) Update(ctx context.Context) {
	if eventUpdateConsumers := v.eventUpdateConsumers; eventUpdateConsumers != nil {
		vn := &vnode{vnodeAttrs: v.vnodeAttrs}
		for _, update := range eventUpdateConsumers {
			if update != nil {
				update(ctx, vn, v)
			}
		}
	}
}

func (v *vnode) Use(hook Hook) Hook {
	return v.hooks.use(hook)
}

func (v *vnode) WillRender(oldVNode VNodeAccessor) {
	if oldVNode != nil {
		v.hooks = oldVNode.Hooks()
	}
	v.hooks.Reset()
}

func (v *vnode) DidMount(ctx context.Context) {
	if eventMountConsumers := v.eventMountConsumers; eventMountConsumers != nil {
		for _, mount := range eventMountConsumers {
			if mount != nil {
				mount(ctx, v)
			}
		}
	}
	v.hooks.commit()
}

func (v *vnode) CommitAsync() {
	v.hooks.CommitAsync()
}

func (v *vnode) Destroy() error {
	for _, c := range v.childVNodeAccessors {
		_ = c.Destroy()
	}
	v.childVNodeAccessors = nil
	v.hooks.destroy()
	return nil
}

func (v *vnode) SetKey(key any) {
	v.key = key
}

func SameComponent(type1 Component, type2 Component) bool {
	if w1, ok := type1.(Element); ok {
		if w2, ok := type2.(Element); ok {
			return w1.DisplayName() == w2.DisplayName()
		}
	}

	t1 := reflect.TypeOf(type1)
	for t1.Kind() == reflect.Ptr {
		t1 = t1.Elem()
	}

	t2 := reflect.TypeOf(type2)
	for t2.Kind() == reflect.Ptr {
		t2 = t2.Elem()
	}

	return t1 == t2
}

func UseHook[T Hook](v VNodeAccessor, h T) T {
	return v.Use(h).(T)
}
