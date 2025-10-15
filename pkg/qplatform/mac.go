//go:build darwin
// +build darwin

package qplatform

import (
	qt "github.com/mappu/miqt/qt6"
)

// Remove a Qt/miqt window from the taskbar and pager on X11.
func HideWindowFromTaskbar(widget *qt.QWidget) {}
