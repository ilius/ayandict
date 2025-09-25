package qdictmgr

import (
	"log/slog"
	"sync"

	"github.com/ilius/ayandict/v3/pkg/dictmgr/internal/dicts"
	qt "github.com/mappu/miqt/qt6"
)

func NewDictFlagsCheckboxes(hide func()) *DictFlagsCheckboxes {
	widget := qt.NewQWidget2()
	hbox := qt.NewQHBoxLayout2()
	widget.SetLayout(hbox.Layout())
	hbox.SetSpacing(10) // TODO: parameterize

	w := &DictFlagsCheckboxes{
		QWidget: widget,
		hbox:    hbox,
	}

	w.addCheckBox("Fuzzy", dicts.FlagNoFuzzy)
	w.addCheckBox("Start with", dicts.FlagNoStartWith)
	w.addCheckBox("Regex", dicts.FlagNoRegex)
	w.addCheckBox("Glob", dicts.FlagNoGlob)
	w.addCheckBox("Word Match", dicts.FlagNoWordMatch)

	hbox.AddSpacing(30) // TODO: parameterize
	hideButton := qt.NewQPushButton3("Hide")
	hideButton.OnClicked(func() {
		hide()
	})
	hbox.AddWidget(hideButton.QWidget)

	return w
}

type DictFlagsCheckboxes struct {
	*qt.QWidget

	hbox *qt.QHBoxLayout

	checkList []*qt.QPushButton

	ds *dicts.DictionarySettings

	flagsMutex sync.Mutex
}

func (w *DictFlagsCheckboxes) SetActiveDictSetting(ds *dicts.DictionarySettings) {
	w.ds = ds
	w.checkList[0].SetChecked(ds.Fuzzy())
	w.checkList[1].SetChecked(ds.StartWith())
	w.checkList[2].SetChecked(ds.Regex())
	w.checkList[3].SetChecked(ds.Glob())
}

func (w *DictFlagsCheckboxes) addCheckBox(label string, flag uint16) {
	check := qt.NewQPushButton3(label)
	check.SetCheckable(true)
	check.SetChecked(true)
	w.hbox.AddWidget3(check.QWidget, 1, 0)
	w.checkList = append(w.checkList, check)
	w.checkConnectClicked(check, flag)
}

func (w *DictFlagsCheckboxes) checkConnectClicked(check *qt.QPushButton, flag uint16) {
	check.OnClicked(func() {
		ds := w.ds
		if ds == nil {
			return
		}
		w.flagsMutex.Lock()
		defer w.flagsMutex.Unlock()
		slog.Debug("(before) flags", "flags", ds.Flags, "symbol", ds.Symbol)
		if check.IsChecked() {
			ds.Flags &= ^flag
		} else {
			ds.Flags |= flag
		}
		slog.Debug("(after)  flags", "flags", ds.Flags, "symbol", ds.Symbol, "flag", flag)
	})
}
