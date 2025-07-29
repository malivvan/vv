package goldmark_test

import (
	"testing"

	. "github.com/malivvan/vv/pkg/cui/mdview/goldmark"
	"github.com/malivvan/vv/pkg/cui/mdview/goldmark/parser"
	"github.com/malivvan/vv/pkg/cui/mdview/goldmark/testutil"
)

func TestAttributeAndAutoHeadingID(t *testing.T) {
	markdown := New(
		WithParserOptions(
			parser.WithAttribute(),
			parser.WithAutoHeadingID(),
		),
	)
	testutil.DoTestCaseFile(markdown, "_test/options.txt", t)
}
