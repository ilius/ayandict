package application

import (
	"log/slog"
	"time"

	"github.com/ilius/ayandict/v3/pkg/qplatform"
	qt "github.com/mappu/miqt/qt6"
)

func (app *Application) setupDekstopWidget() {
	slog.Info("setupDekstopWidget")
	pixmap, err := loadPNGPixmap(iconPixName)
	if err != nil {
		slog.Error("failed to load icon image", "err", err)
		return
	}
	widget := qt.NewQWidget2()

	flags := qt.WindowStaysOnBottomHint | qt.Window | qt.FramelessWindowHint
	widget.SetWindowFlags(flags)

	layout := qt.NewQHBoxLayout(widget)

	label := qt.NewQLabel2()
	label.SetPixmap(pixmap)
	layout.AddWidget(label.QWidget)

	w := &DesktopWidget{
		QWidget: widget,
		app:     app,
	}
	app.desktopWidget = w
	w.init()
	w.Show()
}

type DesktopWidget struct {
	*qt.QWidget

	app *Application

	// dragPos: position of mouse relative to window, when drag starts
	dragPos *qt.QPoint

	// dragPosGlobal: position of mouse globally (in screen), when drag starts
	dragPosGlobal *qt.QPoint
}

func (w *DesktopWidget) init() {
	w.OnMousePressEvent(w.onMousePress)
	w.OnMouseMoveEvent(w.onMouseMove)
	w.OnMouseReleaseEvent(w.onMouseRelease)

	go func() {
		time.Sleep(200 * time.Millisecond)
		qplatform.HideWindowFromTaskbar(w.QWidget)
	}()

	// timer := qt.NewQTimer()
	// timer.SetSingleShot(true)
	// timer.OnTimeout(func() {
	// 	slog.Info("------ OnTimeout: Lower")
	// 	w.Lower()
	// })
	// timer.Start(200)
}

func (w *DesktopWidget) popupMenu(event *qt.QMouseEvent) {
	menu := qt.NewQMenu2()
	for _, action := range w.app.statusIconActions {
		menu.AddAction(action)
	}
	// menu.SetFont(app.systemDefaultFont)
	menu.Popup(event.GlobalPos())
}

func (w *DesktopWidget) onMousePress(super func(*qt.QMouseEvent), event *qt.QMouseEvent) {
	switch event.Button() {
	case qt.RightButton:
		w.popupMenu(event)
	case qt.MiddleButton:
		w.dragPos = event.Pos()
		w.dragPosGlobal = event.GlobalPos()
		if !qplatform.CanMoveWindow() {
			w.WindowHandle().StartSystemMove()
		}
	case qt.LeftButton:
		w.dragPos = event.Pos()
		w.dragPosGlobal = event.GlobalPos()
		// On Wayland, we have to use QWindow.StartSystemMove
		// But then MouseReleaseEvent is never triggered and we can't detect a "click"
		// That's why we use a QTimer: after 100ms (can change this time with config)
		// we check if w.dragPosGlobal is the same pointer, then MouseRelease didn't
		// happen, so we assume user is trying to drag-move the window rather than clicking
		if !qplatform.CanMoveWindow() {
			posGlobal := w.dragPosGlobal
			timer := qt.NewQTimer()
			timer.SetSingleShot(true)
			timer.OnTimeout(func() {
				if w.dragPosGlobal == posGlobal {
					w.WindowHandle().StartSystemMove()
				}
			})
			timer.Start(conf.DesktopWidgetClickTime)
		}
	default:
		super(event)
	}
}

func (w *DesktopWidget) onMouseMove(super func(*qt.QMouseEvent), event *qt.QMouseEvent) {
	// event.Button() == qt.NoButton
	if w.dragPos == nil {
		super(event)
		return
	}
	w.Move(
		event.GlobalX()-w.dragPos.X(),
		event.GlobalY()-w.dragPos.Y(),
	)
}

func absInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func (w *DesktopWidget) onMouseRelease(super func(event *qt.QMouseEvent), event *qt.QMouseEvent) {
	if event.Button() != qt.LeftButton {
		super(event)
		return
	}
	if w.dragPosGlobal != nil {
		pos := event.GlobalPos()
		dx := absInt(pos.X() - w.dragPosGlobal.X())
		dy := absInt(pos.Y() - w.dragPosGlobal.Y())
		if dx < 2 && dy < 2 {
			w.app.statusIconActivate()
		}
		w.dragPosGlobal = nil
	}
	w.dragPos = nil
	super(event)
}
