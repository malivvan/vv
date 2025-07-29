package ansicsi

import "io"

const (
	SGRReset                   = 0  // Default rendition (implementation-defined), cancels the effect of any preceding occurrence of SGR
	SGRBold                    = 1  // Bold or increased intensity
	SGRFaint                   = 2  // Faint, decreased intensity or second color
	SGRItalic                  = 3  // Italicized
	SGRUnderline               = 4  // Singly underlined
	SGRSlowBlink               = 5  // Slowly blinking (less then 150 per minute)
	SGRRapidBlink              = 6  // Rapidly blinking (150 per minute or more)
	SGRInverse                 = 7  // Swap foreground and background colors
	SGRConceal                 = 8  // Concealed characters
	SGRStrikethrough           = 9  // Crossed-out (characters still legible but marked as to be deleted)
	SGRDefaultFont             = 10 // Primary (default) font
	SGRAlternativeFont1        = 11 // First alternative font
	SGRAlternativeFont2        = 12 // Second alternative font
	SGRAlternativeFont3        = 13 // Third alternative font
	SGRAlternativeFont4        = 14 // Fourth alternative font
	SGRAlternativeFont5        = 15 // Fifth alternative font
	SGRAlternativeFont6        = 16 // Sixth alternative font
	SGRAlternativeFont7        = 17 // Seventh alternative font
	SGRAlternativeFont8        = 18 // Eighth alternative font
	SGRAlternativeFont9        = 19 // Ninth alternative font
	SGRFraktur                 = 20 // Fraktur (Gothic)
	SGRDoubleUnderline         = 21 // Doubly underlined
	SGRNormalWeight            = 22 // Normal color or normal intensity (neither bold nor faint)
	SGRNoItalicOrFraktur       = 23 // Not italicized, not fraktur
	SGRNoUnderline             = 24 // Not underlined (neither singly nor doubly)
	SGRNoBlink                 = 25 // Steady (not blinking)
	SGRProportionalSpacing     = 26 // Proportional spacing as specified in CCITT Recommendation T.61
	SGRNoInverse               = 27 // Disable foreground and background color swap
	SGRNoConceal               = 28 // Revealed characters
	SGRNoStrikethrough         = 29 // Not crossed out
	SGRForegroundBlack         = 30 // Black foreground color
	SGRForegroundRed           = 31 // Red foreground color
	SGRForegroundGreen         = 32 // Green foreground color
	SGRForegroundYellow        = 33 // Yellow foreground color
	SGRForegroundBlue          = 34 // Blue foreground color
	SGRForegroundMagenta       = 35 // Magenta foreground color
	SGRForegroundCyan          = 36 // Cyan foreground color
	SGRForegroundWhite         = 37 // White foreground color
	SGRForegroundColor         = 38 // Set the foreground color to a 256-color or true color value
	SGRForegroundDefault       = 39 // Default foreground color (implementation-defined)
	SGRBackgroundBlack         = 40 // Black background color
	SGRBackgroundRed           = 41 // Red background color
	SGRBackgroundGreen         = 42 // Green background color
	SGRBackgroundYellow        = 43 // Yellow background color
	SGRBackgroundBlue          = 44 // Blue background color
	SGRBackgroundMagenta       = 45 // Magenta background color
	SGRBackgroundCyan          = 46 // Cyan background color
	SGRBackgroundWhite         = 47 // White background color
	SGRBackgroundColor         = 48 // Set the background color to a 256-color or true color value
	SGRBackgroundDefault       = 49 // Default background color (implementation-defined)
	SGRNoProportionalSpacing   = 50 // Doisable proportional spacing
	SGRFrame                   = 51 // Framed
	SGREncircle                = 52 // Encircled
	SGROverline                = 53 // Overlined
	SGRNoFrameOrEncircle       = 54 // Not framed or encircled
	SGRNoOverline              = 55 // Not overlined
	SGRUnderlineColor          = 58 // Set the underline color to a 256-color or true color value (non-standard)
	SGRDefaultUnderlineColor   = 59 // Default underline color
	SGRIdeogramUnderline       = 60 // Ideogram underline or right side line
	SGRIdeogramDoubleUnderline = 61 // Ideogram double underline or double line on the right side
	SGRIdeogramOverline        = 62 // Ideogram overline or left side line
	SGRIdeogramDoubleOverline  = 63 // Ideogram double overline or double line on the left side
	SGRIdeogramStress          = 64 // Ideogram stress marking
	SGRIdeogramReset           = 65 // Cancels the effect of the rendition aspects established by ideogram parameter values
)

// SetGraphicsRendition represents a single Set Graphics Rendition control function.
type SetGraphicsRendition struct {
	// Command describes the graphics rendition aspect that this call affects.
	Command int
	// Parameters are the parameters (if any) to the command.
	Parameters []int
}

func (sgr *SetGraphicsRendition) Encode(w io.Writer) (int, error) {
	params := make([]int, 0, len(sgr.Parameters)+1)
	params = append(params, sgr.Command)
	params = append(params, sgr.Parameters...)
	return encodeCommand(w, params, nil, 0x6d)
}

func (sgr *SetGraphicsRendition) decodeParameters(params []int) bool {
	if len(params) < 1 {
		return true
	}

	command, params := params[0], params[1:]
	switch command {
	case SGRForegroundColor, SGRBackgroundColor, SGRUnderlineColor:
		if len(params) < 1 {
			return false
		}

		depth := params[0]
		switch depth {
		case 2:
			if len(params) != 4 {
				return false
			}
		case 5:
			if len(params) != 2 {
				return false
			}
		default:
			return false
		}
	default:
		if command < 0 || command > 65 {
			return false
		}
	}

	sgr.Command, sgr.Parameters = command, params
	return true
}
