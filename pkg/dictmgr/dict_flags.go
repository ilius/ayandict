package dictmgr

import (
	"log"
	"sync"

	"github.com/ilius/qt/widgets"
)

func NewDictFlagsCheckboxes(hide func()) *DictFlagsCheckboxes {
	widget := widgets.NewQWidget(nil, 0)
	hbox := widgets.NewQHBoxLayout2(nil)
	widget.SetLayout(hbox)
	hbox.SetSpacing(10) // FIXME

	w := &DictFlagsCheckboxes{
		QWidget: widget,
		hbox:    hbox,
	}

	w.addCheckBox("Fuzzy", FlagNoFuzzy)
	w.addCheckBox("Start with", FlagNoStartWith)
	w.addCheckBox("Regex", FlagNoRegex)
	w.addCheckBox("Glob", FlagNoGlob)

	hbox.AddSpacing(30) // FIXME
	hideButton := widgets.NewQPushButton2("Hide", nil)
	hideButton.ConnectClicked(func(bool) {
		hide()
	})
	hbox.AddWidget(hideButton, 0, 0)

	return w
}

type DictFlagsCheckboxes struct {
	*widgets.QWidget

	hbox *widgets.QHBoxLayout

	checkList []*widgets.QPushButton

	ds *DictSettings

	flagsMutex sync.Mutex
}

func (w *DictFlagsCheckboxes) SetActiveDictSetting(ds *DictSettings) {
	w.ds = ds
	w.checkList[0].SetChecked(ds.Fuzzy())
	w.checkList[1].SetChecked(ds.StartWith())
	w.checkList[2].SetChecked(ds.Regex())
	w.checkList[3].SetChecked(ds.Glob())
}

func (w *DictFlagsCheckboxes) addCheckBox(label string, flag uint16) {
	check := widgets.NewQPushButton2(label, nil)
	check.SetCheckable(true)
	check.SetChecked(true)
	w.hbox.AddWidget(check, 1, 0)
	w.checkList = append(w.checkList, check)
	w.checkConnectClicked(check, flag)
}

func (w *DictFlagsCheckboxes) checkConnectClicked(check *widgets.QPushButton, flag uint16) {
	check.ConnectClicked(func(checked bool) {
		ds := w.ds
		if ds == nil {
			return
		}
		w.flagsMutex.Lock()
		defer w.flagsMutex.Unlock()
		log.Printf("(before) flags = %x for %v", ds.Flags, ds.Symbol)
		if checked {
			ds.Flags &= ^flag
		} else {
			ds.Flags |= flag
		}
		log.Printf("(after)  flags = %x for %v, flag=%x", ds.Flags, ds.Symbol, flag)
	})
}
