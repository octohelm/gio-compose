package modifier

import (
	"github.com/octohelm/gio-compose/pkg/event"
	"github.com/octohelm/gio-compose/pkg/event/textinput"
	"github.com/octohelm/gio-compose/pkg/modifier"
)

func OnValueChange(action func(v string)) modifier.Modifier[any] {
	return &valueChangeWatcher{
		action: textinput.OnChange(func(data *textinput.ChangeData) {
			action(data.Value)
		}),
	}
}

type valueChangeWatcher struct {
	action event.Handler[textinput.Input]
}

func (w *valueChangeWatcher) Modify(target any) {
	if p, ok := target.(textinput.EventsAccessor); ok {
		p.InputEvents().Add(w.action)
	}
}
