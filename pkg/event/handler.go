package event

type Handler[T comparable] interface {
	Type() T
	Metadata() []any
	Handle(data any)
}

func NewHandler[T comparable](typ T, handle func(), metadata ...any) Handler[T] {
	if handle == nil {
		return nil
	}

	return &handler[T, any]{
		typ:      typ,
		metadata: metadata,
		handle: func(evt event[T, any]) {
			handle()
		},
	}
}

type HandleFunc[D any] func(data D)

func NewHandlerWithEventData[T comparable, D any](typ T, handle HandleFunc[D], metadata ...any) Handler[T] {
	if handle == nil {
		return nil
	}

	return &handler[T, D]{
		typ:      typ,
		metadata: metadata,
		handle: func(evt event[T, D]) {
			handle(evt.data)
		},
	}
}

type handler[T comparable, D any] struct {
	typ      T
	metadata []any
	handle   func(evt event[T, D])
}

func (h *handler[T, D]) Type() T {
	return h.typ
}

func (h *handler[T, D]) Metadata() []any {
	return h.metadata
}

func (h *handler[T, D]) Handle(data any) {
	if d, ok := data.(D); ok {
		h.handle(event[T, D]{
			typ:  h.typ,
			data: d,
		})
	} else {
		h.handle(event[T, D]{
			typ: h.typ,
		})
	}
}

type event[T comparable, D any] struct {
	typ  T
	data D
}
