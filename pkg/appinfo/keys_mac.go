//go:build darwin || mackeybindings
// +build darwin mackeybindings

package appinfo

func init() {
	KeyBindings = [][3]string{
		{" Space ", "  query is not focused  ", "Focus query entry"},
		{"Esc", "query is focused", "Unfocus query entry"},
		{"Esc", "query is not focused", "Clear query and results"},
		{"＋  or  ＝", "", "Zoom in"},
		{"－", "", "Zoom out"},
		{"⌘ ←    or    ⌥ ←", "", "Go back in history"},
		{"⌘ →    or    ⌥ →", "", "Go forward in history"},
		{"⌘ Ｑ", "", "Quit / exit application"},
		{"F1", "", "Show About window"},
	}
}

// ⌘      Command, Cmd
// ⌃      Control, Ctl, Ctrl
// ⌥      Option, Opt, (PC) Alt
// ⇧      Shift
// ⇪      Caps lock
// ↩️      Return, Carriage Return
// ↵      Return, Carriage Return
// ⏎      Return, Carriage Return
// ⌤      Enter
// ⌫      Delete, Backspace
// ⌦      Forward Delete
// ⎋      Escape, Esc
// →      Right arrow
// ←      Left arrow
// ↑      Up arrow
// ↓      Down arrow
// ⇞      Page Up, PgUp
// ⇟      Page Down, PgDn
// ↖️      Home
// ↘️      End
// ⌧      Clear
// ⇥      Tab, Tab Right, Horizontal Tab
// ⇤      Shift Tab, Tab Left, Back-tab
// ␢      Space, Blank
// ␣      Space, Blank
