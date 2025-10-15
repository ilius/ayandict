//go:build !(windows || darwin) && !wayland
// +build !windows,!darwin,!wayland

package application

/*
#cgo pkg-config: x11
#include <X11/Xlib.h>
#include <X11/Xatom.h>
#include <X11/Xutil.h>
#include <stdlib.h>
#include <string.h>

void hide_from_taskbar(unsigned long wid) {
    Display *dpy = XOpenDisplay(NULL);
    if (!dpy) return;

    Window root = DefaultRootWindow(dpy);
    Atom net_wm_state = XInternAtom(dpy, "_NET_WM_STATE", False);
    Atom skip_taskbar = XInternAtom(dpy, "_NET_WM_STATE_SKIP_TASKBAR", False);
    Atom skip_pager   = XInternAtom(dpy, "_NET_WM_STATE_SKIP_PAGER", False);

    XEvent e;
    memset(&e, 0, sizeof(e));
    e.xclient.type = ClientMessage;
    e.xclient.serial = 0;
    e.xclient.send_event = True;
    e.xclient.message_type = net_wm_state;
    e.xclient.window = (Window)wid;
    e.xclient.format = 32;
    e.xclient.data.l[0] = 1; // _NET_WM_STATE_ADD
    e.xclient.data.l[1] = skip_taskbar;
    e.xclient.data.l[2] = skip_pager;
    e.xclient.data.l[3] = 0;
    e.xclient.data.l[4] = 0;

    // Send the request to the root window
    XSendEvent(
        dpy,
        root,
        False,
        SubstructureRedirectMask | SubstructureNotifyMask,
        &e
    );

    XFlush(dpy);
    XCloseDisplay(dpy);
}
*/
import "C"

import (
	"log/slog"

	qt "github.com/mappu/miqt/qt6"
)

// Remove a Qt/miqt window from the taskbar and pager on X11.
func hideWindowFromTaskbar(widget *qt.QWidget) {
	defer func() {
		if r := recover(); r != nil {
			slog.Warn("X11: error while trying to remove window from taskbar", "err", r)
		}
	}()

	winID := widget.WinId()
	slog.Info("X11: hide_from_taskbar", "winID", winID)
	C.hide_from_taskbar(C.ulong(winID))
}
