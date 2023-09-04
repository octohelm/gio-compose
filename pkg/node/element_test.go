package node

import (
	"bytes"
	"testing"

	testingx "github.com/octohelm/x/testing"
)

func TestNode(t *testing.T) {
	w := &Fragment{}
	root := &Element{Name: "AppRoot"}

	AppendChild(w, root)

	t.Run("Should append child", func(t *testing.T) {
		c := &Element{Name: "Child"}
		AppendChild(root, c)
		testingx.Expect(t, toXMLString(w), testingx.Be("<AppRoot><Child/></AppRoot>"))

		t.Run("Should replace", func(t *testing.T) {
			InsertBefore(root, &Element{Name: "Child1"}, c)
			RemoveChild(root, c)

			testingx.Expect(t, toXMLString(w), testingx.Be("<AppRoot><Child1/></AppRoot>"))

			t.Run("Should append child again", func(t *testing.T) {
				c2 := &Element{Name: "Child2"}
				AppendChild(root, c2)
				testingx.Expect(t, toXMLString(w), testingx.Be("<AppRoot><Child1/><Child2/></AppRoot>"))

				t.Run("Should insert before", func(t *testing.T) {
					c3 := &Element{Name: "Child3"}
					InsertBefore(root, c3, c2)
					testingx.Expect(t, toXMLString(w), testingx.Be("<AppRoot><Child1/><Child3/><Child2/></AppRoot>"))

					t.Run("Should remove child", func(t *testing.T) {
						RemoveChild(root, c3)
						testingx.Expect(t, toXMLString(w), testingx.Be("<AppRoot><Child1/><Child2/></AppRoot>"))
					})
				})
			})
		})

	})

}

func toXMLString(n Node) string {
	b := bytes.NewBuffer(nil)
	RenderToString(b, n)
	return b.String()
}
