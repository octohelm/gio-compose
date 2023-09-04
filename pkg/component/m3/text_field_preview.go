//go:build preview

package m3

import (
	. "github.com/octohelm/gio-compose/pkg/compose"
	"github.com/octohelm/gio-compose/pkg/compose/modifier"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

//go:generate giocompose preview . TextFieldPreview
type TextFieldPreview struct {
}

func (TextFieldPreview) Build(BuildContext) VNode {
	return Column(
		modifier.PaddingAll(40),
		modifier.Gap(20),
	).Children(
		H(FilledTextField[string]{
			LeadingIcon:  Icon(icons.ActionSearch),
			Label:        Text("Label"),
			TrailingIcon: Icon(icons.NavigationCancel),
			Supporting:   Text("Supporting"),
		}),

		H(FilledTextField[string]{
			LeadingIcon:  Icon(icons.ActionSearch),
			Label:        Text("Label"),
			TrailingIcon: Icon(icons.NavigationCancel),
			Supporting:   Text("Supporting"),
			Value:        "Input",
		}),

		H(FilledTextField[string]{
			LeadingIcon:  Icon(icons.ActionSearch),
			Label:        Text("Label"),
			TrailingIcon: Icon(icons.NavigationCancel),
			Supporting:   Text("Supporting"),
			HasError:     true,
			Value:        "Input Failure",
		}),

		H(FilledTextField[string]{
			LeadingIcon:  Icon(icons.ActionSearch),
			Label:        Text("Label"),
			TrailingIcon: Icon(icons.NavigationCancel),
			Supporting:   Text("Supporting"),
			Disabled:     true,
		}),

		//
		H(OutlinedTextField[string]{
			LeadingIcon:  Icon(icons.ActionSearch),
			Label:        Text("Label"),
			TrailingIcon: Icon(icons.NavigationCancel),
			Supporting:   Text("Supporting"),
		}),

		H(OutlinedTextField[string]{
			LeadingIcon:  Icon(icons.ActionSearch),
			Label:        Text("Label"),
			TrailingIcon: Icon(icons.NavigationCancel),
			Supporting:   Text("Supporting"),
			Value:        "Input",
		}),

		H(OutlinedTextField[string]{
			LeadingIcon:  Icon(icons.ActionSearch),
			Label:        Text("Label"),
			TrailingIcon: Icon(icons.NavigationCancel),
			Supporting:   Text("Supporting"),
			HasError:     true,
			Value:        "Input Failure",
		}),

		H(OutlinedTextField[string]{
			LeadingIcon:  Icon(icons.ActionSearch),
			Label:        Text("Label"),
			TrailingIcon: Icon(icons.NavigationCancel),
			Supporting:   Text("Supporting"),
			Disabled:     true,
		}),
	)
}
