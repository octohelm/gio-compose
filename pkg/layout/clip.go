package layout

type ClipSetter interface {
	SetClip(enabled bool)
}

var _ ClipSetter = &Clip{}

type Clip struct {
	Enabled bool
}

func (c *Clip) SetClip(enabled bool) {
	c.Enabled = enabled
}
