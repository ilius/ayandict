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
