//go:build darwin
// +build darwin

package appinfo

var KeyBindings = [][3]string{
	{"Space", "query is not focused", "Change keyboard focus to query entry"},
	{"Escape", "query is focused", "Focus leaves the query entry"},
	{"Escape", "query is not focused", "Clear the query and results"},
	{"+ or =", "", "Zoom in"},
	{"–", "", "Zoom out"},
	{"Command + Left or Option + Left", "", "Go back in history"},
	{"Command + Right or Option + Right", "", "Go forward in history"},
	{"Command + Q", "", "Quit / exit application"},
	{"F1", "", "Show About window"},
}
