package compose

import (
	"context"

	"github.com/octohelm/gio-compose/pkg/modifier"

	"github.com/octohelm/gio-compose/pkg/util/contextutil"
)

func FormController(modifiers ...modifier.Modifier[any]) VNode {
	return H(&formController{}, modifiers...)
}

type SubmitEventProvider interface {
	WatchSubmitEvent(fn func())
	TriggerSubmit()
}

type formController struct {
	submitEventConsumers []func()
}

func (f *formController) WatchSubmitEvent(fn func()) {
	f.submitEventConsumers = append(f.submitEventConsumers, fn)
}

func (f *formController) TriggerSubmit() {
	for i := range f.submitEventConsumers {
		if h := f.submitEventConsumers[i]; h != nil {
			h()
		}
	}
}

func (f *formController) Build(buildContext BuildContext) VNode {
	for _, m := range buildContext.Modifiers() {
		if m != nil {
			m.Modify(f)
		}
	}

	return Provider(func(ctx context.Context) context.Context {
		return SubmitEventProviderContext.Inject(ctx, f)
	}).Children(
		buildContext.ChildVNodes()...,
	)
}

var SubmitEventProviderContext = contextutil.New[SubmitEventProvider](
	contextutil.Defaulter(func() SubmitEventProvider {
		return &formController{}
	}),
)
