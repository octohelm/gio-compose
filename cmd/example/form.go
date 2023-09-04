package main

import (
	"fmt"

	"github.com/octohelm/gio-compose/pkg/component/m3"
	. "github.com/octohelm/gio-compose/pkg/compose"
	"github.com/octohelm/gio-compose/pkg/compose/modifier"
)

func init() {
	AddShowCase("ux/form", Form{})
}

type Form struct {
}

func (Form) Build(b BuildContext) VNode {
	formData := UseState(b, UserInfo{})

	return Column(
		modifier.FillMaxSize(),
		modifier.PaddingAll(20),
		modifier.Gap(20),
	).Children(
		FormController(modifier.OnSubmit(func() {
			fmt.Println("Submitted", formData.Value())
		})).Children(
			H(m3.FilledTextField[string]{
				Label: Text("First name"),
				Value: formData.Value().FirstName,
				OnValueChange: func(v string) {
					formData.UpdateFunc(func(info UserInfo) UserInfo {
						info.FirstName = v
						return info
					})
				},
			}),
			H(m3.FilledTextField[string]{
				Label: Text("Last name"),
				Value: formData.Value().LastName,
				OnValueChange: func(v string) {
					formData.UpdateFunc(func(info UserInfo) UserInfo {
						info.LastName = v
						return info
					})
				},
			}),
		),
	)
}

type UserInfo struct {
	FirstName string
	LastName  string
}
