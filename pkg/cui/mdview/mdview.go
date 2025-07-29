package mdview

import (
	"bytes"
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/malivvan/vv/pkg/cui"
	"regexp"
	"sort"
	"strings"
	"sync"
	"unicode"
	"unicode/utf8"

	"github.com/alecthomas/chroma"
	"github.com/gdamore/tcell/v2"
	"github.com/malivvan/vv/pkg/cui/mdview/ansicsi"
	"github.com/malivvan/vv/pkg/cui/mdview/goldmark"
	"github.com/malivvan/vv/pkg/cui/mdview/goldmark/ast"
	"github.com/malivvan/vv/pkg/cui/mdview/goldmark/extension"
	goldmark_parser "github.com/malivvan/vv/pkg/cui/mdview/goldmark/parser"
	goldmark_renderer "github.com/malivvan/vv/pkg/cui/mdview/goldmark/renderer"
	"github.com/malivvan/vv/pkg/cui/mdview/goldmark/text"
	"github.com/malivvan/vv/pkg/cui/mdview/goldmark/util"
	"github.com/malivvan/vv/pkg/cui/mdview/indexer"
	"github.com/malivvan/vv/pkg/cui/mdview/renderer"
	"github.com/rivo/uniseg"
)

func cellStyle(default_ tcell.Style, styles ...chroma.StyleEntry) tcell.Style {
	if len(styles) == 0 {
		return default_
	}

	style := default_
	for _, s := range styles {
		style = style.Foreground(tcell.NewRGBColor(int32(s.Colour.Red()), int32(s.Colour.Blue()), int32(s.Colour.Green())))
		style = style.Background(tcell.NewRGBColor(int32(s.Background.Red()), int32(s.Background.Blue()), int32(s.Background.Green())))

		if s.Bold != chroma.Pass {
			style = style.Bold(s.Bold == chroma.Yes)
		}
		if s.Italic != chroma.Pass {
			style = style.Italic(s.Italic == chroma.Yes)
		}
		if s.Underline != chroma.Pass {
			style = style.Underline(s.Underline == chroma.Yes)
		}
	}
	return style
}

type grapheme struct {
	start int
	end   int
	runes []rune
	style tcell.Style
}

func (g *grapheme) len() int {
	return g.end - g.start
}

func (g *grapheme) isSpace() bool {
	// TODO: text segmentation word boundaries https://unicode.org/reports/tr29/#Word_Boundaries
	for _, r := range g.runes {
		if !unicode.IsSpace(r) {
			return false
		}
	}
	return true
}

func stringGraphemes(s string) []grapheme {
	var graphemes []grapheme

	it := uniseg.NewGraphemes(string(s))
	for it.Next() {
		graphemes = append(graphemes, grapheme{
			runes: it.Runes(),
		})
	}
	return graphemes
}

type line struct {
	start     int
	end       int
	graphemes []grapheme
}

type lineWriter struct {
	byteOffset   int
	buf          bytes.Buffer
	lines        []line
	longestLine  int
	defaultStyle tcell.Style
	style        tcell.Style
}

func (w *lineWriter) updateStyle(sgr *ansicsi.SetGraphicsRendition) {
	switch sgr.Command {
	case ansicsi.SGRReset:
		w.style = w.defaultStyle
	case ansicsi.SGRBold:
		w.style = w.style.Bold(true)
	case ansicsi.SGRFaint:
		w.style = w.style.Dim(true)
	case ansicsi.SGRItalic:
		w.style = w.style.Italic(true)
	case ansicsi.SGRUnderline:
		w.style = w.style.Underline(true)
	case ansicsi.SGRNormalWeight:
		w.style = w.style.Bold(false).Dim(false)
	case ansicsi.SGRNoItalicOrFraktur:
		w.style = w.style.Italic(false)
	case ansicsi.SGRNoUnderline:
		w.style = w.style.Underline(false)
	case ansicsi.SGRForegroundBlack, ansicsi.SGRForegroundRed, ansicsi.SGRForegroundGreen, ansicsi.SGRForegroundYellow, ansicsi.SGRForegroundBlue, ansicsi.SGRForegroundMagenta, ansicsi.SGRForegroundCyan, ansicsi.SGRForegroundWhite:
		w.style = w.style.Foreground(tcell.Color(int32(sgr.Command & 0x7)))
	case ansicsi.SGRForegroundDefault:
		fg, _, _ := w.defaultStyle.Decompose()
		w.style = w.style.Foreground(fg)
	case ansicsi.SGRBackgroundBlack, ansicsi.SGRBackgroundRed, ansicsi.SGRBackgroundGreen, ansicsi.SGRBackgroundYellow, ansicsi.SGRBackgroundBlue, ansicsi.SGRBackgroundMagenta, ansicsi.SGRBackgroundCyan, ansicsi.SGRBackgroundWhite:
		w.style = w.style.Background(tcell.Color(int32(sgr.Command & 0x7)))
	case ansicsi.SGRBackgroundDefault:
		_, bg, _ := w.defaultStyle.Decompose()
		w.style = w.style.Background(bg)
	case ansicsi.SGRForegroundColor:
		switch sgr.Parameters[0] {
		case 2:
			w.style = w.style.Foreground(tcell.NewRGBColor(int32(sgr.Parameters[1]), int32(sgr.Parameters[2]), int32(sgr.Parameters[3])))
		case 5:
			w.style = w.style.Foreground(tcell.Color(int32(sgr.Parameters[1] & 0xff)))
		}
	case ansicsi.SGRBackgroundColor:
		switch sgr.Parameters[0] {
		case 2:
			w.style = w.style.Background(tcell.NewRGBColor(int32(sgr.Parameters[1]), int32(sgr.Parameters[2]), int32(sgr.Parameters[3])))
		case 5:
			w.style = w.style.Background(tcell.Color(int32(sgr.Parameters[1] & 0xff)))
		}
	}
}

func (w *lineWriter) flushLine() {
	w.lines = append(w.lines, line{start: w.byteOffset})
	l := &w.lines[len(w.lines)-1]

	appendGraphemes := func(b []byte) {
		graphemes := uniseg.NewGraphemes(string(b))
		for graphemes.Next() {
			start, end := graphemes.Positions()
			sz := end - start
			l.graphemes = append(l.graphemes, grapheme{
				start: w.byteOffset,
				end:   w.byteOffset + sz,
				runes: graphemes.Runes(),
				style: w.style,
			})
			w.byteOffset += sz
		}
	}

	buf := w.buf.Bytes()
	for start, end := 0, 0; ; {
		if cmd, sz := ansicsi.Decode(buf[end:]); sz != 0 {
			appendGraphemes(buf[start:end])

			if sgr, ok := cmd.(*ansicsi.SetGraphicsRendition); ok {
				w.updateStyle(sgr)
			}

			start = end + sz
			end = start
			w.byteOffset += sz
			continue
		}

		if end >= len(buf) {
			appendGraphemes(buf[start:])
			break
		}

		end++
	}
	if len(l.graphemes)-1 > w.longestLine {
		w.longestLine = len(l.graphemes) - 1
	}

	l.end = w.byteOffset
	w.buf.Reset()
}

func (w *lineWriter) Write(b []byte) (int, error) {
	n := len(b)
	for {
		newline := bytes.IndexByte(b, '\n')
		if newline == -1 {
			w.buf.Write(b)
			return n, nil
		}

		w.buf.Write(b[:newline])
		w.buf.WriteByte(' ')
		w.flushLine()
		b = b[newline+1:]
	}
}

func isLink(n ast.Node) (bool, bool) {
	switch n.Kind() {
	case ast.KindAutoLink, ast.KindImage, ast.KindLink:
		return true, true
	default:
		return false, false
	}
}

func isHeading(n ast.Node) (bool, bool) {
	return false, n.Kind() == ast.KindHeading
}

type MdView struct {
	sync.Mutex
	*cui.Box

	// The colorscheme to use, if any.
	theme *chroma.Style

	// The name of the document.
	name string

	// The raw Markdown.
	markdown []byte

	// The parsed Markdown.
	document ast.Node

	// Node span tree.
	spanTree *renderer.NodeSpan

	// The document index.
	index *indexer.DocumentIndex

	// The selection, if any.
	selection *renderer.NodeSpan

	// The selected span, if any.
	selectionStart, selectionEnd int

	// True if the selected span should be highlighted.
	highlightSelection bool

	// The processed line index. This is nil if the buffer has changed and needs
	// to be re-indexed.
	lines []line

	// The last width for which the current table is drawn.
	lastWidth int

	// The screen width of the longest line in the index (not the buffer).
	longestLine int

	// The index of the first line shown in the text view.
	lineOffset int

	// The number of characters to be skipped on each line (not in wrap mode).
	columnOffset int

	// The height of the content the last time the text view was drawn.
	pageSize int

	// If set to true, lines that are longer than the available width are wrapped
	// onto the next line. If set to false, any characters beyond the available
	// width are discarded.
	wrap bool

	// If set to true, render a gutter with the document name and view position.
	showGutter bool
}

func NewMdView(theme *chroma.Style) *MdView {
	return &MdView{
		Box:   cui.NewBox(),
		theme: theme,
		wrap:  true,
	}
}

// Clear removes all text from the buffer.
func (mv *MdView) Clear() *MdView {
	mv.lines = nil
	mv.markdown = nil
	mv.document = nil
	return mv
}

func (mv *MdView) GetMarkdown() []byte {
	return mv.markdown
}

// SetText sets the text of this text view to the provided string. Previously
// contained text will be removed.
func (mv *MdView) SetText(name, markdown string) *MdView {
	mv.Clear()
	mv.name = name
	mv.markdown = []byte(markdown)
	parser := goldmark.DefaultParser()
	parser.AddOptions(goldmark_parser.WithParagraphTransformers(
		util.Prioritized(extension.NewTableParagraphTransformer(), 200),
	))
	mv.document = parser.Parse(text.NewReader(mv.markdown))
	if doc, ok := mv.document.(*ast.Document); ok {
		mv.index = indexer.Index(doc, mv.markdown)
	}
	return mv
}

// SetWrap sets the flag that, if true, leads to lines that are longer than the
// available width being wrapped onto the next line. If false, any characters
// beyond the available width are not displayed.
func (mv *MdView) SetWrap(wrap bool) *MdView {
	if mv.wrap != wrap {
		mv.lines = nil
	}
	mv.wrap = wrap
	return mv
}

// SetGutter sets the gutter flag, that, if true, instructs the view to render a
// gutter in its bottommost line with the document name and view position.
func (mv *MdView) SetGutter(showGutter bool) *MdView {
	mv.showGutter = showGutter
	return mv
}

// reindexBuffer re-indexes the buffer such that we can use it to easily draw
// the buffer onto the screen. Each line in the index will contain a pointer
// into the buffer from which on we will print text. It will also contain the
// color with which the line starts.
func (mv *MdView) render(width int) {
	if mv.lines != nil {
		return // Nothing has changed. We can still use the set of lines.
	}
	mv.lines = nil

	if mv.document == nil {
		return // No content.
	}

	// Re-render the Markdown into lines.
	wrap := 0
	if mv.wrap {
		wrap = width
	}

	r := renderer.New(
		renderer.WithTheme(mv.theme),
		renderer.WithHyperlinks(true),
		renderer.WithWordWrap(wrap),
		renderer.WithSoftBreak(wrap != 0))

	w := lineWriter{
		style:        tcell.StyleDefault.Foreground(cui.Styles.PrimaryTextColor),
		defaultStyle: tcell.StyleDefault.Foreground(cui.Styles.PrimaryTextColor),
	}
	renderer := goldmark_renderer.NewRenderer(goldmark_renderer.WithNodeRenderers(util.Prioritized(r, 100)))
	if err := renderer.Render(&w, mv.markdown, mv.document); err != nil {
		msg := []rune(fmt.Sprintf("error rendering Markdown: %v", err))
		graphemes := make([]grapheme, len(msg))
		for i, r := range msg {
			graphemes[i] = grapheme{runes: []rune{r}, style: w.defaultStyle}
		}

		mv.lines = []line{{graphemes: graphemes}}
		return
	}
	if w.buf.Len() > 0 {
		w.flushLine()
	}

	mv.spanTree, mv.lines, mv.longestLine = r.SpanTree(), w.lines, w.longestLine
}

// Draw draws this primitive onto the screen.
func (mv *MdView) Draw(screen tcell.Screen) {
	mv.Lock()
	defer mv.Unlock()
	mv.Box.Draw(screen)

	// Get the available size.
	x, y, width, height := mv.GetInnerRect()

	textHeight := height
	if mv.showGutter {
		textHeight = height - 1
	}

	mv.pageSize = textHeight

	// If the width has changed and we're word-wrapping, we need to re-render.
	if width != mv.lastWidth && mv.wrap {
		mv.lines = nil
	}
	mv.lastWidth = width

	// Re-render.
	mv.render(width)

	// If we don't have any lines, there's nothing to draw.
	if mv.lines == nil {
		return
	}

	// Adjust line offset.
	if mv.lineOffset+textHeight > len(mv.lines) {
		mv.lineOffset = len(mv.lines) - textHeight
	}
	if mv.lineOffset < 0 {
		mv.lineOffset = 0
	}

	// Adjust column offset.
	if mv.columnOffset+width > mv.longestLine {
		mv.columnOffset = mv.longestLine - width
	}
	if mv.columnOffset < 0 {
		mv.columnOffset = 0
	}

	// Draw the buffer.
	lastLine := mv.lineOffset + textHeight
	if lastLine > len(mv.lines) {
		lastLine = len(mv.lines)
	}

	// Pull the current style from the end of the preceding line, if any.
	var style tcell.Style
	if mv.lineOffset < len(mv.lines) {
		for i := mv.lineOffset - 1; i >= 0; i-- {
			l := mv.lines[i]
			if len(l.graphemes) > 0 {
				style = l.graphemes[len(l.graphemes)-1].style.Underline(false)
				break
			}
		}
	}

	for i, line := range mv.lines[mv.lineOffset:lastLine] {
		cy := y + i

		if mv.columnOffset > len(line.graphemes) {
			for j := 0; j < width; j++ {
				screen.SetContent(x+j, cy, ' ', nil, style)
			}
			continue
		}

		lastColumn := mv.columnOffset + width
		if lastColumn > len(line.graphemes) {
			lastColumn = len(line.graphemes)
		}
		for j, r := range line.graphemes[mv.columnOffset:lastColumn] {
			cellStyle := r.style
			if mv.selected(r.start) {
				cellStyle = cellStyle.Reverse(true)
			}

			screen.SetContent(x+j, cy, r.runes[0], r.runes[1:], cellStyle)
			style = r.style.Underline(false)
		}
		for j := lastColumn - mv.columnOffset; j < width; j++ {
			screen.SetContent(x+j, cy, ' ', nil, style)
		}
	}

	// Draw the gutter if necessary.
	if mv.showGutter && width >= len("100% ") {
		// Layout: "name {pad} pct% "
		//
		// The document position must be shown. The name will be truncated if necessary.

		nameGraphemes := stringGraphemes(mv.name)
		switch {
		case width-len(" 100% ") > len(nameGraphemes):
			// OK
		case width-len("... 100% ") > 0:
			nameGraphemes = nameGraphemes[:width-len("... 100% ")]
			for _, c := range "..." {
				nameGraphemes = append(nameGraphemes, grapheme{
					runes: []rune{c},
				})
			}
		default:
			nameGraphemes = nil
		}

		defaultStyle := tcell.StyleDefault.Foreground(cui.Styles.PrimaryTextColor)

		style := cellStyle(defaultStyle, mv.theme.Get(chroma.Generic), mv.theme.Get(chroma.Comment))
		col := 0
		for _, r := range nameGraphemes {
			screen.SetContent(x+col, y+height-1, r.runes[0], r.runes[1:], style)
			col++
		}
		for ; col < width-len(" 100% "); col++ {
			screen.SetContent(x+col, y+height-1, ' ', nil, style)
		}

		style = cellStyle(defaultStyle, mv.theme.Get(chroma.Generic), mv.theme.Get(chroma.Text))
		pct := fmt.Sprintf(" % 3d%% ", lastLine*100/len(mv.lines))
		for _, c := range pct {
			screen.SetContent(x+col, y+height-1, c, nil, style)
			col++
		}
	}
}

// InputHandler returns the handler for this primitive.
func (mv *MdView) InputHandler() func(event *tcell.EventKey, setFocus func(p cui.Primitive)) {
	return mv.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p cui.Primitive)) {
		key := event.Key()

		if key == tcell.KeyEscape || key == tcell.KeyEnter || key == tcell.KeyTab || key == tcell.KeyBacktab {
			return
		}

		switch key {
		case tcell.KeyRune:
			switch event.Rune() {
			case 'g': // Home.
				mv.lineOffset = 0
				mv.columnOffset = 0
			case 'G': // End.
				mv.columnOffset = 0
			case 'j': // Down.
				mv.lineOffset++
			case 'k': // Up.
				mv.lineOffset--
			case 'h': // Lefmv.
				mv.columnOffset--
			case 'l': // Righmv.
				mv.columnOffset++
			case '[': // Previous link.
				mv.SelectPrevious(isLink)
			case ']': // Next link.
				mv.SelectNext(isLink)
			case '{': // Previous heading.
				mv.SelectPrevious(isHeading)
			case '}': // Next heading.
				mv.SelectNext(isHeading)
			}
		case tcell.KeyCtrlLeftSq:
			mv.SelectPrevious(func(_ ast.Node) (bool, bool) { return true, true })
		case tcell.KeyCtrlRightSq:
			mv.SelectNext(func(_ ast.Node) (bool, bool) { return true, true })
		case tcell.KeyHome:
			mv.lineOffset = 0
			mv.columnOffset = 0
		case tcell.KeyEnd:
			mv.columnOffset = 0
		case tcell.KeyUp:
			mv.lineOffset--
		case tcell.KeyDown:
			mv.lineOffset++
		case tcell.KeyLeft:
			mv.columnOffset--
		case tcell.KeyRight:
			mv.columnOffset++
		case tcell.KeyPgDn, tcell.KeyCtrlF:
			mv.lineOffset += mv.pageSize
		case tcell.KeyPgUp, tcell.KeyCtrlB:
			mv.lineOffset -= mv.pageSize
		}
	})
}

// Focus is called when this primitive receives focus.
func (mv *MdView) Focus(delegate func(p cui.Primitive)) {
	// Implemented here with locking because this is used by layout primitives.
	mv.Lock()
	defer mv.Unlock()

	mv.Box.Focus(delegate)
}

// HasFocus returns whether or not this primitive has focus.
func (mv *MdView) HasFocus() bool {
	// Implemented here with locking because this may be used in the "changed"
	// callback.
	mv.Lock()
	defer mv.Unlock()

	return mv.Box.HasFocus()
}

func (mv *MdView) scrollToOffset(offset int) {
	// TODO(pdg): only scroll if the offset is not already in view? Maybe vertically
	//            center the line containing the target offset? Whatever's easiest on the
	//            reading experienace.

	start, end := 0, len(mv.lines)
	for start != end {
		i := start + (end-start)/2
		l := mv.lines[i]
		switch {
		case offset < l.start:
			end = i
		case offset >= l.end:
			start = i
		default:
			mv.lineOffset = i
			return
		}
	}
}

type Selector func(n ast.Node) (highlight, ok bool)

func decodeLastValidRune(b []byte) (rune, int) {
	runeStart := len(b) - 1
	for runeStart >= 0 {
		if utf8.RuneStart(b[runeStart]) {
			break
		}
	}

	b = b[runeStart:]
	r, sz := utf8.DecodeRune(b)
	if sz == 0 {
		return r, 0
	}
	return r, len(b)
}

func (mv *MdView) grapheme(at int) *grapheme {
	li := sort.Search(len(mv.lines), func(i int) bool {
		l := &mv.lines[i]
		return l.end >= at
	})
	if li >= len(mv.lines) || at < mv.lines[li].start {
		return nil
	}

	l := &mv.lines[li]
	gi := sort.Search(len(l.graphemes), func(i int) bool {
		g := &l.graphemes[i]
		return g.end >= at
	})
	if gi >= len(l.graphemes) || at < l.graphemes[gi].start {
		return nil
	}

	return &l.graphemes[gi]
}

func (mv *MdView) calculateSelectionSpan(selection *renderer.NodeSpan) {
	// Trim leading and trailing whitespace.
	start, end := selection.Start, selection.End

	for start < end {
		g := mv.grapheme(start)
		if g == nil || !g.isSpace() {
			break
		}
		start = g.end + 1
	}

	for end > start {
		g := mv.grapheme(end - 1)
		if g == nil || !g.isSpace() {
			break
		}
		end = g.start
	}

	mv.selectionStart, mv.selectionEnd = start, end
}

func (mv *MdView) selected(offset int) bool {
	return mv.selection != nil && mv.selectionStart <= offset && offset < mv.selectionEnd
}

func (mv *MdView) Selection() *renderer.NodeSpan {
	return mv.selection
}

// SelectPrevious selects the first node before the current selection that matches the given selector.
func (mv *MdView) SelectPrevious(selector Selector) bool {
	cursor := mv.selection
	if cursor == nil {
		cursor = mv.spanTree
	}
	cursor = cursor.Prev
	if cursor == nil {
		return false
	}

	for cursor = cursor.Prev; cursor != nil; cursor = cursor.Prev {
		if highlight, ok := selector(cursor.Node); ok {
			mv.SelectSpan(cursor, highlight)
			return true
		}
	}

	return false
}

// SelectNext selects the first node after the current selection that matches the given selector.
func (mv *MdView) SelectNext(selector Selector) bool {
	cursor := mv.selection
	if cursor == nil {
		cursor = mv.spanTree
	}
	cursor = cursor.Next
	if cursor == nil {
		return false
	}

	for cursor = cursor.Next; cursor != nil; cursor = cursor.Next {
		if highlight, ok := selector(cursor.Node); ok {
			mv.SelectSpan(cursor, highlight)
			return true
		}
	}

	return false
}

// SelectAnchor selects the next heading with the given anchor.
func (mv *MdView) SelectAnchor(anchor string) bool {
	if mv.index == nil {
		return false
	}

	sections, ok := mv.index.Lookup(anchor)
	if !ok {
		return false
	}
	selector := func(node ast.Node) (bool, bool) {
		for _, s := range sections {
			if s.Start == node {
				return true, true
			}
		}
		return false, false
	}

	if !mv.SelectNext(selector) {
		mv.selection = nil
		mv.SelectNext(selector)
	}
	return true
}

// SelectSpan selects the given node span.
func (mv *MdView) SelectSpan(span *renderer.NodeSpan, highlight bool) {
	mv.highlightSelection = highlight
	mv.selection = span
	mv.calculateSelectionSpan(span)
	mv.scrollToOffset(span.Start)
}

const helpText = `Ctrl+o: open the selected URL in the default browser

]: select the next URL

[: select the previous URL

}: select the next heading

{: select the previous heading

<: go back to the previous selection
`

func textDimensions(text string) (int, int) {
	s, w, h := "", 0, 0
	for len(text) != 0 {
		i := strings.Index(text, "\n")
		if i == -1 {
			s, text = text, ""
		} else {
			s, text = text[:i], text[i+1:]
		}
		if sw := cui.TaggedStringWidth(s); sw > w {
			w = sw
		}
		h++
	}
	return w, h
}

type textDialog struct {
	x, y, w, h int
	visible    bool
	text       string
	textWidth  int
	textHeight int
	textView   *cui.TextView
	persistent bool
}

func (td *textDialog) GetVisible() bool {
	return td.visible
}

func (td *textDialog) SetVisible(v bool) {
	td.visible = v
}

func (td *textDialog) GetFocusable() cui.Focusable {
	return td.textView
}

func newTextDialog(text, title string) *textDialog {
	textView := cui.NewTextView()
	textView.SetBorder(true)
	textView.SetTitle(title)
	textView.SetWrap(true)
	textView.SetWordWrap(true)
	textView.SetText(text)

	tw, th := textDimensions(text)

	return &textDialog{
		text:       text,
		textWidth:  tw,
		textHeight: th,
		textView:   textView,
	}
}

func (td *textDialog) SetText(text string) {
	td.textView.SetText(text)
	td.textWidth, td.textHeight = textDimensions(text)
}

func (td *textDialog) Draw(screen tcell.Screen) {
	screenWidth, screenHeight := screen.Size()

	w, h := screenWidth/2, screenHeight/2

	// We add 2 below to account for the text view's border.
	textWidth, textHeight := td.textWidth+2, td.textHeight+2

	if w > textWidth {
		w = textWidth
	} else {
		// We're going to be wrapping the text. Recalculate the text height with word wrap.
		textHeight = len(cui.WordWrap(td.text, w)) + 2
	}
	if h > textHeight {
		h = textHeight
	}

	x, y := (screenWidth-w)/2, (screenHeight-h)/2
	td.textView.SetRect(x, y, w, h)
	td.textView.Draw(screen)
}

func (td *textDialog) GetRect() (int, int, int, int) {
	return td.x, td.y, td.w, td.h
}

func (td *textDialog) SetRect(x, y, w, h int) {
	td.x, td.y, td.w, td.h = x, y, w, h
}

func (td *textDialog) InputHandler() func(event *tcell.EventKey, setFocus func(p cui.Primitive)) {
	return td.textView.InputHandler()
}

func (td *textDialog) MouseHandler() func(action cui.MouseAction, event *tcell.EventMouse, setFocus func(p cui.Primitive)) (consumed bool, capture cui.Primitive) {
	return td.textView.MouseHandler()
}

func (td *textDialog) Focus(delegate func(p cui.Primitive)) {
	td.textView.Focus(delegate)
}

func (td *textDialog) Blur() {
	td.textView.Blur()
}

func (td *textDialog) HasFocus() bool {
	return td.textView.HasFocus()
}

func getDocumentAnchor(url string) (string, bool) {
	if !strings.HasPrefix(url, "#") {
		return "", false
	}
	return url[1:], true
}

func sendToClipboard(value string) {
	if !clipboard.Unsupported {
		clipboard.WriteAll(value)
	}
}

type Panel struct {
	view *MdView

	app *cui.Application

	hasFocus      bool
	focused       cui.Primitive
	lastFocused   cui.Primitive
	inputHandler  func(event *tcell.EventKey, setFocus func(p cui.Primitive))
	visibleDialog *textDialog

	helpDialog *textDialog
	rootPages  *cui.Pages

	backstack []*renderer.NodeSpan
	query     *regexp.Regexp
}

func (r *Panel) GetVisible() bool {
	return r.view.GetVisible()
}

func (r *Panel) SetVisible(v bool) {
	r.view.SetVisible(v)
}

func (r *Panel) GetFocusable() cui.Focusable {
	return r.view
}

func New(name, source string, theme *chroma.Style, app *cui.Application) *Panel {
	r := &Panel{
		view:       NewMdView(theme),
		app:        app,
		helpDialog: newTextDialog(helpText, "Help"),
	}

	r.view.SetText(name, source)
	r.view.SetGutter(true)

	rootPages := cui.NewPages()
	rootPages.AddAndSwitchToPage("markdown", r.view, true)
	rootPages.AddPage("help", r.helpDialog, true, false)
	r.rootPages = rootPages

	r.focused = r.view

	return r
}

func (r *Panel) Draw(screen tcell.Screen) {
	r.rootPages.Draw(screen)
}

func (r *Panel) GetRect() (int, int, int, int) {
	return r.rootPages.GetRect()
}

func (r *Panel) SetRect(xc, yc, width, height int) {
	r.rootPages.SetRect(xc, yc, width, height)
}

func (r *Panel) focusedLink() string {
	if span := r.view.Selection(); span != nil {
		switch node := span.Node.(type) {
		case *ast.AutoLink:
			return string(node.URL(r.view.GetMarkdown()))
		case *ast.Link:
			return string(node.Destination)
		}
	}
	return ""
}

func (r *Panel) InputHandler() func(event *tcell.EventKey, setFocus func(p cui.Primitive)) {
	return func(event *tcell.EventKey, setFocus func(p cui.Primitive)) {
		event = func() *tcell.EventKey {
			if r.visibleDialog != nil {
				if event.Key() == tcell.KeyEscape || event.Rune() == 'h' && r.visibleDialog == r.helpDialog {
					r.hideDialog()
					return nil
				}
				return event
			}

			switch event.Key() {
			case tcell.KeyCtrlO:
				link := r.focusedLink()
				if anchor, ok := getDocumentAnchor(link); ok {
					selection := r.view.Selection()
					if r.view.SelectAnchor(anchor) && selection != nil {
						r.backstack = append(r.backstack, selection)
					}
				} else {
					//if err := openInBrowser(link); err != nil {
					//	r.showErrorDialog("opening issue", err)
					//}
				}
			case tcell.KeyRune:
				switch event.Rune() {
				case '<':
					if len(r.backstack) != 0 {
						last := r.backstack[len(r.backstack)-1]
						r.backstack = r.backstack[:len(r.backstack)-1]
						r.view.SelectSpan(last, true)
					}
				case 'h':
					// Show the help
					r.showDialog(r.helpDialog)
					return nil
				}
			}
			return event
		}()
		if event != nil && r.focused != nil {
			if handler := r.focused.InputHandler(); handler != nil {
				handler(event, r.setFocus)
			}
		}
	}
}

func (r *Panel) MouseHandler() func(action cui.MouseAction, event *tcell.EventMouse, setFocus func(p cui.Primitive)) (consumed bool, capture cui.Primitive) {
	return r.rootPages.MouseHandler()
}

func (r *Panel) Focus(delegate func(p cui.Primitive)) {
	r.hasFocus = true
	if r.focused != nil {
		r.focused.Focus(r.setFocus)
	}
}

func (r *Panel) Blur() {
	r.hasFocus = false
	if r.focused != nil {
		r.focused.Blur()
	}
}

func (r *Panel) HasFocus() bool {
	return r.hasFocus
}

func (r *Panel) setFocus(p cui.Primitive) {
	r.lastFocused = r.focused

	var doSetFocus func(p cui.Primitive)
	doSetFocus = func(p cui.Primitive) {
		if r.focused != nil {
			r.focused.Blur()
		}

		r.focused = p
		r.inputHandler = p.InputHandler()
		r.focused.Focus(doSetFocus)
	}

	doSetFocus(p)
}

func (r *Panel) showDialog(d *textDialog) {
	if r.visibleDialog != nil && r.visibleDialog.persistent {
		return
	}

	r.rootPages.AddPage("dialog", d, true, true)
	r.setFocus(d)
	r.visibleDialog = d
}

func (r *Panel) hideDialog() {
	if r.visibleDialog != nil && r.visibleDialog.persistent {
		return
	}

	r.rootPages.HidePage("dialog")
	r.setFocus(r.lastFocused)
	r.visibleDialog = nil
}

func (r *Panel) showErrorDialog(action string, err error) {
	r.showDialog(newTextDialog(fmt.Sprintf("Error %v: %v", action, err.Error()), "Error"))
}
