package extension

import (
	"testing"

	"github.com/malivvan/vv/pkg/cui/mdview/goldmark"
	"github.com/malivvan/vv/pkg/cui/mdview/goldmark/renderer/html"
	"github.com/malivvan/vv/pkg/cui/mdview/goldmark/testutil"
)

func TestTypographer(t *testing.T) {
	markdown := goldmark.New(
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
		),
		goldmark.WithExtensions(
			Typographer,
		),
	)
	testutil.DoTestCaseFile(markdown, "_test/typographer.txt", t)
}
