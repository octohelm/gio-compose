package node

func IsFragment(v Node) bool {
	if _, ok := v.(*Fragment); ok {
		return ok
	}
	return false
}

type Fragment struct {
	Element
}

func (Fragment) DisplayName() string {
	return "Fragment"
}

func (Fragment) String() string {
	return "Fragment"
}

var _ Node = &Element{}

type Element struct {
	Name                                                    string
	parent, firstChild, lastChild, prevSibling, nextSibling Node
}

func (e *Element) SetDisplayName(name string) {
	e.Name = name
}

func (e *Element) DisplayName() string {
	if e.Name == "" {
		return "Unknown"
	}
	return e.Name
}

func (e *Element) String() string {
	return e.DisplayName()
}

func (e *Element) SetParentNode(n Node) {
	e.parent = n
}

func (e *Element) SetFirstChild(n Node) {
	e.firstChild = n
}

func (e *Element) SetLastChild(n Node) {
	e.lastChild = n
}

func (e *Element) SetPreviousSibling(n Node) {
	e.prevSibling = n
}

func (e *Element) SetNextSibling(n Node) {
	e.nextSibling = n
}

func (e *Element) ParentNode() Node {
	if e == nil {
		return nil
	}
	if e.parent == nil {
		return nil
	}
	return e.parent
}

func (e *Element) FirstChild() Node {
	if e == nil {
		return nil
	}
	if e.firstChild == nil {
		return nil
	}
	return e.firstChild
}

func (e *Element) LastChild() Node {
	if e == nil {
		return nil
	}
	if e.lastChild == nil {
		return nil
	}
	return e.lastChild
}

func (e *Element) PreviousSibling() Node {
	if e == nil {
		return nil
	}
	if e.prevSibling == nil {
		return nil
	}
	return e.prevSibling
}

func (e *Element) NextSibling() Node {
	if e == nil {
		return nil
	}
	if e.nextSibling == nil {
		return nil
	}
	return e.nextSibling
}
