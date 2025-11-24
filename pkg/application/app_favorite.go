package application

func (app *Application) HasFavorite(term string) bool {
	return app.favoritesWidget.HasFavorite(term)
}

func (app *Application) SetFavoriteFromPopup(term string, checked bool) {
	app.favoritesWidget.SetFavorite(term, checked)
	if term == app.entry.Text() {
		app.queryFavoriteButton.SetChecked(checked)
	}
	if app.resultList.Active != nil && term == app.resultList.Active.Terms()[0] {
		app.favoriteButton.SetChecked(checked)
	}
}

func (app *Application) queryFavoriteButtonClicked(checked bool) {
	term := app.entry.Text()
	if term == "" {
		app.queryFavoriteButton.SetChecked(false)
		return
	}
	app.favoritesWidget.SetFavorite(term, checked)
	if app.resultList.Active != nil && term == app.resultList.Active.Terms()[0] {
		app.favoriteButton.SetChecked(checked)
	}
}

func (app *Application) favoriteButtonClicked(checked bool) {
	if app.resultList.Active == nil {
		app.favoriteButton.SetChecked(false)
		return
	}
	term := app.resultList.Active.Terms()[0]
	app.favoritesWidget.SetFavorite(term, checked)
	if term == app.entry.Text() {
		app.queryFavoriteButton.SetChecked(checked)
	}
}
