package hotkey

import (
	"fmt"
	"strings"

	"golang.design/x/hotkey"
)

// ParseHotkey parses a hotkey string (e.g. "Ctrl+Alt+Space") into modifiers and key
func ParseHotkey(s string) ([]hotkey.Modifier, hotkey.Key, error) {
	parts := strings.Split(s, "+")
	if len(parts) == 0 {
		return nil, 0, fmt.Errorf("empty hotkey string")
	}

	var mods []hotkey.Modifier
	var key hotkey.Key
	var err error

	for i, part := range parts {
		part = strings.TrimSpace(part)
		// Last part is usually the key, unless it's a modifier itself (unlikely for global hotkey)
		if i == len(parts)-1 {
			key, err = parseKey(part)
			if err != nil {
				return nil, 0, err
			}
		} else {
			mod, err := parseModifier(part)
			if err != nil {
				return nil, 0, err
			}
			mods = append(mods, mod)
		}
	}

	return mods, key, nil
}

func parseModifier(s string) (hotkey.Modifier, error) {
	switch strings.ToLower(s) {
	case "ctrl", "control":
		return hotkey.ModCtrl, nil
	case "alt":
		return hotkey.ModAlt, nil
	case "shift":
		return hotkey.ModShift, nil
	case "win", "cmd", "command", "super":
		return hotkey.ModWin, nil
	default:
		return 0, fmt.Errorf("unknown modifier: %s", s)
	}
}

func parseKey(s string) (hotkey.Key, error) {
	s = strings.ToUpper(s)
	switch s {
	case "SPACE":
		return hotkey.KeySpace, nil
	case "ENTER", "RETURN":
		return hotkey.KeyReturn, nil
	case "ESC", "ESCAPE":
		return hotkey.KeyEscape, nil
	case "TAB":
		return hotkey.KeyTab, nil
	case "UP":
		return hotkey.KeyUp, nil
	case "DOWN":
		return hotkey.KeyDown, nil
	case "LEFT":
		return hotkey.KeyLeft, nil
	case "RIGHT":
		return hotkey.KeyRight, nil
	// Add other special keys as needed
	default:
		if len(s) == 1 {
			r := rune(s[0])
			if r >= 'A' && r <= 'Z' {
				return hotkey.Key(r), nil
			}
			if r >= '0' && r <= '9' {
				return hotkey.Key(r), nil
			}
		}
		// F1-F24
		if strings.HasPrefix(s, "F") {
			// Simple F-key parsing could go here if needed
		}
		return 0, fmt.Errorf("unknown key: %s", s)
	}
}
