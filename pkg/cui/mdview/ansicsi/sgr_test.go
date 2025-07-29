package ansicsi

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Default rendition (implementation-defined), cancels the effect of any preceding occurrence of SGR
func TestSGR_Reset(t *testing.T) {
	command := []byte("\x1b[0m")
	cmd, size := Decode(command)
	sgr, ok := cmd.(*SetGraphicsRendition)
	assert.True(t, ok)
	assert.Equal(t, 4, size)
	assert.Equal(t, SGRReset, sgr.Command)

	var b bytes.Buffer
	encodedSize, err := sgr.Encode(&b)
	assert.NoError(t, err)
	assert.Equal(t, size, encodedSize)
	assert.Equal(t, command, b.Bytes())
}

// bold or increased intensity
func TestSGR_Bold(t *testing.T) {
	command := []byte("\x1b[1m")
	cmd, size := Decode(command)
	sgr, ok := cmd.(*SetGraphicsRendition)
	assert.True(t, ok)
	assert.Equal(t, 4, size)
	assert.Equal(t, SGRBold, sgr.Command)

	var b bytes.Buffer
	encodedSize, err := sgr.Encode(&b)
	assert.NoError(t, err)
	assert.Equal(t, size, encodedSize)
	assert.Equal(t, command, b.Bytes())
}

// faint, decreased intensity or second color
func TestSGR_Faint(t *testing.T) {
	command := []byte("\x1b[2m")
	cmd, size := Decode(command)
	sgr, ok := cmd.(*SetGraphicsRendition)
	assert.True(t, ok)
	assert.Equal(t, 4, size)
	assert.Equal(t, SGRFaint, sgr.Command)

	var b bytes.Buffer
	encodedSize, err := sgr.Encode(&b)
	assert.NoError(t, err)
	assert.Equal(t, size, encodedSize)
	assert.Equal(t, command, b.Bytes())
}

// italicized
func TestSGR_Italic(t *testing.T) {
	command := []byte("\x1b[3m")
	cmd, size := Decode(command)
	sgr, ok := cmd.(*SetGraphicsRendition)
	assert.True(t, ok)
	assert.Equal(t, 4, size)
	assert.Equal(t, SGRItalic, sgr.Command)

	var b bytes.Buffer
	encodedSize, err := sgr.Encode(&b)
	assert.NoError(t, err)
	assert.Equal(t, size, encodedSize)
	assert.Equal(t, command, b.Bytes())
}

// singly underlined
func TestSGR_Underline(t *testing.T) {
	command := []byte("\x1b[4m")
	cmd, size := Decode(command)
	sgr, ok := cmd.(*SetGraphicsRendition)
	assert.True(t, ok)
	assert.Equal(t, 4, size)
	assert.Equal(t, SGRUnderline, sgr.Command)

	var b bytes.Buffer
	encodedSize, err := sgr.Encode(&b)
	assert.NoError(t, err)
	assert.Equal(t, size, encodedSize)
	assert.Equal(t, command, b.Bytes())
}

// slowly blinking (less then 150 per minute)
func TestSGR_SlowBlink(t *testing.T) {
	command := []byte("\x1b[5m")
	cmd, size := Decode(command)
	sgr, ok := cmd.(*SetGraphicsRendition)
	assert.True(t, ok)
	assert.Equal(t, 4, size)
	assert.Equal(t, SGRSlowBlink, sgr.Command)

	var b bytes.Buffer
	encodedSize, err := sgr.Encode(&b)
	assert.NoError(t, err)
	assert.Equal(t, size, encodedSize)
	assert.Equal(t, command, b.Bytes())
}

// rapidly blinking (150 per minute or more)
func TestSGR_RapidBlink(t *testing.T) {
	command := []byte("\x1b[6m")
	cmd, size := Decode(command)
	sgr, ok := cmd.(*SetGraphicsRendition)
	assert.True(t, ok)
	assert.Equal(t, 4, size)
	assert.Equal(t, SGRRapidBlink, sgr.Command)

	var b bytes.Buffer
	encodedSize, err := sgr.Encode(&b)
	assert.NoError(t, err)
	assert.Equal(t, size, encodedSize)
	assert.Equal(t, command, b.Bytes())
}

// swap foreground and background colors
func TestSGR_Inverse(t *testing.T) {
	command := []byte("\x1b[7m")
	cmd, size := Decode(command)
	sgr, ok := cmd.(*SetGraphicsRendition)
	assert.True(t, ok)
	assert.Equal(t, 4, size)
	assert.Equal(t, SGRInverse, sgr.Command)

	var b bytes.Buffer
	encodedSize, err := sgr.Encode(&b)
	assert.NoError(t, err)
	assert.Equal(t, size, encodedSize)
	assert.Equal(t, command, b.Bytes())
}

// concealed characters
func TestSGR_Conceal(t *testing.T) {
	command := []byte("\x1b[8m")
	cmd, size := Decode(command)
	sgr, ok := cmd.(*SetGraphicsRendition)
	assert.True(t, ok)
	assert.Equal(t, 4, size)
	assert.Equal(t, SGRConceal, sgr.Command)

	var b bytes.Buffer
	encodedSize, err := sgr.Encode(&b)
	assert.NoError(t, err)
	assert.Equal(t, size, encodedSize)
	assert.Equal(t, command, b.Bytes())
}

// crossed-out (characters still legible but marked as to be deleted)
func TestSGR_Strikethrough(t *testing.T) {
	command := []byte("\x1b[9m")
	cmd, size := Decode(command)
	sgr, ok := cmd.(*SetGraphicsRendition)
	assert.True(t, ok)
	assert.Equal(t, 4, size)
	assert.Equal(t, SGRStrikethrough, sgr.Command)

	var b bytes.Buffer
	encodedSize, err := sgr.Encode(&b)
	assert.NoError(t, err)
	assert.Equal(t, size, encodedSize)
	assert.Equal(t, command, b.Bytes())
}

// primary (default) font
func TestSGR_DefaultFont(t *testing.T) {
	command := []byte("\x1b[10m")
	cmd, size := Decode(command)
	sgr, ok := cmd.(*SetGraphicsRendition)
	assert.True(t, ok)
	assert.Equal(t, 5, size)
	assert.Equal(t, SGRDefaultFont, sgr.Command)

	var b bytes.Buffer
	encodedSize, err := sgr.Encode(&b)
	assert.NoError(t, err)
	assert.Equal(t, size, encodedSize)
	assert.Equal(t, command, b.Bytes())
}

// first alternative font
func TestSGR_AlternativeFont1(t *testing.T) {
	command := []byte("\x1b[11m")
	cmd, size := Decode(command)
	sgr, ok := cmd.(*SetGraphicsRendition)
	assert.True(t, ok)
	assert.Equal(t, 5, size)
	assert.Equal(t, SGRAlternativeFont1, sgr.Command)

	var b bytes.Buffer
	encodedSize, err := sgr.Encode(&b)
	assert.NoError(t, err)
	assert.Equal(t, size, encodedSize)
	assert.Equal(t, command, b.Bytes())
}

// second alternative font
func TestSGR_AlternativeFont2(t *testing.T) {
	command := []byte("\x1b[12m")
	cmd, size := Decode(command)
	sgr, ok := cmd.(*SetGraphicsRendition)
	assert.True(t, ok)
	assert.Equal(t, 5, size)
	assert.Equal(t, SGRAlternativeFont2, sgr.Command)

	var b bytes.Buffer
	encodedSize, err := sgr.Encode(&b)
	assert.NoError(t, err)
	assert.Equal(t, size, encodedSize)
	assert.Equal(t, command, b.Bytes())
}

// third alternative font
func TestSGR_AlternativeFont3(t *testing.T) {
	command := []byte("\x1b[13m")
	cmd, size := Decode(command)
	sgr, ok := cmd.(*SetGraphicsRendition)
	assert.True(t, ok)
	assert.Equal(t, 5, size)
	assert.Equal(t, SGRAlternativeFont3, sgr.Command)

	var b bytes.Buffer
	encodedSize, err := sgr.Encode(&b)
	assert.NoError(t, err)
	assert.Equal(t, size, encodedSize)
	assert.Equal(t, command, b.Bytes())
}

// fourth alternative font
func TestSGR_AlternativeFont4(t *testing.T) {
	command := []byte("\x1b[14m")
	cmd, size := Decode(command)
	sgr, ok := cmd.(*SetGraphicsRendition)
	assert.True(t, ok)
	assert.Equal(t, 5, size)
	assert.Equal(t, SGRAlternativeFont4, sgr.Command)

	var b bytes.Buffer
	encodedSize, err := sgr.Encode(&b)
	assert.NoError(t, err)
	assert.Equal(t, size, encodedSize)
	assert.Equal(t, command, b.Bytes())
}

// fifth alternative font
func TestSGR_AlternativeFont5(t *testing.T) {
	command := []byte("\x1b[15m")
	cmd, size := Decode(command)
	sgr, ok := cmd.(*SetGraphicsRendition)
	assert.True(t, ok)
	assert.Equal(t, 5, size)
	assert.Equal(t, SGRAlternativeFont5, sgr.Command)

	var b bytes.Buffer
	encodedSize, err := sgr.Encode(&b)
	assert.NoError(t, err)
	assert.Equal(t, size, encodedSize)
	assert.Equal(t, command, b.Bytes())
}

// sixth alternative font
func TestSGR_AlternativeFont6(t *testing.T) {
	command := []byte("\x1b[16m")
	cmd, size := Decode(command)
	sgr, ok := cmd.(*SetGraphicsRendition)
	assert.True(t, ok)
	assert.Equal(t, 5, size)
	assert.Equal(t, SGRAlternativeFont6, sgr.Command)

	var b bytes.Buffer
	encodedSize, err := sgr.Encode(&b)
	assert.NoError(t, err)
	assert.Equal(t, size, encodedSize)
	assert.Equal(t, command, b.Bytes())
}

// seventh alternative font
func TestSGR_AlternativeFont7(t *testing.T) {
	command := []byte("\x1b[17m")
	cmd, size := Decode(command)
	sgr, ok := cmd.(*SetGraphicsRendition)
	assert.True(t, ok)
	assert.Equal(t, 5, size)
	assert.Equal(t, SGRAlternativeFont7, sgr.Command)

	var b bytes.Buffer
	encodedSize, err := sgr.Encode(&b)
	assert.NoError(t, err)
	assert.Equal(t, size, encodedSize)
	assert.Equal(t, command, b.Bytes())
}

// eighth alternative font
func TestSGR_AlternativeFont8(t *testing.T) {
	command := []byte("\x1b[18m")
	cmd, size := Decode(command)
	sgr, ok := cmd.(*SetGraphicsRendition)
	assert.True(t, ok)
	assert.Equal(t, 5, size)
	assert.Equal(t, SGRAlternativeFont8, sgr.Command)

	var b bytes.Buffer
	encodedSize, err := sgr.Encode(&b)
	assert.NoError(t, err)
	assert.Equal(t, size, encodedSize)
	assert.Equal(t, command, b.Bytes())
}

// ninth alternative font
func TestSGR_AlternativeFont9(t *testing.T) {
	command := []byte("\x1b[19m")
	cmd, size := Decode(command)
	sgr, ok := cmd.(*SetGraphicsRendition)
	assert.True(t, ok)
	assert.Equal(t, 5, size)
	assert.Equal(t, SGRAlternativeFont9, sgr.Command)

	var b bytes.Buffer
	encodedSize, err := sgr.Encode(&b)
	assert.NoError(t, err)
	assert.Equal(t, size, encodedSize)
	assert.Equal(t, command, b.Bytes())
}

// Fraktur (Gothic)
func TestSGR_Fraktur(t *testing.T) {
	command := []byte("\x1b[20m")
	cmd, size := Decode(command)
	sgr, ok := cmd.(*SetGraphicsRendition)
	assert.True(t, ok)
	assert.Equal(t, 5, size)
	assert.Equal(t, SGRFraktur, sgr.Command)

	var b bytes.Buffer
	encodedSize, err := sgr.Encode(&b)
	assert.NoError(t, err)
	assert.Equal(t, size, encodedSize)
	assert.Equal(t, command, b.Bytes())
}

// doubly underlined
func TestSGR_DoubleUnderline(t *testing.T) {
	command := []byte("\x1b[21m")
	cmd, size := Decode(command)
	sgr, ok := cmd.(*SetGraphicsRendition)
	assert.True(t, ok)
	assert.Equal(t, 5, size)
	assert.Equal(t, SGRDoubleUnderline, sgr.Command)

	var b bytes.Buffer
	encodedSize, err := sgr.Encode(&b)
	assert.NoError(t, err)
	assert.Equal(t, size, encodedSize)
	assert.Equal(t, command, b.Bytes())
}

// normal color or normal intensity (neither bold nor faint)
func TestSGR_NormalWeight(t *testing.T) {
	command := []byte("\x1b[22m")
	cmd, size := Decode(command)
	sgr, ok := cmd.(*SetGraphicsRendition)
	assert.True(t, ok)
	assert.Equal(t, 5, size)
	assert.Equal(t, SGRNormalWeight, sgr.Command)

	var b bytes.Buffer
	encodedSize, err := sgr.Encode(&b)
	assert.NoError(t, err)
	assert.Equal(t, size, encodedSize)
	assert.Equal(t, command, b.Bytes())
}

// not italicized, not fraktur
func TestSGR_NoItalicOrFraktur(t *testing.T) {
	command := []byte("\x1b[23m")
	cmd, size := Decode(command)
	sgr, ok := cmd.(*SetGraphicsRendition)
	assert.True(t, ok)
	assert.Equal(t, 5, size)
	assert.Equal(t, SGRNoItalicOrFraktur, sgr.Command)

	var b bytes.Buffer
	encodedSize, err := sgr.Encode(&b)
	assert.NoError(t, err)
	assert.Equal(t, size, encodedSize)
	assert.Equal(t, command, b.Bytes())
}

// not underlined (neither singly nor doubly)
func TestSGR_NoUnderline(t *testing.T) {
	command := []byte("\x1b[24m")
	cmd, size := Decode(command)
	sgr, ok := cmd.(*SetGraphicsRendition)
	assert.True(t, ok)
	assert.Equal(t, 5, size)
	assert.Equal(t, SGRNoUnderline, sgr.Command)

	var b bytes.Buffer
	encodedSize, err := sgr.Encode(&b)
	assert.NoError(t, err)
	assert.Equal(t, size, encodedSize)
	assert.Equal(t, command, b.Bytes())
}

// steady (not blinking)
func TestSGR_NoBlink(t *testing.T) {
	command := []byte("\x1b[25m")
	cmd, size := Decode(command)
	sgr, ok := cmd.(*SetGraphicsRendition)
	assert.True(t, ok)
	assert.Equal(t, 5, size)
	assert.Equal(t, SGRNoBlink, sgr.Command)

	var b bytes.Buffer
	encodedSize, err := sgr.Encode(&b)
	assert.NoError(t, err)
	assert.Equal(t, size, encodedSize)
	assert.Equal(t, command, b.Bytes())
}

// proportional spacing as specified in CCITT Recommendation T.61
func TestSGR_ProportionalSpacing(t *testing.T) {
	command := []byte("\x1b[26m")
	cmd, size := Decode(command)
	sgr, ok := cmd.(*SetGraphicsRendition)
	assert.True(t, ok)
	assert.Equal(t, 5, size)
	assert.Equal(t, SGRProportionalSpacing, sgr.Command)

	var b bytes.Buffer
	encodedSize, err := sgr.Encode(&b)
	assert.NoError(t, err)
	assert.Equal(t, size, encodedSize)
	assert.Equal(t, command, b.Bytes())
}

// disable foreground and background color swap
func TestSGR_NoInverse(t *testing.T) {
	command := []byte("\x1b[27m")
	cmd, size := Decode(command)
	sgr, ok := cmd.(*SetGraphicsRendition)
	assert.True(t, ok)
	assert.Equal(t, 5, size)
	assert.Equal(t, SGRNoInverse, sgr.Command)

	var b bytes.Buffer
	encodedSize, err := sgr.Encode(&b)
	assert.NoError(t, err)
	assert.Equal(t, size, encodedSize)
	assert.Equal(t, command, b.Bytes())
}

// revealed characters
func TestSGR_NoConceal(t *testing.T) {
	command := []byte("\x1b[28m")
	cmd, size := Decode(command)
	sgr, ok := cmd.(*SetGraphicsRendition)
	assert.True(t, ok)
	assert.Equal(t, 5, size)
	assert.Equal(t, SGRNoConceal, sgr.Command)

	var b bytes.Buffer
	encodedSize, err := sgr.Encode(&b)
	assert.NoError(t, err)
	assert.Equal(t, size, encodedSize)
	assert.Equal(t, command, b.Bytes())
}

// not crossed out
func TestSGR_NoStrikethrough(t *testing.T) {
	command := []byte("\x1b[29m")
	cmd, size := Decode(command)
	sgr, ok := cmd.(*SetGraphicsRendition)
	assert.True(t, ok)
	assert.Equal(t, 5, size)
	assert.Equal(t, SGRNoStrikethrough, sgr.Command)

	var b bytes.Buffer
	encodedSize, err := sgr.Encode(&b)
	assert.NoError(t, err)
	assert.Equal(t, size, encodedSize)
	assert.Equal(t, command, b.Bytes())
}

// black foreground color
func TestSGR_ForegroundBlack(t *testing.T) {
	command := []byte("\x1b[30m")
	cmd, size := Decode(command)
	sgr, ok := cmd.(*SetGraphicsRendition)
	assert.True(t, ok)
	assert.Equal(t, 5, size)
	assert.Equal(t, SGRForegroundBlack, sgr.Command)

	var b bytes.Buffer
	encodedSize, err := sgr.Encode(&b)
	assert.NoError(t, err)
	assert.Equal(t, size, encodedSize)
	assert.Equal(t, command, b.Bytes())
}

// red foreground color
func TestSGR_ForegroundRed(t *testing.T) {
	command := []byte("\x1b[31m")
	cmd, size := Decode(command)
	sgr, ok := cmd.(*SetGraphicsRendition)
	assert.True(t, ok)
	assert.Equal(t, 5, size)
	assert.Equal(t, SGRForegroundRed, sgr.Command)

	var b bytes.Buffer
	encodedSize, err := sgr.Encode(&b)
	assert.NoError(t, err)
	assert.Equal(t, size, encodedSize)
	assert.Equal(t, command, b.Bytes())
}

// green foreground color
func TestSGR_ForegroundGreen(t *testing.T) {
	command := []byte("\x1b[32m")
	cmd, size := Decode(command)
	sgr, ok := cmd.(*SetGraphicsRendition)
	assert.True(t, ok)
	assert.Equal(t, 5, size)
	assert.Equal(t, SGRForegroundGreen, sgr.Command)

	var b bytes.Buffer
	encodedSize, err := sgr.Encode(&b)
	assert.NoError(t, err)
	assert.Equal(t, size, encodedSize)
	assert.Equal(t, command, b.Bytes())
}

// yellow foreground color
func TestSGR_ForegroundYellow(t *testing.T) {
	command := []byte("\x1b[33m")
	cmd, size := Decode(command)
	sgr, ok := cmd.(*SetGraphicsRendition)
	assert.True(t, ok)
	assert.Equal(t, 5, size)
	assert.Equal(t, SGRForegroundYellow, sgr.Command)

	var b bytes.Buffer
	encodedSize, err := sgr.Encode(&b)
	assert.NoError(t, err)
	assert.Equal(t, size, encodedSize)
	assert.Equal(t, command, b.Bytes())
}

// blue foreground color
func TestSGR_ForegroundBlue(t *testing.T) {
	command := []byte("\x1b[34m")
	cmd, size := Decode(command)
	sgr, ok := cmd.(*SetGraphicsRendition)
	assert.True(t, ok)
	assert.Equal(t, 5, size)
	assert.Equal(t, SGRForegroundBlue, sgr.Command)

	var b bytes.Buffer
	encodedSize, err := sgr.Encode(&b)
	assert.NoError(t, err)
	assert.Equal(t, size, encodedSize)
	assert.Equal(t, command, b.Bytes())
}

// magenta foreground color
func TestSGR_ForegroundMagenta(t *testing.T) {
	command := []byte("\x1b[35m")
	cmd, size := Decode(command)
	sgr, ok := cmd.(*SetGraphicsRendition)
	assert.True(t, ok)
	assert.Equal(t, 5, size)
	assert.Equal(t, SGRForegroundMagenta, sgr.Command)

	var b bytes.Buffer
	encodedSize, err := sgr.Encode(&b)
	assert.NoError(t, err)
	assert.Equal(t, size, encodedSize)
	assert.Equal(t, command, b.Bytes())
}

// cyan foreground color
func TestSGR_ForegroundCyan(t *testing.T) {
	command := []byte("\x1b[36m")
	cmd, size := Decode(command)
	sgr, ok := cmd.(*SetGraphicsRendition)
	assert.True(t, ok)
	assert.Equal(t, 5, size)
	assert.Equal(t, SGRForegroundCyan, sgr.Command)

	var b bytes.Buffer
	encodedSize, err := sgr.Encode(&b)
	assert.NoError(t, err)
	assert.Equal(t, size, encodedSize)
	assert.Equal(t, command, b.Bytes())
}

// white foreground color
func TestSGR_ForegroundWhite(t *testing.T) {
	command := []byte("\x1b[37m")
	cmd, size := Decode(command)
	sgr, ok := cmd.(*SetGraphicsRendition)
	assert.True(t, ok)
	assert.Equal(t, 5, size)
	assert.Equal(t, SGRForegroundWhite, sgr.Command)

	var b bytes.Buffer
	encodedSize, err := sgr.Encode(&b)
	assert.NoError(t, err)
	assert.Equal(t, size, encodedSize)
	assert.Equal(t, command, b.Bytes())
}

// set the foreground color to a 256-color or true color value
func TestSGR_ForegroundColor_256(t *testing.T) {
	command := []byte("\x1b[38;5;128m")
	cmd, size := Decode(command)
	sgr, ok := cmd.(*SetGraphicsRendition)
	assert.True(t, ok)
	assert.Equal(t, 11, size)
	assert.Equal(t, SGRForegroundColor, sgr.Command)
	assert.Equal(t, []int{5, 128}, sgr.Parameters)

	var b bytes.Buffer
	encodedSize, err := sgr.Encode(&b)
	assert.NoError(t, err)
	assert.Equal(t, size, encodedSize)
	assert.Equal(t, command, b.Bytes())
}

// set the foreground color to a 256-color or true color value
func TestSGR_ForegroundColor_Truecolor(t *testing.T) {
	command := []byte("\x1b[38;2;32;64;128m")
	cmd, size := Decode(command)
	sgr, ok := cmd.(*SetGraphicsRendition)
	assert.True(t, ok)
	assert.Equal(t, 17, size)
	assert.Equal(t, SGRForegroundColor, sgr.Command)
	assert.Equal(t, []int{2, 32, 64, 128}, sgr.Parameters)

	var b bytes.Buffer
	encodedSize, err := sgr.Encode(&b)
	assert.NoError(t, err)
	assert.Equal(t, size, encodedSize)
	assert.Equal(t, command, b.Bytes())
}

// default foreground color (implementation-defined)
func TestSGR_ForegroundDefault(t *testing.T) {
	command := []byte("\x1b[39m")
	cmd, size := Decode(command)
	sgr, ok := cmd.(*SetGraphicsRendition)
	assert.True(t, ok)
	assert.Equal(t, 5, size)
	assert.Equal(t, SGRForegroundDefault, sgr.Command)

	var b bytes.Buffer
	encodedSize, err := sgr.Encode(&b)
	assert.NoError(t, err)
	assert.Equal(t, size, encodedSize)
	assert.Equal(t, command, b.Bytes())
}

// black background color
func TestSGR_BackgroundBlack(t *testing.T) {
	command := []byte("\x1b[40m")
	cmd, size := Decode(command)
	sgr, ok := cmd.(*SetGraphicsRendition)
	assert.True(t, ok)
	assert.Equal(t, 5, size)
	assert.Equal(t, SGRBackgroundBlack, sgr.Command)

	var b bytes.Buffer
	encodedSize, err := sgr.Encode(&b)
	assert.NoError(t, err)
	assert.Equal(t, size, encodedSize)
	assert.Equal(t, command, b.Bytes())
}

// red background color
func TestSGR_BackgroundRed(t *testing.T) {
	command := []byte("\x1b[41m")
	cmd, size := Decode(command)
	sgr, ok := cmd.(*SetGraphicsRendition)
	assert.True(t, ok)
	assert.Equal(t, 5, size)
	assert.Equal(t, SGRBackgroundRed, sgr.Command)

	var b bytes.Buffer
	encodedSize, err := sgr.Encode(&b)
	assert.NoError(t, err)
	assert.Equal(t, size, encodedSize)
	assert.Equal(t, command, b.Bytes())
}

// green background color
func TestSGR_BackgroundGreen(t *testing.T) {
	command := []byte("\x1b[42m")
	cmd, size := Decode(command)
	sgr, ok := cmd.(*SetGraphicsRendition)
	assert.True(t, ok)
	assert.Equal(t, 5, size)
	assert.Equal(t, SGRBackgroundGreen, sgr.Command)

	var b bytes.Buffer
	encodedSize, err := sgr.Encode(&b)
	assert.NoError(t, err)
	assert.Equal(t, size, encodedSize)
	assert.Equal(t, command, b.Bytes())
}

// yellow background color
func TestSGR_BackgroundYellow(t *testing.T) {
	command := []byte("\x1b[43m")
	cmd, size := Decode(command)
	sgr, ok := cmd.(*SetGraphicsRendition)
	assert.True(t, ok)
	assert.Equal(t, 5, size)
	assert.Equal(t, SGRBackgroundYellow, sgr.Command)

	var b bytes.Buffer
	encodedSize, err := sgr.Encode(&b)
	assert.NoError(t, err)
	assert.Equal(t, size, encodedSize)
	assert.Equal(t, command, b.Bytes())
}

// blue background color
func TestSGR_BackgroundBlue(t *testing.T) {
	command := []byte("\x1b[44m")
	cmd, size := Decode(command)
	sgr, ok := cmd.(*SetGraphicsRendition)
	assert.True(t, ok)
	assert.Equal(t, 5, size)
	assert.Equal(t, SGRBackgroundBlue, sgr.Command)

	var b bytes.Buffer
	encodedSize, err := sgr.Encode(&b)
	assert.NoError(t, err)
	assert.Equal(t, size, encodedSize)
	assert.Equal(t, command, b.Bytes())
}

// magenta background color
func TestSGR_BackgroundMagenta(t *testing.T) {
	command := []byte("\x1b[45m")
	cmd, size := Decode(command)
	sgr, ok := cmd.(*SetGraphicsRendition)
	assert.True(t, ok)
	assert.Equal(t, 5, size)
	assert.Equal(t, SGRBackgroundMagenta, sgr.Command)

	var b bytes.Buffer
	encodedSize, err := sgr.Encode(&b)
	assert.NoError(t, err)
	assert.Equal(t, size, encodedSize)
	assert.Equal(t, command, b.Bytes())
}

// cyan background color
func TestSGR_BackgroundCyan(t *testing.T) {
	command := []byte("\x1b[46m")
	cmd, size := Decode(command)
	sgr, ok := cmd.(*SetGraphicsRendition)
	assert.True(t, ok)
	assert.Equal(t, 5, size)
	assert.Equal(t, SGRBackgroundCyan, sgr.Command)

	var b bytes.Buffer
	encodedSize, err := sgr.Encode(&b)
	assert.NoError(t, err)
	assert.Equal(t, size, encodedSize)
	assert.Equal(t, command, b.Bytes())
}

// white background color
func TestSGR_BackgroundWhite(t *testing.T) {
	command := []byte("\x1b[47m")
	cmd, size := Decode(command)
	sgr, ok := cmd.(*SetGraphicsRendition)
	assert.True(t, ok)
	assert.Equal(t, 5, size)
	assert.Equal(t, SGRBackgroundWhite, sgr.Command)

	var b bytes.Buffer
	encodedSize, err := sgr.Encode(&b)
	assert.NoError(t, err)
	assert.Equal(t, size, encodedSize)
	assert.Equal(t, command, b.Bytes())
}

// set the background color to a 256-color or true color value
func TestSGR_BackgroundColor_256(t *testing.T) {
	command := []byte("\x1b[48;5;128m")
	cmd, size := Decode(command)
	sgr, ok := cmd.(*SetGraphicsRendition)
	assert.True(t, ok)
	assert.Equal(t, 11, size)
	assert.Equal(t, SGRBackgroundColor, sgr.Command)
	assert.Equal(t, []int{5, 128}, sgr.Parameters)

	var b bytes.Buffer
	encodedSize, err := sgr.Encode(&b)
	assert.NoError(t, err)
	assert.Equal(t, size, encodedSize)
	assert.Equal(t, command, b.Bytes())
}

// set the background color to a 256-color or true color value
func TestSGR_BackgroundColor_Truecolor(t *testing.T) {
	command := []byte("\x1b[48;2;32;64;128m")
	cmd, size := Decode(command)
	sgr, ok := cmd.(*SetGraphicsRendition)
	assert.True(t, ok)
	assert.Equal(t, 17, size)
	assert.Equal(t, SGRBackgroundColor, sgr.Command)
	assert.Equal(t, []int{2, 32, 64, 128}, sgr.Parameters)

	var b bytes.Buffer
	encodedSize, err := sgr.Encode(&b)
	assert.NoError(t, err)
	assert.Equal(t, size, encodedSize)
	assert.Equal(t, command, b.Bytes())
}

// default background color (implementation-defined)
func TestSGR_BackgroundDefault(t *testing.T) {
	command := []byte("\x1b[49m")
	cmd, size := Decode(command)
	sgr, ok := cmd.(*SetGraphicsRendition)
	assert.True(t, ok)
	assert.Equal(t, 5, size)
	assert.Equal(t, SGRBackgroundDefault, sgr.Command)

	var b bytes.Buffer
	encodedSize, err := sgr.Encode(&b)
	assert.NoError(t, err)
	assert.Equal(t, size, encodedSize)
	assert.Equal(t, command, b.Bytes())
}

// doisable proportional spacing
func TestSGR_NoProportionalSpacing(t *testing.T) {
	command := []byte("\x1b[50m")
	cmd, size := Decode(command)
	sgr, ok := cmd.(*SetGraphicsRendition)
	assert.True(t, ok)
	assert.Equal(t, 5, size)
	assert.Equal(t, SGRNoProportionalSpacing, sgr.Command)

	var b bytes.Buffer
	encodedSize, err := sgr.Encode(&b)
	assert.NoError(t, err)
	assert.Equal(t, size, encodedSize)
	assert.Equal(t, command, b.Bytes())
}

// framed
func TestSGR_Frame(t *testing.T) {
	command := []byte("\x1b[51m")
	cmd, size := Decode(command)
	sgr, ok := cmd.(*SetGraphicsRendition)
	assert.True(t, ok)
	assert.Equal(t, 5, size)
	assert.Equal(t, SGRFrame, sgr.Command)

	var b bytes.Buffer
	encodedSize, err := sgr.Encode(&b)
	assert.NoError(t, err)
	assert.Equal(t, size, encodedSize)
	assert.Equal(t, command, b.Bytes())
}

// encircled
func TestSGR_Encircle(t *testing.T) {
	command := []byte("\x1b[52m")
	cmd, size := Decode(command)
	sgr, ok := cmd.(*SetGraphicsRendition)
	assert.True(t, ok)
	assert.Equal(t, 5, size)
	assert.Equal(t, SGREncircle, sgr.Command)

	var b bytes.Buffer
	encodedSize, err := sgr.Encode(&b)
	assert.NoError(t, err)
	assert.Equal(t, size, encodedSize)
	assert.Equal(t, command, b.Bytes())
}

// overlined
func TestSGR_Overline(t *testing.T) {
	command := []byte("\x1b[53m")
	cmd, size := Decode(command)
	sgr, ok := cmd.(*SetGraphicsRendition)
	assert.True(t, ok)
	assert.Equal(t, 5, size)
	assert.Equal(t, SGROverline, sgr.Command)

	var b bytes.Buffer
	encodedSize, err := sgr.Encode(&b)
	assert.NoError(t, err)
	assert.Equal(t, size, encodedSize)
	assert.Equal(t, command, b.Bytes())
}

// not framed or encircled
func TestSGR_NoFrameOrEncircle(t *testing.T) {
	command := []byte("\x1b[54m")
	cmd, size := Decode(command)
	sgr, ok := cmd.(*SetGraphicsRendition)
	assert.True(t, ok)
	assert.Equal(t, 5, size)
	assert.Equal(t, SGRNoFrameOrEncircle, sgr.Command)

	var b bytes.Buffer
	encodedSize, err := sgr.Encode(&b)
	assert.NoError(t, err)
	assert.Equal(t, size, encodedSize)
	assert.Equal(t, command, b.Bytes())
}

// not overlined
func TestSGR_NoOverline(t *testing.T) {
	command := []byte("\x1b[55m")
	cmd, size := Decode(command)
	sgr, ok := cmd.(*SetGraphicsRendition)
	assert.True(t, ok)
	assert.Equal(t, 5, size)
	assert.Equal(t, SGRNoOverline, sgr.Command)

	var b bytes.Buffer
	encodedSize, err := sgr.Encode(&b)
	assert.NoError(t, err)
	assert.Equal(t, size, encodedSize)
	assert.Equal(t, command, b.Bytes())
}

// set the underline color to a 256-color or true color value
func TestSGR_UnderlineColor_256(t *testing.T) {
	command := []byte("\x1b[58;5;128m")
	cmd, size := Decode(command)
	sgr, ok := cmd.(*SetGraphicsRendition)
	assert.True(t, ok)
	assert.Equal(t, 11, size)
	assert.Equal(t, SGRUnderlineColor, sgr.Command)
	assert.Equal(t, []int{5, 128}, sgr.Parameters)

	var b bytes.Buffer
	encodedSize, err := sgr.Encode(&b)
	assert.NoError(t, err)
	assert.Equal(t, size, encodedSize)
	assert.Equal(t, command, b.Bytes())
}

// set the underline color to a 256-color or true color value
func TestSGR_UnderlineColor_Truecolor(t *testing.T) {
	command := []byte("\x1b[58;2;32;64;128m")
	cmd, size := Decode(command)
	sgr, ok := cmd.(*SetGraphicsRendition)
	assert.True(t, ok)
	assert.Equal(t, 17, size)
	assert.Equal(t, SGRUnderlineColor, sgr.Command)
	assert.Equal(t, []int{2, 32, 64, 128}, sgr.Parameters)

	var b bytes.Buffer
	encodedSize, err := sgr.Encode(&b)
	assert.NoError(t, err)
	assert.Equal(t, size, encodedSize)
	assert.Equal(t, command, b.Bytes())
}

// default underline color
func TestSGR_DefaultUnderlineColor(t *testing.T) {
	command := []byte("\x1b[59m")
	cmd, size := Decode(command)
	sgr, ok := cmd.(*SetGraphicsRendition)
	assert.True(t, ok)
	assert.Equal(t, 5, size)
	assert.Equal(t, SGRDefaultUnderlineColor, sgr.Command)

	var b bytes.Buffer
	encodedSize, err := sgr.Encode(&b)
	assert.NoError(t, err)
	assert.Equal(t, size, encodedSize)
	assert.Equal(t, command, b.Bytes())
}

// ideogram underline or right side line
func TestSGR_IdeogramUnderline(t *testing.T) {
	command := []byte("\x1b[60m")
	cmd, size := Decode(command)
	sgr, ok := cmd.(*SetGraphicsRendition)
	assert.True(t, ok)
	assert.Equal(t, 5, size)
	assert.Equal(t, SGRIdeogramUnderline, sgr.Command)

	var b bytes.Buffer
	encodedSize, err := sgr.Encode(&b)
	assert.NoError(t, err)
	assert.Equal(t, size, encodedSize)
	assert.Equal(t, command, b.Bytes())
}

// ideogram double underline or double line on the right side
func TestSGR_IdeogramDoubleUnderline(t *testing.T) {
	command := []byte("\x1b[61m")
	cmd, size := Decode(command)
	sgr, ok := cmd.(*SetGraphicsRendition)
	assert.True(t, ok)
	assert.Equal(t, 5, size)
	assert.Equal(t, SGRIdeogramDoubleUnderline, sgr.Command)

	var b bytes.Buffer
	encodedSize, err := sgr.Encode(&b)
	assert.NoError(t, err)
	assert.Equal(t, size, encodedSize)
	assert.Equal(t, command, b.Bytes())
}

// ideogram overline or left side line
func TestSGR_IdeogramOverline(t *testing.T) {
	command := []byte("\x1b[62m")
	cmd, size := Decode(command)
	sgr, ok := cmd.(*SetGraphicsRendition)
	assert.True(t, ok)
	assert.Equal(t, 5, size)
	assert.Equal(t, SGRIdeogramOverline, sgr.Command)

	var b bytes.Buffer
	encodedSize, err := sgr.Encode(&b)
	assert.NoError(t, err)
	assert.Equal(t, size, encodedSize)
	assert.Equal(t, command, b.Bytes())
}

// ideogram double overline or double line on the left side
func TestSGR_IdeogramDoubleOverline(t *testing.T) {
	command := []byte("\x1b[63m")
	cmd, size := Decode(command)
	sgr, ok := cmd.(*SetGraphicsRendition)
	assert.True(t, ok)
	assert.Equal(t, 5, size)
	assert.Equal(t, SGRIdeogramDoubleOverline, sgr.Command)

	var b bytes.Buffer
	encodedSize, err := sgr.Encode(&b)
	assert.NoError(t, err)
	assert.Equal(t, size, encodedSize)
	assert.Equal(t, command, b.Bytes())
}

// ideogram stress marking
func TestSGR_IdeogramStress(t *testing.T) {
	command := []byte("\x1b[64m")
	cmd, size := Decode(command)
	sgr, ok := cmd.(*SetGraphicsRendition)
	assert.True(t, ok)
	assert.Equal(t, 5, size)
	assert.Equal(t, SGRIdeogramStress, sgr.Command)

	var b bytes.Buffer
	encodedSize, err := sgr.Encode(&b)
	assert.NoError(t, err)
	assert.Equal(t, size, encodedSize)
	assert.Equal(t, command, b.Bytes())
}

// cancels the effect of the rendition aspects established by ideogram parameter values
func TestSGR_IdeogramReset(t *testing.T) {
	command := []byte("\x1b[65m")
	cmd, size := Decode(command)
	sgr, ok := cmd.(*SetGraphicsRendition)
	assert.True(t, ok)
	assert.Equal(t, 5, size)
	assert.Equal(t, SGRIdeogramReset, sgr.Command)

	var b bytes.Buffer
	encodedSize, err := sgr.Encode(&b)
	assert.NoError(t, err)
	assert.Equal(t, size, encodedSize)
	assert.Equal(t, command, b.Bytes())
}
