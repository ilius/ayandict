//go:build !(windows || darwin) && !wayland

package qplatform

/*
#cgo pkg-config: x11
#include <X11/Xlib.h>
#include <X11/Xatom.h>
#include <X11/Xutil.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

static const long WM_EVENT_MASK = SubstructureRedirectMask | SubstructureNotifyMask;

// Xlib's default error handler prints "X Error of failed request: ..." and
// calls exit(), killing the process. This connection is opened solely for
// setting the taskbar hints, so any X error here should be ignored.
static int ignore_x_error(Display *dpy, XErrorEvent *ev) {
    return 0;
}

void hide_from_taskbar(Window wid) {
    Display *dpy = XOpenDisplay(NULL);
    if (!dpy) {
        fprintf(stderr, "Cannot open display\n");
        return;
    }

    XErrorHandler old_handler = XSetErrorHandler(ignore_x_error);

    XWindowAttributes attr;
    if (!XGetWindowAttributes(dpy, wid, &attr)) {
        XSetErrorHandler(old_handler);
        XCloseDisplay(dpy);
        return;
    }

    Atom net_wm_state = XInternAtom(dpy, "_NET_WM_STATE", False);
    Atom skip_taskbar = XInternAtom(dpy, "_NET_WM_STATE_SKIP_TASKBAR", False);
    Atom skip_pager   = XInternAtom(dpy, "_NET_WM_STATE_SKIP_PAGER", False);

    if (!net_wm_state || !skip_taskbar || !skip_pager) {
        XSetErrorHandler(old_handler);
        XCloseDisplay(dpy);
        return;
    }

    XEvent e = {0};
    e.xclient.type = ClientMessage;
    e.xclient.message_type = net_wm_state;
    e.xclient.window = wid;
    e.xclient.format = 32;
    e.xclient.data.l[0] = 1;

    e.xclient.data.l[1] = skip_taskbar;
    XSendEvent(dpy, attr.root, False, WM_EVENT_MASK, &e);

    e.xclient.data.l[1] = skip_pager;
    XSendEvent(dpy, attr.root, False, WM_EVENT_MASK, &e);

    XFlush(dpy);
    XSetErrorHandler(old_handler);
    XCloseDisplay(dpy);
}
*/
import "C"

import (
	"log/slog"
	"os"

	qt "github.com/mappu/miqt/qt6"
)

// Remove a Qt/miqt window from the taskbar and pager on X11.
func HideWindowFromTaskbar(widget *qt.QWidget) {
	defer func() {
		if r := recover(); r != nil {
			slog.Warn("X11: error while trying to remove window from taskbar", "r", r)
		}
	}()
	if os.Getenv("DISPLAY") == "" {
		slog.Debug("No DISPLAY variable; X11 not available (probably pure Wayland)")
		return
	}
	if os.Getenv("WAYLAND_DISPLAY") != "" {
		slog.Debug("Running under Wayland; X11 window ID would be invalid")
		return
	}

	winID := widget.WinId()
	slog.Debug("X11: hide_from_taskbar", "winID", winID)
	C.hide_from_taskbar(C.ulong(winID))
}

func CanMoveWindow() bool {
	return os.Getenv("WAYLAND_DISPLAY") == ""
}

// "Scan Selection" and "Scan Clipboard" work on Wayland, but only within the app.
