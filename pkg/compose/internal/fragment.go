package internal

var _ Component = Fragment{}

type Fragment struct{}

func (f Fragment) Build(BuildContext) VNode {
	return nil
}
