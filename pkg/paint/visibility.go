package paint

type VisibleSetter interface {
	SetVisible(visible bool)
}

var _ VisibleSetter = &Visibility{}

type Visibility struct {
	Visible bool
}

func (v *Visibility) SetVisible(visible bool) {
	v.Visible = visible
}
