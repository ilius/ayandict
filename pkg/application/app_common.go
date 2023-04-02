package application

import (
	"log"

	"github.com/ilius/ayandict/pkg/frequency"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

const VERSION = "0.1.0"

const ABOUT = `AyanDict is simple dictionary appliation based on Qt and written in Go.

Copyleft © 2023 Saeed Rasooli
AyanDict is licensed by the GNU General Public License version 3 (or later)
`

const AUTHORS = `Saeed Rasooli <saeed.gnu@gmail.com> (ilius)`

const LICENSE = `AyanDict - A simple simple dictionary appliation

Copyright © 2023 Saeed Rasooli
This program is free software; you can redistribute it
and/or modify it under the terms of the GNU General Public
License as published by the Free Software Foundation; 
either version 3 of the License, or (at your option) any
later version.

This program is distributed in the hope that it will be
useful, but WITHOUT ANY WARRANTY; without even the implied
warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR
PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public
License along with this program. Or on Debian systems,
from /usr/share/common-licenses/GPL.
If not, see http://www.gnu.org/licenses/gpl.txt
`

var expanding = widgets.QSizePolicy__Expanding

// we trim these characters when user right-clicks on a word without selecting it
const punctuation = "!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~،؛؟۔"

// when double-click in QTextBrowser. some punctuations next to words
// are also selected, specially non-ascii ones,
// so we trim them on riht-click -> Query action or on middle-click action
const queryForceTrimChars = "‘’،؛"

var frequencyTable *frequency.FrequencyTable

func loadIcon(filename string) *gui.QIcon {
	data, err := res.ReadFile("res/" + filename)
	if err != nil {
		log.Println(err)
		return nil
	}
	pixmap := gui.NewQPixmap()
	// LoadFromData fails
	// https://github.com/therecipe/qt/issues/1193
	if !pixmap.LoadFromData(data, uint(len(data)), "PNG", 0) {
		log.Println("loadIcon: LoadFromData failed")
		return nil
	}
	// byteArray := core.NewQByteArray()
	// byteArray.SetRawData(data, uint(len(data)))
	// if !pixmap.LoadFromData2(byteArray, "", 0) {
	// 	log.Println("loadIcon: LoadFromData failed")
	// 	return nil
	// }
	return gui.NewQIcon2(pixmap)
}
