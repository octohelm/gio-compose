package main

import "github.com/octohelm/gio-compose/pkg/compose"

type Preview struct {
}

func (Preview) Build(ctx compose.BuildContext) compose.VNode {
	return nil
}
