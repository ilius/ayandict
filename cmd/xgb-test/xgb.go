package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jezek/xgb"
	"github.com/jezek/xgb/xproto"
)

// minimal map from string names to keysyms (ASCII uppercase letters)
var keysymMap = map[string]xproto.Keysym{
	"A": 0x0041, "B": 0x0042, "C": 0x0043, "D": 0x0044, "E": 0x0045,
	"F": 0x0046, "G": 0x0047, "H": 0x0048, "I": 0x0049, "J": 0x004A,
	"K": 0x004B, "L": 0x004C, "M": 0x004D, "N": 0x004E, "O": 0x004F,
	"P": 0x0050, "Q": 0x0051, "R": 0x0052, "S": 0x0053, "T": 0x0054,
	"U": 0x0055, "V": 0x0056, "W": 0x0057, "X": 0x0058, "Y": 0x0059, "Z": 0x005A,
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: hotkey-grabber <Hotkey1> [Hotkey2 ...]")
		fmt.Println("Example: hotkey-grabber Ctrl+Alt+Z")
		return
	}

	X, err := xgb.NewConn()
	if err != nil {
		log.Fatalf("Cannot connect to X server: %v", err)
	}
	defer X.Close()

	setup := xproto.Setup(X)
	root := setup.DefaultScreen(X).Root

	// Enable KeyPress events on the root window
	err = xproto.ChangeWindowAttributesChecked(X, root,
		xproto.CwEventMask,
		[]uint32{xproto.EventMaskKeyPress}).Check()
	if err != nil {
		log.Fatalf("Failed to select input events: %v", err)
	}

	for _, hotkey := range os.Args[1:] {
		modifiers, keycode := parseHotkey(X, hotkey)
		if keycode == 0 {
			fmt.Printf("Invalid hotkey: %s\n", hotkey)
			continue
		}
		err := xproto.GrabKeyChecked(X, true, root,
			uint16(modifiers),
			keycode,
			xproto.GrabModeAsync,
			xproto.GrabModeAsync).Check()
		if err != nil {
			log.Fatalf("Failed to grab hotkey %s: %v", hotkey, err)
		}
		fmt.Printf("Grabbed hotkey: %s (mod=%d, keycode=%d)\n", hotkey, modifiers, keycode)
	}

	fmt.Println("Listening for hotkeys...")
	for {
		ev, err := X.WaitForEvent()
		if err != nil {
			log.Fatal(err)
		}
		switch e := ev.(type) {
		case xproto.KeyPressEvent:
			fmt.Printf("Hotkey pressed! keycode=%d state=%d\n", e.Detail, e.State)
		}
	}
}

func parseHotkey(X *xgb.Conn, combo string) (modifiers int, keycode xproto.Keycode) {
	parts := strings.Split(strings.ToUpper(combo), "+")
	var key string

	for _, p := range parts {
		switch p {
		case "CTRL", "CONTROL":
			modifiers |= int(xproto.ModMaskControl)
		case "ALT":
			modifiers |= int(xproto.ModMask1)
		case "SHIFT":
			modifiers |= int(xproto.ModMaskShift)
		case "SUPER", "META", "WIN":
			modifiers |= int(xproto.ModMask4)
		default:
			key = p
		}
	}
	keysym, ok := keysymMap[key]
	if !ok {
		return modifiers, 0
	}

	// Translate keysym to keycode
	minKeycode := xproto.Setup(X).MinKeycode
	maxKeycode := xproto.Setup(X).MaxKeycode
	reply, err := xproto.GetKeyboardMapping(X, minKeycode, byte(maxKeycode-minKeycode+1)).Reply()
	if err != nil {
		return modifiers, 0
	}

	keysymsPerKeycode := int(reply.KeysymsPerKeycode)
	for i := 0; i < len(reply.Keysyms); i += keysymsPerKeycode {
		if reply.Keysyms[i] == keysym {
			return modifiers, minKeycode + xproto.Keycode(i/keysymsPerKeycode)
		}
	}
	return modifiers, 0
}
