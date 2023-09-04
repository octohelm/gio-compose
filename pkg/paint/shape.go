package paint

import (
	"image"

	"github.com/octohelm/gio-compose/pkg/layout/position"

	"gioui.org/layout"

	"github.com/octohelm/gio-compose/pkg/cmp"
	"github.com/octohelm/gio-compose/pkg/paint/canvas"
	"github.com/octohelm/gio-compose/pkg/paint/size"
	"github.com/octohelm/gio-compose/pkg/unit"
)

type Shape interface {
	Equals(v any) cmp.Result
	Rectangle(gtx layout.Context) image.Rectangle
	Path(gtx layout.Context, positions ...position.Position) *canvas.Path
}

type SizedShape interface {
	Shape

	SizeSetter
	SizedChecker
}

type SizeSetter interface {
	SetSize(v unit.Dp, typs ...size.Type)
}

type SizedChecker interface {
	Sized(typ size.Type) size.SizingType
}
