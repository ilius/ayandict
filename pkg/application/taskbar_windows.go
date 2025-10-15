//go:build windows
// +build windows

package application

import (
	qt "github.com/mappu/miqt/qt6"
)

// Remove a Qt/miqt window from the taskbar and pager on X11.
func hideWindowFromTaskbar(widget *qt.QWidget) {}
