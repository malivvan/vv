package cui

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"unicode"

	"github.com/gdamore/tcell/v2"
)

// Modifier labels
const (
	LabelCtrl  = "ctrl"
	LabelAlt   = "alt"
	LabelMeta  = "meta"
	LabelShift = "shift"
)

// ErrInvalidKeyEvent is the error returned when encoding or decoding a key event fails.
var ErrInvalidKeyEvent = errors.New("invalid key event")

// UnifyEnterKeys is a flag that determines whether or not KPEnter (keypad
// enter) key events are interpreted as Enter key events. When enabled, Ctrl+J
// key events are also interpreted as Enter key events.
var UnifyEnterKeys = true

var fullKeyNames = map[string]string{
	"backspace2": "Backspace",
	"pgup":       "PageUp",
	"pgdn":       "PageDown",
	"esc":        "Escape",
}

var ctrlKeys = map[rune]tcell.Key{
	' ':  tcell.KeyCtrlSpace,
	'a':  tcell.KeyCtrlA,
	'b':  tcell.KeyCtrlB,
	'c':  tcell.KeyCtrlC,
	'd':  tcell.KeyCtrlD,
	'e':  tcell.KeyCtrlE,
	'f':  tcell.KeyCtrlF,
	'g':  tcell.KeyCtrlG,
	'h':  tcell.KeyCtrlH,
	'i':  tcell.KeyCtrlI,
	'j':  tcell.KeyCtrlJ,
	'k':  tcell.KeyCtrlK,
	'l':  tcell.KeyCtrlL,
	'm':  tcell.KeyCtrlM,
	'n':  tcell.KeyCtrlN,
	'o':  tcell.KeyCtrlO,
	'p':  tcell.KeyCtrlP,
	'q':  tcell.KeyCtrlQ,
	'r':  tcell.KeyCtrlR,
	's':  tcell.KeyCtrlS,
	't':  tcell.KeyCtrlT,
	'u':  tcell.KeyCtrlU,
	'v':  tcell.KeyCtrlV,
	'w':  tcell.KeyCtrlW,
	'x':  tcell.KeyCtrlX,
	'y':  tcell.KeyCtrlY,
	'z':  tcell.KeyCtrlZ,
	'\\': tcell.KeyCtrlBackslash,
	']':  tcell.KeyCtrlRightSq,
	'^':  tcell.KeyCtrlCarat,
	'_':  tcell.KeyCtrlUnderscore,
}

// BindDecode decodes a string as a key or combination of keys.
func BindDecode(s string) (mod tcell.ModMask, key tcell.Key, ch rune, err error) {
	if len(s) == 0 {
		return 0, 0, 0, ErrInvalidKeyEvent
	}

	// Special case for plus rune decoding
	if s[len(s)-1:] == "+" {
		key = tcell.KeyRune
		ch = '+'

		if len(s) == 1 {
			return mod, key, ch, nil
		} else if len(s) == 2 {
			return 0, 0, 0, ErrInvalidKeyEvent
		} else {
			s = s[:len(s)-2]
		}
	}

	split := strings.Split(s, "+")
DECODEPIECE:
	for _, piece := range split {
		// BindDecode modifiers
		pieceLower := strings.ToLower(piece)
		switch pieceLower {
		case LabelCtrl:
			mod |= tcell.ModCtrl
			continue
		case LabelAlt:
			mod |= tcell.ModAlt
			continue
		case LabelMeta:
			mod |= tcell.ModMeta
			continue
		case LabelShift:
			mod |= tcell.ModShift
			continue
		}

		// BindDecode key
		for shortKey, fullKey := range fullKeyNames {
			if pieceLower == strings.ToLower(fullKey) {
				pieceLower = shortKey
				break
			}
		}
		switch pieceLower {
		case "backspace":
			key = tcell.KeyBackspace2
			continue
		case "space", "spacebar":
			key = tcell.KeyRune
			ch = ' '
			continue
		}
		for k, keyName := range tcell.KeyNames {
			if pieceLower == strings.ToLower(strings.ReplaceAll(keyName, "-", "+")) {
				key = k
				if key < 0x80 {
					ch = rune(k)
				}
				continue DECODEPIECE
			}
		}

		// BindDecode rune
		if len(piece) > 1 {
			return 0, 0, 0, ErrInvalidKeyEvent
		}

		key = tcell.KeyRune
		ch = rune(piece[0])
	}

	if mod&tcell.ModCtrl != 0 {
		k, ok := ctrlKeys[unicode.ToLower(ch)]
		if ok {
			key = k
			if UnifyEnterKeys && key == ctrlKeys['j'] {
				key = tcell.KeyEnter
			} else if key < 0x80 {
				ch = rune(key)
			}
		}
	}

	return mod, key, ch, nil
}

// BindEncode encodes a key or combination of keys a string.
func BindEncode(mod tcell.ModMask, key tcell.Key, ch rune) (string, error) {
	var b strings.Builder
	var wrote bool

	if mod&tcell.ModCtrl != 0 {
		if key == tcell.KeyBackspace || key == tcell.KeyTab || key == tcell.KeyEnter {
			mod ^= tcell.ModCtrl
		} else {
			for _, ctrlKey := range ctrlKeys {
				if key == ctrlKey {
					mod ^= tcell.ModCtrl
					break
				}
			}
		}
	}

	if key != tcell.KeyRune {
		if UnifyEnterKeys && key == ctrlKeys['j'] {
			key = tcell.KeyEnter
		} else if key < 0x80 {
			ch = rune(key)
		}
	}

	// BindEncode modifiers
	if mod&tcell.ModCtrl != 0 {
		b.WriteString(upperFirst(LabelCtrl))
		wrote = true
	}
	if mod&tcell.ModAlt != 0 {
		if wrote {
			b.WriteRune('+')
		}
		b.WriteString(upperFirst(LabelAlt))
		wrote = true
	}
	if mod&tcell.ModMeta != 0 {
		if wrote {
			b.WriteRune('+')
		}
		b.WriteString(upperFirst(LabelMeta))
		wrote = true
	}
	if mod&tcell.ModShift != 0 {
		if wrote {
			b.WriteRune('+')
		}
		b.WriteString(upperFirst(LabelShift))
		wrote = true
	}

	if key == tcell.KeyRune && ch == ' ' {
		if wrote {
			b.WriteRune('+')
		}
		b.WriteString("Space")
	} else if key != tcell.KeyRune {
		// BindEncode key
		keyName := tcell.KeyNames[key]
		if keyName == "" {
			return "", ErrInvalidKeyEvent
		}
		keyName = strings.ReplaceAll(keyName, "-", "+")
		fullKeyName := fullKeyNames[strings.ToLower(keyName)]
		if fullKeyName != "" {
			keyName = fullKeyName
		}

		if wrote {
			b.WriteRune('+')
		}
		b.WriteString(keyName)
	} else {
		// BindEncode rune
		if wrote {
			b.WriteRune('+')
		}
		b.WriteRune(ch)
	}

	return b.String(), nil
}

func upperFirst(s string) string {
	if len(s) <= 1 {
		return strings.ToUpper(s)
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

type eventHandler func(ev *tcell.EventKey) *tcell.EventKey

// BindConfig maps keys to event handlers and processes key events.
type BindConfig struct {
	handlers map[string]eventHandler
	mutex    *sync.RWMutex
}

// NewBindConfig returns a new input configuration.
func NewBindConfig() *BindConfig {
	c := BindConfig{
		handlers: make(map[string]eventHandler),
		mutex:    new(sync.RWMutex),
	}

	return &c
}

// Set sets the handler for a key event string.
func (c *BindConfig) Set(s string, handler func(ev *tcell.EventKey) *tcell.EventKey) error {
	mod, key, ch, err := BindDecode(s)
	if err != nil {
		return err
	}

	if key == tcell.KeyRune {
		c.SetRune(mod, ch, handler)
	} else {
		c.SetKey(mod, key, handler)
	}
	return nil
}

// SetKey sets the handler for a key.
func (c *BindConfig) SetKey(mod tcell.ModMask, key tcell.Key, handler func(ev *tcell.EventKey) *tcell.EventKey) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if mod&tcell.ModShift != 0 && key == tcell.KeyTab {
		mod ^= tcell.ModShift
		key = tcell.KeyBacktab
	}

	if mod&tcell.ModCtrl == 0 && key != tcell.KeyBackspace && key != tcell.KeyTab && key != tcell.KeyEnter {
		for _, ctrlKey := range ctrlKeys {
			if key == ctrlKey {
				mod |= tcell.ModCtrl
				break
			}
		}
	}

	c.handlers[fmt.Sprintf("%d-%d", mod, key)] = handler
}

// SetRune sets the handler for a rune.
func (c *BindConfig) SetRune(mod tcell.ModMask, ch rune, handler func(ev *tcell.EventKey) *tcell.EventKey) {
	// Some runes are identical to named keys. Set the bind on the matching
	// named key instead.
	switch ch {
	case '\t':
		c.SetKey(mod, tcell.KeyTab, handler)
		return
	case '\n':
		c.SetKey(mod, tcell.KeyEnter, handler)
		return
	}

	if mod&tcell.ModCtrl != 0 {
		k, ok := ctrlKeys[unicode.ToLower(ch)]
		if ok {
			c.SetKey(mod, k, handler)
			return
		}
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.handlers[fmt.Sprintf("%d:%d", mod, ch)] = handler
}

// Capture handles key events.
func (c *BindConfig) Capture(ev *tcell.EventKey) *tcell.EventKey {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if ev == nil {
		return nil
	}

	var keyName string
	if ev.Key() != tcell.KeyRune {
		keyName = fmt.Sprintf("%d-%d", ev.Modifiers(), ev.Key())
	} else {
		keyName = fmt.Sprintf("%d:%d", ev.Modifiers(), ev.Rune())
	}

	handler := c.handlers[keyName]
	if handler != nil {
		return handler(ev)
	}
	return ev
}

// Clear removes all handlers.
func (c *BindConfig) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.handlers = make(map[string]eventHandler)
}
