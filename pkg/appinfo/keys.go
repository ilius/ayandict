package appinfo

var KeyBindings1 = [][3]string{
	{"＋  or  ＝", "", "Zoom in"},
	{"－", "", "Zoom out"},
	{"Ctrl﹢F", "", "Search in article text"},
	{"Ctrl﹢G", "", "Goto next match"},

	{"Alt﹢Down", "", "Goto next result"},
	{"Alt﹢Up", "", "Goto previous result"},

	{"F", "", "Add/Remove Favorite"},
	{"Ctrl﹢Ｄ", "", "Manage Dictionaries"},
	{"Ctrl﹢Ｒ", "", "Reload Config"},
	{"Ctrl﹢Shift﹢R", "", "Reload Config, Dicts, Style"},
	{"Ctrl﹢Delete", "", "Clear History"},
}

var KeyBindings2 = [][3]string{
	{"Esc", "search bar is visible", "Hide search bar"},
	{"Esc", "query is focused", "Unfocus query entry"},
	{"Esc", "none of two above", "Clear query and results"},

	{" Space ", "  query is not focused  ", "Focus query entry"},

	{"Alt﹢Left ｜ Ctrl﹢Left", "", "Go back in history"},
	{"Alt﹢Right ｜ Ctrl﹢Right", "", "Go forward in history"},

	{"F1", "", "Show About window"},
	{"Ctrl﹢Ｑ", "", "Quit / exit application"},

	{"Enter", "Scan popup", "Query in main window"},
}
