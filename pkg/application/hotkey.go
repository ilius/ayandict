package application

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"golang.design/x/hotkey"
)

// Convert a hotkey string like "Ctrl+Shift+R" → ([]Modifier, Key)
func ParseHotkeyString(seq string) ([]hotkey.Modifier, hotkey.Key, error) {
	seq = strings.ToLower(seq)
	parts := strings.Split(seq, "+")
	var mods []hotkey.Modifier
	var key hotkey.Key
	for _, p := range parts {
		p = strings.TrimSpace(p)
		switch p {
		case "ctrl", "control":
			mods = append(mods, hotkey.ModCtrl)
		case "alt", "mod1", "option":
			mods = append(mods, hotkey.Mod1)
		case "shift":
			mods = append(mods, hotkey.ModShift)
		case "cmd", "meta", "win", "super", "mod4":
			mods = append(mods, hotkey.Mod4)
		default:
			if key != 0 {
				return nil, 0, fmt.Errorf("multiple keys")
			}
			key = parseKey(p)
		}
	}
	return mods, key, nil
}

func parseKey(p string) hotkey.Key {
	// Handle letters A-Z
	if len(p) == 1 && unicode.IsLetter(rune(p[0])) {
		r := unicode.ToUpper(rune(p[0]))
		return hotkey.Key(r)
	}
	// Handle digits 0-9
	if len(p) == 1 && unicode.IsDigit(rune(p[0])) {
		return hotkey.Key(p[0])
	}
	// Handle function keys (F1–F24)
	if strings.HasPrefix(p, "f") {
		if n, err := strconv.Atoi(p[1:]); err == nil && n >= 1 && n <= 24 {
			return hotkey.Key(uint16(hotkey.KeyF1) + uint16(n) - 1)
		}
	}
	// Could add arrows, escape, etc.
	return 0
}
