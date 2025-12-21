package favoritebutton

import (
	"log/slog"

	"github.com/ilius/ayandict/v3/pkg/config"
	qt "github.com/mappu/miqt/qt6"
)

const maxMenuWidth = 400

type ApplicationIface interface {
	HasFavorite(term string) bool
	SetFavoriteFromFavoriteButtonMenu(term string, checked bool)
}

func NewImageFavoriteButton(
	conf *config.Config,
	onClick func(bool),
	app ApplicationIface,
) *ImageFavoriteButton {
	filename := conf.FavoriteButtonImage
	activeIcon, err := loadPNGIcon(filename)
	if err != nil {
		slog.Error("error loading icon " + filename + ": " + err.Error())
		panic(err)
	}
	inactiveIcon, err := loadPNGIcon("favorite-64.png")
	if err != nil {
		slog.Error("error loading icon favorite-64.png: " + err.Error())
		panic(err)
	}
	qButton := qt.NewQPushButton4(inactiveIcon, "")
	qButton.OnResizeEvent(func(super func(event *qt.QResizeEvent), event *qt.QResizeEvent) {
		h := event.Size().Height()
		iconSize := h * 2 / 3
		qButton.SetIconSize(qt.NewQSize2(iconSize, iconSize))
		qButton.SetFixedWidth(h)
	})
	qButton.SetStyleSheet("margin: 0px;")
	button := &ImageFavoriteButton{
		QPushButton:  qButton,
		app:          app,
		activeIcon:   activeIcon,
		inactiveIcon: inactiveIcon,
	}
	if !config.PrivateMode {
		button.OnClicked(func() {
			onClick(button.ToggleChecked())
		})
	}
	button.OnContextMenuEvent(button.onContextMenu)
	return button
}

type ImageFavoriteButton struct {
	*qt.QPushButton

	app ApplicationIface

	checked bool
	terms   []string

	activeIcon   *qt.QIcon
	inactiveIcon *qt.QIcon

	inactiveTooltip string
	activeTooltip   string
}

func (b *ImageFavoriteButton) QWidget() *qt.QWidget {
	return b.QPushButton.QWidget
}

func (b *ImageFavoriteButton) SetChecked(checked bool) {
	b.checked = checked
	if checked {
		b.SetIcon(b.activeIcon)
		b.SetToolTip(b.activeTooltip)
	} else {
		b.SetIcon(b.inactiveIcon)
		b.SetToolTip(b.inactiveTooltip)
	}
}

func (b *ImageFavoriteButton) SetTerms(terms []string) {
	b.terms = terms
}

func (b *ImageFavoriteButton) ToggleChecked() bool {
	b.SetChecked(!b.checked)
	return b.checked
}

func (b *ImageFavoriteButton) SetToolTips(inactive string, active string) {
	b.inactiveTooltip = inactive
	b.activeTooltip = active
	b.SetToolTip(inactive)
}

func (b *ImageFavoriteButton) onContextMenu(super func(event *qt.QContextMenuEvent), event *qt.QContextMenuEvent) {
	if len(b.terms) == 0 {
		return
	}
	fm := b.FontMetrics()
	menu := qt.NewQMenu2()
	menu.SetSeparatorsCollapsible(false)
	menu.SetStyleSheet(`
QMenu::item {
    padding-left: 0.25em;
    padding-right: 0.25em;
}`)
	menuWidth := 0
	mainTerm := b.terms[0]
	for _, term := range b.terms {
		action := qt.NewQAction2(term)
		action.SetCheckable(true)
		action.SetChecked(b.app.HasFavorite(term))
		action.OnTriggeredWithChecked(func(checked bool) {
			if term == mainTerm {
				b.SetChecked(checked)
			}
			b.app.SetFavoriteFromFavoriteButtonMenu(term, checked)
		})
		menu.AddAction(action)
		width := fm.HorizontalAdvance(term)
		if width > menuWidth {
			menuWidth = width
		}
	}
	menuWidth += fm.HorizontalAdvance("M")/2 + 32
	if menuWidth > maxMenuWidth {
		menuWidth = maxMenuWidth
	}
	menu.SetMinimumWidth(menuWidth)

	pos := b.Pos()
	pos.SetY(pos.Y() + b.Height())
	menu.Popup(b.ParentWidget().MapToGlobalWithQPoint(pos))
}
