package modifier

import (
	"github.com/octohelm/gio-compose/pkg/layout"
	"github.com/octohelm/gio-compose/pkg/modifier"
)

func DetectLayout(handlers ...layout.PhaseHandler) modifier.Modifier[any] {
	return &layoutEventWatcher{
		handlers: handlers,
	}
}

type layoutEventWatcher struct {
	handlers []layout.PhaseHandler
}

func (g *layoutEventWatcher) Modify(target any) {
	if setter, ok := target.(layout.PhaseHandlersSetter); ok {
		setter.SetPhaseHandlers(g.handlers...)
	}
}
