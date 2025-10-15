package application

import (
	"log/slog"

	qt "github.com/mappu/miqt/qt6"
)

func (app *Application) setupDekstopWidget() {
	slog.Info("setupDekstopWidget")
	pixmap, err := loadPNGPixmap("ayandict-64px.png")
	if err != nil {
		slog.Error("failed to load icon image", "err", err)
		return
	}
	widget := qt.NewQWidget2()
	widget.SetWindowFlags(qt.FramelessWindowHint |
		qt.WindowStaysOnBottomHint | qt.Tool | qt.BypassWindowManagerHint)

	layout := qt.NewQHBoxLayout(widget)

	label := qt.NewQLabel2()
	label.SetPixmap(pixmap)
	layout.AddWidget(label.QWidget)

	w := &DektopWidget{
		QWidget: widget,
		app:     app,
	}
	app.desktopWidget = w
	w.init()
	w.Show()
}

type DektopWidget struct {
	*qt.QWidget

	app *Application

	// dragPos: position of mouse relative to window, when drag starts
	dragPos *qt.QPoint

	// dragPosGlobal: position of mouse globally (in screen), when drag starts
	dragPosGlobal *qt.QPoint
}

func (w *DektopWidget) init() {
	w.OnMousePressEvent(w.onDragMousePress)
	w.OnMouseMoveEvent(w.onDragMouseMove)
	w.OnMouseReleaseEvent(w.onDragMouseRelease)

	// timer := qt.NewQTimer()
	// timer.SetSingleShot(true)
	// timer.OnTimeout(func() {
	// 	slog.Info("------ OnTimeout: Lower")
	// 	w.Lower()
	// })
	// timer.Start(200)
}

func (w *DektopWidget) popupMenu(event *qt.QMouseEvent) {
	menu := qt.NewQMenu2()
	for _, action := range w.app.statusIconActions {
		menu.AddAction(action)
	}
	// menu.SetFont(app.systemDefaultFont)
	menu.Popup(event.Pos())
}

func (w *DektopWidget) onDragMousePress(super func(*qt.QMouseEvent), event *qt.QMouseEvent) {
	switch event.Button() {
	case qt.LeftButton:
		w.dragPos = event.Pos()
		w.dragPosGlobal = event.GlobalPos()
	case qt.RightButton:
		w.popupMenu(event)
	default:
		super(event)
	}
}

func (w *DektopWidget) onDragMouseMove(super func(*qt.QMouseEvent), event *qt.QMouseEvent) {
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

func (w *DektopWidget) onDragMouseRelease(super func(event *qt.QMouseEvent), event *qt.QMouseEvent) {
	if event.Button() != qt.LeftButton {
		super(event)
		return
	}
	if w.dragPosGlobal != nil {
		pos := event.GlobalPos()
		dx := absInt(pos.X() - w.dragPosGlobal.X())
		dy := absInt(pos.Y() - w.dragPosGlobal.Y())
		if dx < 2 && dy < 2 {
			w.app.onStatusIconClick()
		}
		w.dragPosGlobal = nil
	}
	w.dragPos = nil
	super(event)
}
