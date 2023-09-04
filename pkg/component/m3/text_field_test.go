package m3

import (
	"context"
	"fmt"
	"testing"

	"github.com/octohelm/gio-compose/pkg/component/m3/theming"
	"github.com/octohelm/gio-compose/pkg/compose/testutil"

	"gioui.org/app"

	. "github.com/octohelm/gio-compose/pkg/compose"
	"github.com/octohelm/gio-compose/pkg/compose/renderer"
)

func TestFilledTextField(t *testing.T) {
	ctx := theming.Context.Inject(context.Background(), theming.New())

	r := renderer.CreateRoot(app.NewWindow())

	for i := 0; i < 3; i++ {
		r.Render(ctx, FormController().Children(
			H(FilledTextField[string]{
				Label: Text(fmt.Sprintf("Label %v", i)),
				Value: Echo(func() string {
					if i > 0 {
						return fmt.Sprintf("%v", i)
					}
					return ""
				}),
				OnValueChange: func(v string) {

				},
			}),
		))

		r.Act(func() {

		})

		testutil.ExpectNodeRenderedEqual(t, r.RootNode().FirstChild(), `
<Container>
	<Box>
		<Row>
			<Column>
				<Label>
					<Text/>
				</Label>
				<Input/>
			</Column>
		</Row>
		<StateLayer/>
	</Box>
</Container>
`)
	}

}
