package application

import (
	"github.com/ilius/ayandict/v3/pkg/config"
)

func (app *Application) HasFavorite(term string) bool {
	return app.favoritesWidget.HasFavorite(term)
}

func (app *Application) SetFavoriteFromPopup(term string, checked bool) {
	app.favoritesWidget.SetFavorite(term, checked)
	if term == app.entry.Text() {
		app.queryFavoriteButton.SetChecked(checked)
	}
	active := app.resultList.Active()
	if active != nil && term == active.Terms()[0] {
		app.favoriteButton.SetChecked(checked)
	}
	app.resultList.Reload()
}

func (app *Application) queryFavoriteButtonClicked(checked bool) {
	if config.PrivateMode {
		return
	}
	term := app.entry.Text()
	if term == "" {
		app.queryFavoriteButton.SetChecked(false)
		return
	}
	app.favoritesWidget.SetFavorite(term, checked)
	active := app.resultList.Active()
	if active != nil && term == active.Terms()[0] {
		app.favoriteButton.SetChecked(checked)
	}
	app.resultList.Reload()
}

func (app *Application) favoriteButtonClicked(checked bool) {
	if config.PrivateMode {
		return
	}
	if app.resultList.Active == nil {
		app.favoriteButton.SetChecked(false)
		return
	}
	term := app.resultList.Active().Terms()[0]
	app.favoritesWidget.SetFavorite(term, checked)
	if term == app.entry.Text() {
		app.queryFavoriteButton.SetChecked(checked)
	}
	app.resultList.Reload()
}
