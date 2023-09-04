package internal

var _ Component = Portal{}

type Portal struct{}

func (Portal) IsRoot() bool {
	return true
}

func (Portal) Build(c BuildContext) VNode {
	return H(Fragment{}).Children(c.ChildVNodes()...)
}
