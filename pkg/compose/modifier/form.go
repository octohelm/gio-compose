package modifier

import (
	"github.com/octohelm/gio-compose/pkg/compose"
	"github.com/octohelm/gio-compose/pkg/modifier"
)

func OnSubmit(action func()) modifier.Modifier[any] {
	return &submitEventWatcher{action: action}
}

type submitEventWatcher struct {
	action func()
}

func (s *submitEventWatcher) Modify(target any) {
	if p, ok := target.(compose.SubmitEventProvider); ok {
		p.WatchSubmitEvent(s.action)
	}
}
