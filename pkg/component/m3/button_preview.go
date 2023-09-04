//go:build preview

package m3

import (
	. "github.com/octohelm/gio-compose/pkg/compose"
	"github.com/octohelm/gio-compose/pkg/compose/modifier"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

//go:generate giocompose preview . ButtonPreview
type ButtonPreview struct {
}

func (ButtonPreview) Build(BuildContext) VNode {
	return Column(
		modifier.PaddingAll(40),
		modifier.Gap(40),
	).Children(
		Row(modifier.FillMaxWidth(), modifier.Gap(20)).Children(
			H(ElevatedButton{
				Label: Text("Elevated Button"),
			}),
			H(ElevatedButton{
				Label:    Text("Elevated Button"),
				Disabled: true,
			}),
			H(ElevatedButton{
				LeadingIcon: Icon(icons.ActionNoteAdd),
				Label:       Text("提交"),
			}),
			H(ElevatedButton{
				LeadingIcon: Icon(icons.ActionNoteAdd),
				Label:       Text("Elevated Button"),
				Disabled:    true,
			}),
		),
		Row(modifier.FillMaxWidth(), modifier.Gap(20)).Children(
			H(FilledButton{
				Label: Text("Filled Button"),
			}),
			H(FilledButton{
				Label:    Text("Filled Button"),
				Disabled: true,
			}),
		),
		Row(modifier.FillMaxWidth(), modifier.Gap(20)).Children(
			H(FilledTonalButton{
				Label: Text("Filled Tonal Button"),
			}),
			H(FilledTonalButton{
				Label:    Text("Filled Tonal Button"),
				Disabled: true,
			}),
		),
		Row(modifier.FillMaxWidth(), modifier.Gap(20)).Children(
			H(OutlinedButton{
				Label: Text("Outlined Button"),
			}),
			H(OutlinedButton{
				Label:    Text("Outlined Button"),
				Disabled: true,
			}),
		),
		Row(modifier.FillMaxWidth(), modifier.Gap(20)).Children(
			H(TextButton{
				Label: Text("Text Button"),
			}),
			H(TextButton{
				Label:    Text("Text Button"),
				Disabled: true,
			}),
		),
	)
}
