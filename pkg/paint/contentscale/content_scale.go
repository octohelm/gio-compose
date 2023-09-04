package contentscale

type ContentScale int

const (
	None ContentScale = iota
	Inside
	FillBounds
	FillWidth
	FillHeight
	Crop
	Fit
)
