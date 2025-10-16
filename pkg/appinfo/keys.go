//go:build !darwin
// +build !darwin

package appinfo

var KeyBindings = [][3]string{
	{"Space", "query is not focused", "Change keyboard focus to query entry"},
	{"Escape", "query is focused", "Focus leaves the query entry"},
	{"Escape", "query is not focused", "Clear the query and results"},
	{"+ or =", "", "Zoom in"},
	{"–", "", "Zoom out"},
	{"Alt + Left or Ctrl + Left", "", "Go back in history"},
	{"Alt + Right or Ctrl + Right", "", "Go forward in history"},
	{"Ctrl + Q", "", "Quit / exit application"},
	{"F1", "", "Show About window"},
}
