package fuzz

import (
	"bytes"

	"github.com/malivvan/vv/pkg/cui/mdview/goldmark"
	"github.com/malivvan/vv/pkg/cui/mdview/goldmark/extension"
	"github.com/malivvan/vv/pkg/cui/mdview/goldmark/parser"
	"github.com/malivvan/vv/pkg/cui/mdview/goldmark/renderer/html"
)

// Fuzz runs automated fuzzing against goldmark.
func Fuzz(data []byte) int {
	markdown := goldmark.New(
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
			parser.WithAttribute(),
		),
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
			html.WithXHTML(),
		),
		goldmark.WithExtensions(
			extension.DefinitionList,
			extension.Footnote,
			extension.GFM,
			extension.Linkify,
			extension.Table,
			extension.TaskList,
			extension.Typographer,
		),
	)
	var b bytes.Buffer
	if err := markdown.Convert(data, &b); err != nil {
		return 0
	}

	return 1
}
