package desktopwidget

import (
	"log/slog"

	qt "github.com/mappu/miqt/qt6"
)

const iconPixName = "ayandict-64px.png"

func NewDekstopWidget(
	onActivate func(event *qt.QMouseEvent),
	actions []*qt.QAction,
) *DesktopWidget {
	pixmap, err := loadPNGPixmap(iconPixName)
	if err != nil {
		slog.Error("failed to load icon image", "err", err)
		return nil
	}
	widget := qt.NewQWidget2()
	widget.SetWindowFlags(qt.FramelessWindowHint |
		qt.Tool |
		qt.WindowStaysOnBottomHint)
	// widget.SetWindowFlags(qt.WindowStaysOnBottomHint |
	// 	qt.WindowStaysOnBottomHint |
	// 	qt.WindowDoesNotAcceptFocus)
	// if conf.DesktopWidgetBypassWindowManager {
	// 	flag |= qt.BypassWindowManagerHint
	// }
	// widget.SetAttribute2(qt.WA_ShowWithoutActivating, true)
	// widget.SetProperty("NET_WM_WINDOW_TYPE", qt.NewQVariant11("_NET_WM_WINDOW_TYPE_DESKTOP"))

	layout := qt.NewQHBoxLayout(widget)

	label := qt.NewQLabel2()
	label.SetPixmap(pixmap)
	layout.AddWidget(label.QWidget)

	w := &DesktopWidget{
		QWidget:    widget,
		onActivate: onActivate,
		actions:    actions,
	}
	w.init()
	return w
}

type DesktopWidget struct {
	*qt.QWidget

	// set by factory:
	onActivate func(event *qt.QMouseEvent)
	actions    []*qt.QAction

	// dragPos: position of mouse relative to window, when drag starts
	dragPos *qt.QPoint

	// dragPosGlobal: position of mouse globally (in screen), when drag starts
	dragPosGlobal *qt.QPoint
}

func (w *DesktopWidget) init() {
	w.OnMousePressEvent(w.onDragMousePress)
	w.OnMouseMoveEvent(w.onDragMouseMove)
	w.OnMouseReleaseEvent(w.onDragMouseRelease)

	// timer := qt.NewQTimer()
	// timer.SetSingleShot(true)
	// timer.OnTimeout(func() {
	// 	slog.Debug("------ OnTimeout: Lower")
	// 	w.Lower()
	// })
	// timer.Start(200)
}

func (w *DesktopWidget) popupMenu(event *qt.QMouseEvent) {
	if len(w.actions) == 0 {
		return
	}
	menu := qt.NewQMenu2()
	for _, action := range w.actions {
		menu.AddAction(action)
	}
	// menu.SetFont(app.systemDefaultFont)
	menu.Popup(event.Pos())
}

func (w *DesktopWidget) onDragMousePress(super func(*qt.QMouseEvent), event *qt.QMouseEvent) {
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

func (w *DesktopWidget) onDragMouseMove(super func(*qt.QMouseEvent), event *qt.QMouseEvent) {
	if w.dragPos == nil || event.Button() != qt.LeftButton {
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

func (w *DesktopWidget) onDragMouseRelease(super func(event *qt.QMouseEvent), event *qt.QMouseEvent) {
	if event.Button() != qt.LeftButton {
		super(event)
		return
	}
	if w.dragPosGlobal != nil {
		pos := event.GlobalPos()
		dx := absInt(pos.X() - w.dragPosGlobal.X())
		dy := absInt(pos.Y() - w.dragPosGlobal.Y())
		if dx < 2 && dy < 2 {
			w.onActivate(event)
		}
		w.dragPosGlobal = nil
	}
	w.dragPos = nil
	super(event)
}
