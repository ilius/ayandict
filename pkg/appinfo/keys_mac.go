//go:build darwin || mackeybindings
// +build darwin mackeybindings

package appinfo

func init() {
	KeyBindings = [][3]string{
		{"＋  or  ＝", "", "Zoom in"},
		{"－", "", "Zoom out"},
		{"⌘ F", "", "Search in article text"},
		{"⌘ G", "", "Goto next match (via ⌘ F)"},

		{"Esc", "search bar is visible", "Hide search bar"},
		{"Esc", "query is focused", "Unfocus query entry"},
		{"Esc", "none of two above", "Clear query and results"},

		{" Space ", "  query is not focused  ", "Focus query entry"},

		{"⌘ ←    or    ⌥ ←", "", "Go back in history"},
		{"⌘ →    or    ⌥ →", "", "Go forward in history"},

		{"⌥ ↓", "", "Goto next result"},
		{"⌥ ↑", "", "Goto previous result"},

		{"⌘ Ｑ", "", "Quit / exit application"},
		{"F1", "", "Show About window"},
		{"F", "", "Add/Remove Favorite"},

		{"Enter", "Scan popup", "Query in main window"},
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
