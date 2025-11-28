package application

import (
	qt "github.com/mappu/miqt/qt6"
)

type KeyPressIface interface {
	OnKeyPressEvent(func(func(event *qt.QKeyEvent), *qt.QKeyEvent))
}

func plaintextFromHTML(htext string) string {
	doc := qt.NewQTextDocument()
	doc.SetHtml(htext)
	return doc.ToPlainText()
}

func linedLabel(text string, right int, bottom int) *qt.QWidget {
	label := qt.NewQLabel3(text)
	frame := qt.NewQFrame2()
	frame.SetFrameShape(qt.QFrame__Box)    // The border shape
	frame.SetFrameShadow(qt.QFrame__Plain) // Flat (not sunken/raised)
	frame.SetContentsMargins(1, 1, right, bottom)
	frameLayout := qt.NewQHBoxLayout(frame.QWidget)
	frameLayout.AddWidget3(label.QWidget, 1, qt.AlignCenter)
	frameLayout.SetContentsMargins(3, 0, 3, 0)
	return frame.QWidget
}

// func posStr(pos *qt.QPoint) string {
// 	if pos == nil {
// 		return "nil"
// 	}
// 	return fmt.Sprintf("(%v, %v)", pos.X(), pos.Y())
// }

// func sizeStr(size *qt.QSize) string {
// 	if size == nil {
// 		return "nil"
// 	}
// 	return fmt.Sprintf("(%v, %v)", size.Width(), size.Height())
// }
