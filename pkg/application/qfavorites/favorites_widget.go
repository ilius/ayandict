package qfavorites

import (
	"log/slog"

	"github.com/ilius/ayandict/v3/pkg/config"
	"github.com/ilius/ayandict/v3/pkg/favorites"
	qt "github.com/mappu/miqt/qt6"
)

func NewFavoritesWidget(conf *config.Config) *FavoritesWidget {
	widget := qt.NewQListWidget(nil)
	widget.OnItemClicked(func(item *qt.QListWidgetItem) {
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
	*qt.QListWidget
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
	if config.PrivateMode {
		return nil
	}
	return w.Data.Save(favorites.Path())
}

func (w *FavoritesWidget) HasFavorite(item string) bool {
	_, ok := w.Data.Map[item]
	return ok
}

func (w *FavoritesWidget) AddFavorite(item string) {
	if config.PrivateMode {
		return
	}
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
	if config.PrivateMode {
		return
	}
	index := w.Data.Remove(item)
	if index < 0 {
		return
	}
	// the widget order is reversed, so our widget index
	// is N-index-1
	_ = w.TakeItem(w.Count() - index - 1)
	if w.conf.FavoritesAutoSave {
		err := w.Save()
		if err != nil {
			slog.Error("error saving favorites: " + err.Error())
		}
	}
}

func (w *FavoritesWidget) SetFavorite(item string, favorite bool) {
	if config.PrivateMode {
		return
	}
	if favorite {
		w.AddFavorite(item)
	} else {
		w.RemoveFavorite(item)
	}
}
