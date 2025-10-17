package appinfo

var KeyBindings = [][3]string{
	{" Space ", "  query is not focused  ", "Focus query entry"},
	{"Esc", "query is focused", "Unfocus query entry"},
	{"Esc", "query is not focused", "Clear query and results"},
	{"＋  or  ＝", "", "Zoom in"},
	{"－", "", "Zoom out"},
	{"Alt﹢Left ｜ Ctrl﹢Left", "", "Go back in history"},
	{"Alt﹢Right ｜ Ctrl﹢Right", "", "Go forward in history"},
	{"Ctrl﹢Ｑ", "", "Quit / exit application"},
	{"F1", "", "Show About window"},
	{"Enter", "Scan popup", "Query in main window"},
}
