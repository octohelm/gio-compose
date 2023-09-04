package modifier

import (
	"context"
	"testing"

	testingx "github.com/octohelm/x/testing"
)

func TestIterModifier(t *testing.T) {
	t.Run("Iter", func(t *testing.T) {
		modifiers := make([]Modifier[any], 0)

		for m := range Iter[any](
			context.Background(),
			Discord{},
			Modifiers{
				Discord{},
				Modifiers{Discord{}},
			},
			When[any](true, Discord{}),
		) {
			modifiers = append(modifiers, m)
		}

		testingx.Expect(t, modifiers, testingx.HaveLen[[]Modifier[any]](4))
	})

}
