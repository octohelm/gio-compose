package node

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
)

type DisplayNameSetter interface {
	SetDisplayName(name string)
}

type Node interface {
	DisplayName() string

	ParentNode() Node
	FirstChild() Node
	LastChild() Node

	PreviousSibling() Node
	NextSibling() Node

	SetParentNode(n Node)
	SetFirstChild(n Node)
	SetLastChild(n Node)
	SetPreviousSibling(p Node)
	SetNextSibling(p Node)
}

func InsertBefore(parent, newChildNode, oldChildNode Node) Node {
	if parent == nil {
		panic("insert before nil parent")
	}

	if oldChildNode == nil {
		AppendChild(parent, newChildNode)
		return nil
	}

	if newChildNode.ParentNode() != nil || newChildNode.PreviousSibling() != nil || newChildNode.NextSibling() != nil {
		panic(errors.Errorf("insertBefore called for an attached child Element: %s.%p", newChildNode, newChildNode))
	}

	var prev, next Node

	if oldChildNode != nil {
		prev, next = oldChildNode.PreviousSibling(), oldChildNode
	} else {
		prev = parent.LastChild()
	}
	if prev != nil {
		prev.SetNextSibling(newChildNode)
	} else {
		parent.SetFirstChild(newChildNode)
	}
	if next != nil {
		next.SetPreviousSibling(newChildNode)
	} else {
		parent.SetLastChild(newChildNode)
	}

	newChildNode.SetParentNode(parent)
	newChildNode.SetPreviousSibling(prev)
	newChildNode.SetNextSibling(next)

	return oldChildNode
}

func AppendChild(parent Node, child Node) {
	if parent == nil {
		panic("append to nil parent")
	}

	if child.ParentNode() != nil || child.PreviousSibling() != nil || child.NextSibling() != nil {
		panic(errors.Errorf("appendChild called for an attached child Element: %s", Debug(child)))
	}

	lastChild := parent.LastChild()
	if lastChild != nil {
		lastChild.SetNextSibling(child)
	} else {
		parent.SetFirstChild(child)
	}

	parent.SetLastChild(child)
	child.SetParentNode(parent)
	child.SetPreviousSibling(lastChild)
}

func RemoveChild(parent Node, child Node) Node {
	if child == nil {
		return nil
	}

	if cp := child.ParentNode(); cp != parent {
		panic(errors.Errorf("remove a non-child node %s@%s of %s", Debug(child), Debug(cp), Debug(parent)))
	}

	if parent.FirstChild() == child {
		parent.SetFirstChild(child.NextSibling())
	}

	if child.NextSibling() != nil {
		child.NextSibling().SetPreviousSibling(child.PreviousSibling())
	}

	if parent.LastChild() == child {
		parent.SetLastChild(child.PreviousSibling())
	}

	if child.PreviousSibling() != nil {
		child.PreviousSibling().SetNextSibling(child.NextSibling())
	}

	// cleanup
	child.SetParentNode(nil)
	child.SetPreviousSibling(nil)
	child.SetNextSibling(nil)

	return child
}

func IterChildElement(ctx context.Context, n Node) <-chan Node {
	ch := make(chan Node)

	go func() {
		defer close(ch)

		for c := n.FirstChild(); c != nil; c = c.NextSibling() {
			if IsFragment(c) {
				for fc := range IterChildElement(ctx, c) {
					select {
					case <-ctx.Done():
						return
					case ch <- fc:
					}
				}
				continue
			}

			select {
			case <-ctx.Done():
				return
			case ch <- c:
			}
		}
	}()

	return ch
}

func Debug(n Node) string {
	if n == nil {
		return "nil"
	}
	return fmt.Sprintf("%s.%p", n, n)
}
