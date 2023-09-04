//go:build preview

package m3

import (
	. "github.com/octohelm/gio-compose/pkg/compose"
	"github.com/octohelm/gio-compose/pkg/compose/modifier"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

//go:generate giocompose preview . IconButtonPreview
type IconButtonPreview struct {
}

func (IconButtonPreview) Build(BuildContext) VNode {
	return Row(
		modifier.PaddingAll(40),
		modifier.Gap(20),
	).Children(
		H(FilledIconButton{Icon: Icon(icons.ActionNoteAdd)}),
		H(FilledIconButton{Icon: Icon(icons.ActionNoteAdd), UnSelected: true}),
		H(FilledIconButton{Icon: Icon(icons.ActionNoteAdd), Disabled: true}),

		H(FilledTonalIconButton{Icon: Icon(icons.ActionNoteAdd)}),
		H(FilledTonalIconButton{Icon: Icon(icons.ActionNoteAdd), UnSelected: true}),
		H(FilledTonalIconButton{Icon: Icon(icons.ActionNoteAdd), Disabled: true}),

		H(OutlinedIconButton{Icon: Icon(icons.ActionNoteAdd)}),
		H(OutlinedIconButton{Icon: Icon(icons.ActionNoteAdd), UnSelected: true}),
		H(OutlinedIconButton{Icon: Icon(icons.ActionNoteAdd), Disabled: true}),

		H(StandardIconButton{Icon: Icon(icons.ActionNoteAdd)}),
		H(StandardIconButton{Icon: Icon(icons.ActionNoteAdd), UnSelected: true}),
		H(StandardIconButton{Icon: Icon(icons.ActionNoteAdd), Disabled: true}),
	)
}
