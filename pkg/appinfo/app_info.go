package appinfo

import (
	"os/user"
)

const (
	APP_NAME = "ayandict"
	APP_DESC = "AyanDict"
	VERSION  = "v3.0.0beta5"
)

var LOCAL_SOCKET_NAME = APP_NAME + "-" + userId()

const ABOUT = `A simple cross-platform desktop dictionary application based on Qt framework and written in Go that uses StarDict dictionary format.

Copyleft © 2025 Saeed Rasooli
AyanDict is licensed by the GNU Affero General Public License version 3 (or later)
`

const AUTHORS = `Saeed Rasooli <saeed.gnu@gmail.com> (ilius)`

const LICENSE = `AyanDict - A simple dictionary application

Copyright © 2025 Saeed Rasooli

This program is free software: you can redistribute it and/or modify it under the terms of the GNU Affero General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.
`

func userId() string {
	u, err := user.Current()
	if err != nil {
		panic(err)
	}
	return u.Uid
}

var KeyBindings = [][3]string{
	{"Space", "query is not focused", "Change keyboard focus to query entry"},
	{"Escape", "query is focused", "Focus leaves the query entry"},
	{"Escape", "query is not focused", "Clear the query and results"},
	{"+ or =", "", "Zoom in"},
	{"–", "", "Zoom out"},
	{"Alt + Left", "", "Go back in history"},
	{"Alt + Right", "", "Go forward in history"},
	{"Ctrl + Q", "", "Quit / exit application"},
	{"F1", "", "Show About window"},
}
