package ansicsi

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCSI(t *testing.T) {
	input := []byte("\x1b[38;5;128m\x1b[1mHello, world!\x1b[0m")

	expected := []struct {
		size int
		cmd  Command
	}{
		{size: 11, cmd: &SetGraphicsRendition{Command: SGRForegroundColor, Parameters: []int{5, 128}}},
		{size: 4, cmd: &SetGraphicsRendition{Command: SGRBold, Parameters: []int{}}},
		{size: 0},
		{size: 0},
		{size: 0},
		{size: 0},
		{size: 0},
		{size: 0},
		{size: 0},
		{size: 0},
		{size: 0},
		{size: 0},
		{size: 0},
		{size: 0},
		{size: 0},
		{size: 4, cmd: &SetGraphicsRendition{Command: SGRReset, Parameters: []int{}}},
	}

	var buf bytes.Buffer
	for b := input; len(b) > 0; expected = expected[1:] {
		cmd, size := Decode(b)
		if !assert.Equal(t, expected[0].size, size) {
			return
		}
		if size == 0 {
			buf.WriteByte(b[0])
			b = b[1:]
			continue
		}
		b = b[size:]

		if !assert.Equal(t, expected[0].cmd, cmd) {
			return
		}

		encodedSize, err := cmd.Encode(&buf)
		if !assert.NoError(t, err) || !assert.Equal(t, size, encodedSize) {
			return
		}
	}
	assert.Equal(t, input, buf.Bytes())
}
