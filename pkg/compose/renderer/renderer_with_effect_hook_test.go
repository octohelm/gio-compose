package renderer_test

import (
	"context"
	"testing"

	"gioui.org/app"
	testingx "github.com/octohelm/x/testing"

	. "github.com/octohelm/gio-compose/pkg/compose"
	"github.com/octohelm/gio-compose/pkg/compose/renderer"
)

type AppWithEffectHook struct {
	Setup   *func() func()
	Refresh string
}

func (a AppWithEffectHook) Build(ctx BuildContext) VNode {
	UseEffect(ctx, *a.Setup, []any{a.Refresh})
	return nil
}

func TestRenderWithEffectHook(t *testing.T) {
	ctx := context.Background()

	r := renderer.CreateRoot(app.NewWindow())

	counts := struct {
		setup   int
		cleanup int
	}{}

	setup := func() func() {
		counts.setup++
		return func() {
			counts.cleanup++
		}
	}

	r.Render(ctx, Box().Children(H(AppWithEffectHook{
		Setup: &setup,
	})))

	testingx.Expect(t, counts.setup, testingx.Be(1))
	testingx.Expect(t, counts.cleanup, testingx.Be(0))

	t.Run("should not re setup when deps not changed", func(t *testing.T) {
		r.Render(ctx, Box().Children(H(AppWithEffectHook{
			Setup: &setup,
		})))

		testingx.Expect(t, counts.setup, testingx.Be(1))
		testingx.Expect(t, counts.cleanup, testingx.Be(0))

		t.Run("should re setup when deps changed", func(t *testing.T) {
			r.Render(ctx, Box().Children(H(AppWithEffectHook{
				Setup:   &setup,
				Refresh: "1",
			})))

			testingx.Expect(t, counts.setup, testingx.Be(2))
			testingx.Expect(t, counts.cleanup, testingx.Be(1))

			t.Run("should cleanup when destroyed", func(t *testing.T) {
				r.Render(ctx, Box())

				testingx.Expect(t, counts.setup, testingx.Be(2))
				testingx.Expect(t, counts.cleanup, testingx.Be(2))
			})
		})
	})

}
