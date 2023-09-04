package contextutil

import (
	"context"

	"github.com/pkg/errors"
)

type Context[T any] interface {
	Inject(ctx context.Context, v T) context.Context
	Extract(ctx context.Context) T
}

type ContextOptFunc[T any] func(opt *ctx[T])

func Defaulter[T any](d func() T) ContextOptFunc[T] {
	return func(opt *ctx[T]) {
		opt.defaulter = d
	}
}

func New[T any](optFuncs ...ContextOptFunc[T]) Context[T] {
	c := &ctx[T]{}

	for _, fn := range optFuncs {
		fn(c)
	}

	return c
}

type ctx[T any] struct {
	defaulter func() T
}

func (c *ctx[T]) Inject(ctx context.Context, v T) context.Context {
	return context.WithValue(ctx, c, v)
}

func (c *ctx[T]) Extract(ctx context.Context) T {
	if v, ok := ctx.Value(c).(T); ok {
		return v
	}

	if c.defaulter != nil {
		return c.defaulter()
	}

	panic(errors.Errorf("need ctx for %T", c))
}
