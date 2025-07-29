package extension

import (
	"testing"

	"github.com/malivvan/vv/pkg/cui/mdview/goldmark"
	"github.com/malivvan/vv/pkg/cui/mdview/goldmark/renderer/html"
	"github.com/malivvan/vv/pkg/cui/mdview/goldmark/testutil"
)

func TestFootnote(t *testing.T) {
	markdown := goldmark.New(
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
		),
		goldmark.WithExtensions(
			Footnote,
		),
	)
	testutil.DoTestCaseFile(markdown, "_test/footnote.txt", t)
}
