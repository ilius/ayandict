package main

import (
	"os"

	qt "github.com/mappu/miqt/qt6"
)

func main() {
	_ = qt.NewQApplication(os.Args[1:])

	sourceEdit := qt.NewQTextEdit2()
	htmlBrowser := qt.NewQTextBrowser(nil)

	sourceEdit.OnTextChanged(func() {
		htmlBrowser.SetHtml(sourceEdit.ToPlainText())
	})

	window := qt.NewQWidget(nil)
	window.SetWindowTitle("QTextBrowser Preview")
	window.Resize(800, 800)
	layout := qt.NewQVBoxLayout2()
	window.SetLayout(layout.Layout())
	layout.AddWidget(sourceEdit.QWidget)
	layout.AddSpacing(20)
	layout.AddWidget(htmlBrowser.QWidget)

	window.Show()
	_ = qt.QApplication_Exec()
}
