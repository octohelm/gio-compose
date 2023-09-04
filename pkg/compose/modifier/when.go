package modifier

import (
	"github.com/octohelm/gio-compose/pkg/modifier"
)

func When[T any](when bool, modifiers ...modifier.Modifier[T]) modifier.Modifier[T] {
	return modifier.When(when, modifiers...)
}
