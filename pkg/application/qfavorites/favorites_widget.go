package qfavorites

import (
	"log/slog"

	"github.com/ilius/ayandict/v2/pkg/config"
	"github.com/ilius/ayandict/v2/pkg/favorites"
	"github.com/ilius/qt/widgets"
)

func NewFavoritesWidget(conf *config.Config) *FavoritesWidget {
	widget := widgets.NewQListWidget(nil)
	widget.ConnectItemClicked(func(item *widgets.QListWidgetItem) {
		widget.ItemActivated(item)
	})
	return &FavoritesWidget{
		QListWidget: widget,
		Data: &favorites.Favorites{
			Map: map[string]int{},
		},
		conf: conf,
	}
}

type FavoritesWidget struct {
	*widgets.QListWidget
	Data *favorites.Favorites
	conf *config.Config
}

func (w *FavoritesWidget) Load() error {
	err := w.Data.Load(favorites.Path())
	if err != nil {
		return err
	}
	for _, term := range w.Data.List {
		w.InsertItem2(0, term)
	}
	return nil
}

func (w *FavoritesWidget) Save() error {
	return w.Data.Save(favorites.Path())
}

func (w *FavoritesWidget) HasFavorite(item string) bool {
	_, ok := w.Data.Map[item]
	return ok
}

func (w *FavoritesWidget) AddFavorite(item string) {
	w.Data.Add(item)
	w.InsertItem2(0, item)
	if w.conf.FavoritesAutoSave {
		err := w.Save()
		if err != nil {
			slog.Error("error saving favorites: " + err.Error())
		}
	}
}

func (w *FavoritesWidget) RemoveFavorite(item string) {
	index := w.Data.Remove(item)
	if index < 0 {
		return
	}
	// the widget order is reversed, so our widget index
	// is N-index-1
	w.TakeItem(w.Count() - index - 1)
	if w.conf.FavoritesAutoSave {
		err := w.Save()
		if err != nil {
			slog.Error("error saving favorites: " + err.Error())
		}
	}
}

func (w *FavoritesWidget) SetFavorite(item string, favorite bool) {
	if favorite {
		w.AddFavorite(item)
	} else {
		w.RemoveFavorite(item)
	}
}
