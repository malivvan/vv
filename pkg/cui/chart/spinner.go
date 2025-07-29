package chart

import (
	"github.com/gdamore/tcell/v2"
	"github.com/malivvan/vv/pkg/cui"
)

// Spinner represents a spinner widget.
type Spinner struct {
	*cui.Box

	counter      int
	currentStyle SpinnerStyle

	styles map[SpinnerStyle][]rune
}

func (s *Spinner) GetVisible() bool {
	return true
}

func (s *Spinner) SetVisible(v bool) {

}

func (s *Spinner) InputHandler() func(event *tcell.EventKey, setFocus func(p cui.Primitive)) {
	return func(event *tcell.EventKey, setFocus func(p cui.Primitive)) {
		// No input handling for spinner
	}
}

func (s *Spinner) Focus(delegate func(p cui.Primitive)) {
	// Spinner does not take focus, but we implement this method to satisfy the Focusable interface.
	// If you want to use the spinner in a focusable context, you can delegate focus to another primitive.
	delegate(s)
}

func (s *Spinner) GetFocusable() cui.Focusable {
	return s
}

func (s *Spinner) MouseHandler() func(action cui.MouseAction, event *tcell.EventMouse, setFocus func(p cui.Primitive)) (consumed bool, capture cui.Primitive) {
	return func(action cui.MouseAction, event *tcell.EventMouse, setFocus func(p cui.Primitive)) (consumed bool, capture cui.Primitive) {
		// No mouse handling for spinner
		return false, nil
	}
}

type SpinnerStyle int

const (
	SpinnerDotsCircling SpinnerStyle = iota
	SpinnerDotsUpDown
	SpinnerBounce
	SpinnerLine
	SpinnerCircleQuarters
	SpinnerSquareCorners
	SpinnerCircleHalves
	SpinnerCorners
	SpinnerArrows
	SpinnerHamburger
	SpinnerStack
	SpinnerGrowHorizontal
	SpinnerGrowVertical
	SpinnerStar
	SpinnerBoxBounce
	spinnerCustom // non-public constant to indicate that a custom style has been set by the user.
)

// NewSpinner returns a new spinner widget.
func NewSpinner() *Spinner {
	return &Spinner{
		Box:          cui.NewBox(),
		currentStyle: SpinnerDotsCircling,
		styles: map[SpinnerStyle][]rune{
			SpinnerDotsCircling:   []rune(`⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏`),
			SpinnerDotsUpDown:     []rune(`⠋⠙⠚⠞⠖⠦⠴⠲⠳⠓`),
			SpinnerBounce:         []rune(`⠄⠆⠇⠋⠙⠸⠰⠠⠰⠸⠙⠋⠇⠆`),
			SpinnerLine:           []rune(`|/-\`),
			SpinnerCircleQuarters: []rune(`◴◷◶◵`),
			SpinnerSquareCorners:  []rune(`◰◳◲◱`),
			SpinnerCircleHalves:   []rune(`◐◓◑◒`),
			SpinnerCorners:        []rune(`⌜⌝⌟⌞`),
			SpinnerArrows:         []rune(`⇑⇗⇒⇘⇓⇙⇐⇖`),
			SpinnerHamburger:      []rune(`☰☱☳☷☶☴`),
			SpinnerStack:          []rune(`䷀䷪䷡䷊䷒䷗䷁䷖䷓䷋䷠䷫`),
			SpinnerGrowHorizontal: []rune(`▉▊▋▌▍▎▏▎▍▌▋▊▉`),
			SpinnerGrowVertical:   []rune(`▁▃▄▅▆▇▆▅▄▃`),
			SpinnerStar:           []rune(`✶✸✹✺✹✷`),
			SpinnerBoxBounce:      []rune(`▌▀▐▄`),
		},
	}
}

// Draw draws this primitive onto the screen.
func (s *Spinner) Draw(screen tcell.Screen) {
	s.Box.Draw(screen)
	x, y, width, _ := s.Box.GetInnerRect()
	cui.Print(screen, []byte(s.getCurrentFrame()), x, y, width, cui.AlignLeft, tcell.ColorDefault)
}

// Pulse updates the spinner to the next frame.
func (s *Spinner) Pulse() {
	s.counter++
}

// Reset sets the frame counter to 0.
func (s *Spinner) Reset() {
	s.counter = 0
}

// SetStyle sets the spinner style.
func (s *Spinner) SetStyle(style SpinnerStyle) *Spinner {
	s.currentStyle = style

	return s
}

func (s *Spinner) getCurrentFrame() string {
	frames := s.styles[s.currentStyle]
	if len(frames) == 0 {
		return ""
	}

	return string(frames[s.counter%len(frames)])
}

// SetCustomStyle sets a list of runes as custom frames to show as the spinner.
func (s *Spinner) SetCustomStyle(frames []rune) *Spinner {
	s.styles[spinnerCustom] = frames
	s.currentStyle = spinnerCustom

	return s
}
