package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/ilius/ayandict/pkg/config"
	"github.com/ilius/ayandict/pkg/stardict"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

const dictsJsonFilename = "dicts.json"

func loadDictsOrder() (map[string]int, error) {
	m := map[string]int{}
	fpath := filepath.Join(config.GetConfigDir(), dictsJsonFilename)
	jsonBytes, err := ioutil.ReadFile(fpath)
	if err != nil {
		if os.IsNotExist(err) {
			return m, nil
		}
		return m, err
	}
	err = json.Unmarshal(jsonBytes, &m)
	if err != nil {
		return m, err
	}
	return m, nil
}

func NewDictManagerDialog(
	app *widgets.QApplication,
	parent widgets.QWidget_ITF,
) *widgets.QDialog {
	infos := stardict.GetInfoList()

	window := widgets.NewQDialog(parent, core.Qt__Dialog)
	window.Resize2(400, 400)

	listWidget := widgets.NewQListWidget(nil)

	mainHBox := widgets.NewQHBoxLayout2(window)
	mainHBox.AddWidget(listWidget, 10, core.Qt__AlignJustify)

	for _, info := range infos {
		item := widgets.NewQListWidgetItem2(info.BookName(), listWidget, 0)
		if info.Disabled {
			item.SetCheckState(core.Qt__Unchecked)
		} else {
			item.SetCheckState(core.Qt__Checked)
		}
	}

	return window
}
