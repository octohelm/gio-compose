package testutil

import (
	"bytes"
	"io"
	"testing"

	"github.com/tdewolff/minify/v2"

	"github.com/octohelm/gio-compose/pkg/node"
	testingx "github.com/octohelm/x/testing"
	"github.com/tdewolff/minify/v2/xml"
)

func ExpectNodeRenderedEqual(t testing.TB, n node.Node, expectString string) {
	t.Helper()

	b := bytes.NewBuffer(nil)
	node.RenderToString(b, n)
	testingx.Expect(t, b.String(), testingx.Equal(minifyXML(bytes.NewBufferString(expectString))))
}

func minifyXML(r io.Reader) string {
	w := bytes.NewBuffer(nil)
	_ = xml.Minify(minify.New(), w, r, nil)
	return w.String()
}
