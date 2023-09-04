package renderer

import (
	"context"
	"image"
	"image/color"
	"io"

	"github.com/octohelm/gio-compose/pkg/modifier"

	"gioui.org/app"
	"gioui.org/io/system"
	giolayout "gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	giopaint "gioui.org/op/paint"
	"github.com/octohelm/gio-compose/pkg/compose"
	"github.com/octohelm/gio-compose/pkg/compose/internal"
	"github.com/octohelm/gio-compose/pkg/layout"
	"github.com/octohelm/gio-compose/pkg/node"
	"github.com/octohelm/gio-compose/pkg/paint"
	"github.com/octohelm/gio-compose/pkg/util/debounce"
)

func WithGuideline(enabled bool) RootOptionFunc {
	return func(r *root) {
		r.guidelineEnabled = enabled
	}
}

func WithBoundingRectPrinter(w io.Writer) RootOptionFunc {
	return func(r *root) {
		r.boundingRectPrinter = w
	}
}

type RootOptionFunc func(r *root)

type Root interface {
	Act(fn func())
	Loop() error

	Render(ctx context.Context, vnode compose.VNode)
	RootNode() node.Node
	WindowNode() node.Node
}

func CreateRoot(w *app.Window, optFns ...RootOptionFunc) Root {
	ctx := context.Background()

	windowWidget := CreateRootWidget(ctx, "Window")
	rootWidget := CreateRootWidget(ctx, "Root")

	node.AppendChild(windowWidget, rootWidget)

	r := &root{
		w:       w,
		window:  windowWidget,
		appRoot: compose.Portal(rootWidget).(internal.VNodeAccessor),
		invalid: debounce.Debouncing{BufferTime: bufferTime},
	}

	for i := range optFns {
		optFns[i](r)
	}

	r.commitQueue.Start()

	return r
}

type root struct {
	window  node.Node
	appRoot internal.VNodeAccessor

	w *app.Window

	commitQueue commitQueue
	invalid     debounce.Debouncing

	guidelineEnabled    bool
	boundingRectPrinter io.Writer
}

func (r *root) WindowNode() node.Node {
	return r.window
}

func (r *root) RootNode() node.Node {
	return r.appRoot.Node()
}

func (r *root) Act(fn func()) {
	fn()
	commitAllAsync(r.appRoot)
	r.commitQueue.ForceCommit()
}

func commitAllAsync(v internal.VNodeAccessor) {
	for _, c := range v.ChildVNodeAccessors() {
		commitAllAsync(c)
	}
	v.(internal.HookAsyncCommitter).CommitAsync()
}

func (r *root) Close() error {
	return r.commitQueue.Close()
}

func (r *root) Loop() error {
	ops := &op.Ops{}

	for e := range r.w.Events() {
		switch evt := e.(type) {
		case system.DestroyEvent:
			return evt.Err
		case system.FrameEvent:
			gtx := giolayout.NewContext(ops, evt)

			r.window.(compose.Element).Layout(gtx)

			if r.guidelineEnabled {
				r.drawGuideline(gtx)
			}

			if r.boundingRectPrinter != nil {
				layout.PrintBoundingRectTo(r.boundingRectPrinter, r.window)
			}

			evt.Frame(gtx.Ops)
		}
	}

	return nil
}

func (r *root) drawGuideline(gtx layout.Context) {
	size := gtx.Dp(10)
	rect := paint.Group(gtx.Ops, func() {
		giopaint.FillShape(gtx.Ops, color.NRGBA{A: 0x10}, clip.Rect{Max: image.Pt(size, size)}.Op())
	})

	for x := 0; x <= gtx.Constraints.Max.X/size; x++ {
		for y := 0; y <= gtx.Constraints.Max.Y/size; y++ {
			if (x%2 == 0 && y%2 == 1) || (x%2 == 1 && y%2 == 0) {
				func() {
					defer op.Offset(image.Pt(x*size, y*size)).Push(gtx.Ops).Pop()
					rect.Add(gtx.Ops)
				}()
			}

		}
	}
}

func (r *root) Render(ctx context.Context, vnode compose.VNode) {
	ctx = WindowContext.Inject(ctx, r.window)

	nextRoot := compose.Portal(r.appRoot.Node()).Children(vnode).(internal.VNodeAccessor)
	r.patchVNode(ctx, r.appRoot, nextRoot)
	r.appRoot = nextRoot
	r.commitQueue.ForceCommit()
}

func (r *root) sameVNode(vnode1 internal.VNodeAccessor, vnode2 internal.VNodeAccessor) bool {
	if vnode1 == vnode2 {
		return true
	}
	return internal.SameComponent(vnode1.Type(), vnode2.Type()) && vnode1.Key() == vnode2.Key()
}

func (r *root) patchVNode(ctx context.Context, oldVNode internal.VNodeAccessor, vnode internal.VNodeAccessor) {
	if vnode == oldVNode {
		return
	}

	switch vnode.Type().(type) {
	case internal.Element, internal.Fragment:
		vnode.PutChildVNodes(vnode.ChildVNodes()...)
		r.mount(ctx, oldVNode, vnode)
		r.commitQueue.Commit(func() {
			vnode.DidMount(ctx)
		})
	default:
		// only component need to render
		r.renderVNode(ctx, oldVNode, vnode)
	}
}

func (r *root) renderVNode(ctx context.Context, oldVNode internal.VNodeAccessor, vnode internal.VNodeAccessor) {
	childCtx := ctx

	if cp, ok := vnode.Type().(internal.ContextProvider); ok {
		childCtx = cp.GetChildContext(ctx)
	}

	vnode.WillRender(oldVNode)

	vnode.OnUpdate(func(ctx context.Context, vnode internal.VNodeAccessor, oldVNode internal.VNodeAccessor) {
		r.commitQueue.Commit(func() {
			r.renderVNode(ctx, oldVNode, vnode)
			// after render should replace old vnode of parent
			oldVNode.Parent().ReplaceChildVNodeAccessor(oldVNode, vnode)
		})
	})

	vnode.PutChildVNodes(vnode.Type().Build(&buildContext{
		Context: childCtx,
		vnode:   vnode,
	}))

	r.mount(childCtx, oldVNode, vnode)
	r.commitQueue.Commit(func() {
		vnode.DidMount(ctx)
	})
}

func (r *root) mount(ctx context.Context, oldVNode internal.VNodeAccessor, vnode internal.VNodeAccessor) {
	if oldVNode == nil {
		switch x := vnode.Type().(type) {
		case compose.Element:
			vnode.BindNode(x.New(ctx))
			r.updateVNode(ctx, vnode)
		default:
			if !vnode.IsRoot() {
				vnode.BindNode(&node.Fragment{})
			}
		}

		r.addVNodes(
			ctx,
			vnode.Node(),
			nil,
			vnode.ChildVNodeAccessors(),
			0,
			len(vnode.ChildVNodeAccessors())-1,
		)

		return
	}

	vnode.BindNode(oldVNode.Node())

	r.updateVNode(ctx, vnode)

	mounted := vnode.Node()
	if mounted == nil {
		return
	}

	childVNodes := vnode.ChildVNodeAccessors()
	oldChildVNodes := oldVNode.ChildVNodeAccessors()

	if len(oldChildVNodes) != 0 && len(childVNodes) != 0 {
		r.patchVNodes(ctx, mounted, oldChildVNodes, childVNodes)
	} else if len(childVNodes) != 0 {
		r.addVNodes(ctx, mounted, nil, childVNodes, 0, len(childVNodes)-1)
	} else if len(oldChildVNodes) != 0 {
		r.removeVNodes(ctx, mounted, oldChildVNodes, 0, len(oldChildVNodes)-1)
	}
}

func (r *root) updateVNode(ctx context.Context, vnode internal.VNodeAccessor) {
	if w, ok := vnode.Node().(internal.Element); ok && w != nil {
		if changed := w.Update(ctx, vnode.Modifiers()...); changed {
			r.invalid.Do(func() {
				r.w.Invalidate()
			})
		}
	}
}

func getVNode(l []internal.VNodeAccessor, n, i int) internal.VNodeAccessor {
	if maxIndex := n - 1; maxIndex >= i && i >= 0 {
		return l[i]
	}
	return nil
}

func (r *root) patchVNodes(ctx context.Context, parentNode node.Node, oldChildren []internal.VNodeAccessor, newChildren []internal.VNodeAccessor) {
	// modify from https://github.com/snabbdom/snabbdom/blob/master/src/init.ts
	oldStartIdx := 0
	newStartIdx := 0
	oldN := len(oldChildren)
	oldEndIdx := oldN - 1
	oldStartVNode := getVNode(oldChildren, oldN, 0)
	oldEndVNode := getVNode(oldChildren, oldN, oldEndIdx)
	newN := len(newChildren)
	newEndIdx := newN - 1
	newStartVNode := getVNode(newChildren, newN, 0)
	newEndVNode := getVNode(newChildren, newN, newEndIdx)

	var oldKeyToIdx map[any]int
	var idxInOld int
	var elmToMove internal.VNodeAccessor

	for oldStartIdx <= oldEndIdx && newStartIdx <= newEndIdx {
		if oldStartVNode == nil {
			oldStartIdx++
			oldStartVNode = getVNode(oldChildren, oldN, oldStartIdx) // VNode might have been moved left
		} else if oldEndVNode == nil {
			oldEndIdx--
			oldEndVNode = getVNode(oldChildren, oldN, oldEndIdx)
		} else if newStartVNode == nil {
			newStartIdx++
			newStartVNode = getVNode(newChildren, newN, newStartIdx)
		} else if newEndVNode == nil {
			newEndIdx--
			newEndVNode = getVNode(newChildren, newN, newEndIdx)
		} else if r.sameVNode(oldStartVNode, newStartVNode) {
			r.patchVNode(ctx, oldStartVNode, newStartVNode)
			oldStartIdx++
			oldStartVNode = getVNode(oldChildren, oldN, oldStartIdx)
			newStartIdx++
			newStartVNode = getVNode(newChildren, newN, newStartIdx)
		} else if r.sameVNode(oldEndVNode, newEndVNode) {
			r.patchVNode(ctx, oldEndVNode, newEndVNode)
			oldEndIdx--
			oldEndVNode = getVNode(oldChildren, oldN, oldEndIdx)
			newEndIdx--
			newEndVNode = getVNode(newChildren, newN, newEndIdx)
		} else if r.sameVNode(oldStartVNode, newEndVNode) {
			// VNode moved right
			r.patchVNode(ctx, oldStartVNode, newEndVNode)
			r.commitNode(opMoveRight, parentNode, oldStartVNode.Node(), oldEndVNode.Node().NextSibling())
			oldStartIdx++
			oldStartVNode = getVNode(oldChildren, oldN, oldStartIdx)
			newEndIdx--
			newEndVNode = getVNode(newChildren, newN, newEndIdx)
		} else if r.sameVNode(oldEndVNode, newStartVNode) {
			// VNode moved left
			r.patchVNode(ctx, oldEndVNode, newStartVNode)
			r.commitNode(opMoveLeft, parentNode, oldEndVNode.Node(), oldStartVNode.Node())

			oldEndIdx--
			oldEndVNode = getVNode(oldChildren, oldN, oldEndIdx)
			newStartIdx++
			newStartVNode = getVNode(newChildren, newN, newStartIdx)
		} else {
			if oldKeyToIdx == nil {
				oldKeyToIdx = createKeyToOldIdx(oldChildren, oldStartIdx, oldEndIdx)
			}

			if key := newStartVNode.Key(); key != nil {
				idx, oldKeyIdxExists := oldKeyToIdx[key]
				if oldKeyIdxExists {
					idxInOld = idx
				}
			}

			if idxInOld == 0 {
				r.patchVNode(ctx, nil, newStartVNode)
				r.commitNode(opInsertFirst, parentNode, newStartVNode.Node(), oldStartVNode.Node())
			} else {
				elmToMove = getVNode(oldChildren, oldN, idxInOld)

				if !internal.SameComponent(elmToMove.Type(), newStartVNode.Type()) {
					r.patchVNode(ctx, nil, newStartVNode)
					r.commitNode(opInsert, parentNode, newStartVNode.Node(), oldStartVNode.Node())
				} else {
					r.patchVNode(ctx, elmToMove, newStartVNode)
					oldChildren[idxInOld] = nil
					r.commitNode(opInsert, parentNode, elmToMove.Node(), oldStartVNode.Node())
				}
			}

			newStartIdx++
			newStartVNode = getVNode(newChildren, newN, newStartIdx)
		}
	}

	if newStartIdx <= newEndIdx {
		var beforeNode node.Node
		if n := getVNode(newChildren, newN, newEndIdx+1); n != nil {
			beforeNode = n.Node()
		}
		r.addVNodes(ctx, parentNode, beforeNode, newChildren, newStartIdx, newEndIdx)
	}

	if oldStartIdx <= oldEndIdx {
		r.removeVNodes(ctx, parentNode, oldChildren, oldStartIdx, oldEndIdx)
	}
}

func (r *root) addVNodes(ctx context.Context, parentNode node.Node, beforeNode node.Node, vnodes []internal.VNodeAccessor, startIdx int, endIdx int) {
	for startIdx <= endIdx {
		vn := vnodes[startIdx]
		r.patchVNode(ctx, nil, vn)
		if !vn.IsRoot() {
			r.commitNode(opInsert, parentNode, vn.Node(), beforeNode)
		}
		startIdx++
	}
}

func (r *root) removeVNodes(ctx context.Context, parentNode node.Node, vnodes []internal.VNodeAccessor, startIdx int, endIdx int) {
	for startIdx <= endIdx {
		vn := vnodes[startIdx]
		if !vn.IsRoot() {
			r.commitNode(opRemove, parentNode, vn.Node(), nil)
		}
		r.destroyNode(ctx, vn)
		startIdx++
	}
}

func (r *root) destroyNode(ctx context.Context, vn internal.VNodeAccessor) {
	_ = vn.Destroy()

	// vn destroy should trigger re-render
	r.invalid.Do(func() {
		r.w.Invalidate()
	})
}

func createKeyToOldIdx(childVNodes []internal.VNodeAccessor, beginIdx int, endIdx int) map[any]int {
	keyToIndexMap := map[any]int{}
	for i := beginIdx; i <= endIdx; i++ {
		vn := childVNodes[i]
		key := vn.Key()
		if key != nil {
			keyToIndexMap[key] = i
		}
	}
	return keyToIndexMap
}

func (r *root) commitNode(op operator, parent, child, before node.Node) {
	switch op {
	case opRemove:
		node.RemoveChild(parent, child)
	default:
		if p := child.ParentNode(); p != nil && p == parent {
			node.RemoveChild(parent, child)
		}
		node.InsertBefore(parent, child, before)
	}
}

type buildContext struct {
	context.Context

	vnode internal.VNodeAccessor
}

func (b *buildContext) VNode() internal.VNodeAccessor {
	return b.vnode
}

func (b *buildContext) ChildVNodes() []compose.VNode {
	return b.vnode.ChildVNodes()
}

func (b *buildContext) Modifiers() modifier.Modifiers {
	return b.vnode.Modifiers()
}

func (b *buildContext) RawContext() context.Context {
	return b.Context
}

type operator int

const (
	opRemove operator = iota
	opInsert
	opInsertFirst
	opMoveRight
	opMoveLeft
)

func (o operator) String() string {
	switch o {
	case opInsertFirst:
		return "INSERT_FIRST"
	case opMoveRight:
		return "MOVE_RIGHT"
	case opMoveLeft:
		return "MOVE_LEFT"
	case opInsert:
		return "INSERT"
	}

	return "REMOVE"
}
