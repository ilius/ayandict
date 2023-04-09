package application

import (
	"fmt"
	"log"
	"os"

	// "github.com/therecipe/qt/webengine"

	"github.com/ilius/ayandict/pkg/favorites"
	"github.com/ilius/ayandict/pkg/frequency"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

var dictManager *DictManager

func Run() {
	app := widgets.NewQApplication(len(os.Args), os.Args)
	LoadConfig(app)

	LoadUserStyle(app)
	initDicts()

	frequencyTable = frequency.NewFrequencyView(conf.MostFrequentMaxSize)

	// icon := gui.NewQIcon5("./img/icon.png")

	window := widgets.NewQMainWindow(nil, 0)
	window.SetWindowTitle("AyanDict")
	window.Resize2(600, 400)

	entry := widgets.NewQLineEdit(nil)
	entry.SetPlaceholderText("Type search query and press Enter")
	// to reduce inner margins:
	entry.SetTextMargins(0, -3, 0, -3)

	okButton := widgets.NewQPushButton2("OK", nil)

	queryFavoriteButton := NewPNGIconTextButton("", "favorite.png")
	queryFavoriteButton.SetCheckable(true)

	queryBox := widgets.NewQFrame(nil, 0)
	queryBoxLayout := widgets.NewQHBoxLayout2(queryBox)
	queryBoxLayout.SetContentsMargins(5, 5, 5, 0)
	queryBoxLayout.SetSpacing(10)
	queryBoxLayout.AddWidget(widgets.NewQLabel2("Query:", nil, 0), 0, 0)
	queryBoxLayout.AddWidget(entry, 0, 0)
	queryBoxLayout.AddWidget(queryFavoriteButton, 0, 0)
	queryBoxLayout.AddWidget(okButton, 0, 0)
	// queryBoxLayout.SetSpacing(10)

	headerLabel := CreateHeaderLabel(app)

	// FIXME: putting headerLabel in a HBox while WordWrap is on
	// makes it not expand. Since I could not fix this, I'm putting
	// Favorite button to the bottomBox for now
	// headerBox := widgets.NewQWidget(nil, 0)
	// headerBoxLayout := widgets.NewQHBoxLayout2(headerBox)
	// headerBoxLayout.SetSizeConstraint(widgets.QLayout__SetMinimumSize)
	// headerBoxLayout.SetContentsMargins(0, 0, 0, 0)
	// headerBoxLayout.SetSpacing(10)
	// headerBoxLayout.AddWidget(headerLabel, 1, core.Qt__AlignLeft)
	// headerBox.SetSizePolicy2(expanding, widgets.QSizePolicy__Minimum)
	// headerBoxLayout.AddWidget(favoriteButton, 0, core.Qt__AlignLeft)
	// favoriteButton.SetSizePolicy2(widgets.QSizePolicy__Minimum, widgets.QSizePolicy__Minimum)

	articleView := NewArticleView(app)

	historyView := NewHistoryView()

	frequencyTable.SetHorizontalHeaderItem(
		0,
		widgets.NewQTableWidgetItem2("Query", 0),
	)
	frequencyTable.SetHorizontalHeaderItem(
		1,
		widgets.NewQTableWidgetItem2("Count", 0),
	)
	if !conf.MostFrequentDisable {
		frequencyTable.LoadFromFile(frequencyFilePath())
	}
	// TODO: save the width of 2 columns

	favoritesWidget := favorites.NewFavoritesWidget(conf)
	{
		err := favoritesWidget.Load()
		if err != nil {
			// conf.FavoritesAutoSave = false
			fmt.Println(err)
		}
	}

	miscBox := widgets.NewQFrame(nil, 0)
	miscLayout := widgets.NewQVBoxLayout2(miscBox)
	miscLayout.SetContentsMargins(0, 0, 0, 0)

	saveHistoryButton := widgets.NewQPushButton2("Save History", nil)
	miscLayout.AddWidget(saveHistoryButton, 0, 0)
	clearHistoryButton := widgets.NewQPushButton2("Clear History", nil)
	miscLayout.AddWidget(clearHistoryButton, 0, 0)

	reloadDictsButton := widgets.NewQPushButton2("Reload Dicts", nil)
	miscLayout.AddWidget(reloadDictsButton, 0, 0)
	closeDictsButton := widgets.NewQPushButton2("Close Dicts", nil)
	miscLayout.AddWidget(closeDictsButton, 0, 0)
	reloadStyleButton := widgets.NewQPushButton2("Reload Style", nil)
	miscLayout.AddWidget(reloadStyleButton, 0, 0)

	bottomBox := widgets.NewQHBoxLayout()
	bottomBox.SetContentsMargins(0, 0, 0, 0)
	bottomBox.SetSpacing(10)

	bottomBoxStyleOpt := widgets.NewQStyleOptionButton()
	style := app.Style()

	newIconTextButton := func(label string, pix widgets.QStyle__StandardPixmap) *widgets.QPushButton {
		return widgets.NewQPushButton3(
			style.StandardIcon(
				pix, bottomBoxStyleOpt, nil,
			),
			label, nil,
		)
	}

	favoriteButton := NewPNGIconTextButton("Favorite", "favorite.png")

	favoriteButton.SetCheckable(true)
	bottomBox.AddWidget(favoriteButton, 0, core.Qt__AlignLeft)

	dictsButton := newIconTextButton("Dictionaries", widgets.QStyle__SP_FileDialogDetailedView)
	bottomBox.AddWidget(dictsButton, 0, core.Qt__AlignLeft)

	aboutButton := newIconTextButton("About", widgets.QStyle__SP_MessageBoxInformation)
	bottomBox.AddWidget(aboutButton, 0, core.Qt__AlignLeft)

	bottomBox.AddStretch(1)

	openConfigButton := NewPNGIconTextButton("Config", "preferences-system-22.png")
	bottomBox.AddWidget(openConfigButton, 0, 0)
	reloadConfigButton := newIconTextButton("Reload", widgets.QStyle__SP_BrowserReload)
	bottomBox.AddWidget(reloadConfigButton, 0, 0)

	bottomBox.AddStretch(1)

	clearButton := widgets.NewQPushButton2("Clear", nil)
	bottomBox.AddWidget(clearButton, 0, core.Qt__AlignRight)

	leftMainWidget := widgets.NewQWidget(nil, 0)
	leftMainLayout := widgets.NewQVBoxLayout2(leftMainWidget)
	leftMainLayout.SetContentsMargins(0, 0, 0, 0)
	leftMainLayout.SetSpacing(0)
	leftMainLayout.AddWidget(queryBox, 0, 0)
	leftMainLayout.AddSpacing(5)
	leftMainLayout.AddWidget(headerLabel, 0, 0)
	leftMainLayout.AddSpacing(5)
	leftMainLayout.AddWidget(articleView, 0, 0)
	leftMainLayout.AddSpacing(5)
	leftMainLayout.AddLayout(bottomBox, 0)

	activityTypeCombo := widgets.NewQComboBox(nil)
	activityTypeCombo.AddItems([]string{
		"Recent",
		"Most Frequent",
		"Favorites",
	})

	frequencyTable.Hide()
	favoritesWidget.Hide()

	activityWidget := widgets.NewQWidget(nil, 0)
	activityLayout := widgets.NewQVBoxLayout2(activityWidget)
	activityLayout.SetContentsMargins(5, 5, 5, 5)
	activityLayout.AddWidget(activityTypeCombo, 0, 0)
	activityLayout.AddWidget(historyView, 0, 0)
	activityLayout.AddWidget(frequencyTable, 0, 0)
	activityLayout.AddWidget(favoritesWidget, 0, 0)

	activityTypeCombo.ConnectCurrentIndexChanged(func(index int) {
		switch index {
		case 0:
			historyView.Show()
			frequencyTable.Hide()
			favoritesWidget.Hide()
		case 1:
			historyView.Hide()
			frequencyTable.Show()
			favoritesWidget.Hide()
		case 2:
			historyView.Hide()
			frequencyTable.Hide()
			favoritesWidget.Show()
		}
	})

	onResultDisplay := func(terms []string) {
		isFav := favoritesWidget.HasFavorite(terms[0])
		favoriteButton.SetChecked(isFav)
	}

	leftPanel := widgets.NewQWidget(nil, 0)
	leftPanelLayout := widgets.NewQVBoxLayout2(leftPanel)
	leftPanelLayout.AddWidget(widgets.NewQLabel2("Results", nil, 0), 0, 0)
	resultList := NewResultListWidget(
		articleView,
		headerLabel,
		onResultDisplay,
	)
	leftPanelLayout.AddWidget(resultList, 0, 0)

	postQuery := func(query string) {
		if query == "" {
			queryFavoriteButton.SetChecked(false)
			return
		}
		queryFavoriteButton.SetChecked(favoritesWidget.HasFavorite(query))
	}

	queryArgs := &QueryArgs{
		ArticleView: articleView,
		ResultList:  resultList,
		HeaderLabel: headerLabel,
		HistoryView: historyView,
		PostQuery:   postQuery,
	}

	doQuery := func(query string) {
		onQuery(query, queryArgs, false)
		entry.SetText(query)
	}
	headerLabel.doQuery = doQuery
	articleView.doQuery = doQuery
	historyView.doQuery = doQuery

	rightPanel := widgets.NewQTabWidget(nil)
	rightPanel.AddTab(activityWidget, "Activity")
	rightPanel.AddTab(miscBox, "Misc")

	mainSplitter := widgets.NewQSplitter(nil)
	mainSplitter.SetSizePolicy2(expanding, expanding)
	mainSplitter.AddWidget(leftPanel)
	mainSplitter.AddWidget(leftMainWidget)
	mainSplitter.AddWidget(rightPanel)
	mainSplitter.SetStretchFactor(0, 1)
	mainSplitter.SetStretchFactor(1, 5)
	mainSplitter.SetStretchFactor(2, 1)

	window.SetCentralWidget(mainSplitter)

	app.SetFont(ConfigFont(), "")

	resetQuery := func() {
		entry.SetText("")
		resultList.Clear()
		articleView.SetHtml("")
		headerLabel.SetText("")
	}

	entry.ConnectReturnPressed(func() {
		onQuery(entry.Text(), queryArgs, false)
	})
	okButton.ConnectClicked(func(bool) {
		onQuery(entry.Text(), queryArgs, false)
	})
	aboutButton.ConnectClicked(func(bool) {
		aboutClicked(window)
	})
	resultList.ConnectKeyPressEvent(func(event *gui.QKeyEvent) {
		switch event.Text() {
		case " ":
			entry.SetFocus(core.Qt__ShortcutFocusReason)
			return
		}
		resultList.KeyPressEventDefault(event)
	})

	articleView.SetupCustomHandlers()
	articleView.ConnectKeyPressEvent(func(event *gui.QKeyEvent) {
		switch event.Text() {
		case " ":
			entry.SetFocus(core.Qt__ShortcutFocusReason)
		}
	})

	historyView.SetupCustomHandlers()
	historyView.ConnectKeyPressEvent(func(event *gui.QKeyEvent) {
		switch event.Text() {
		case " ":
			entry.SetFocus(core.Qt__ShortcutFocusReason)
			return
		}
		historyView.KeyPressEventDefault(event)
	})

	frequencyTable.ConnectItemActivated(func(item *widgets.QTableWidgetItem) {
		key := frequencyTable.Keys[item.Row()]
		doQuery(key)
		newRow := frequencyTable.KeyMap[key]
		// item.Column() panics!
		frequencyTable.SetCurrentCell(newRow, 0)
	})
	favoritesWidget.ConnectItemActivated(func(item *widgets.QListWidgetItem) {
		doQuery(item.Text())
	})
	reloadDictsButton.ConnectClicked(func(checked bool) {
		reloadDicts()
	})
	closeDictsButton.ConnectClicked(func(checked bool) {
		closeDicts()
	})
	openConfigButton.ConnectClicked(func(checked bool) {
		OpenConfig()
	})
	reloadConfigButton.ConnectClicked(func(checked bool) {
		ReloadConfig(app)
		onQuery(entry.Text(), queryArgs, false)
	})
	reloadStyleButton.ConnectClicked(func(checked bool) {
		LoadUserStyle(app)
		onQuery(entry.Text(), queryArgs, false)
	})
	saveHistoryButton.ConnectClicked(func(checked bool) {
		SaveHistory()
		SaveFrequency()
	})
	clearHistoryButton.ConnectClicked(func(checked bool) {
		historyView.ClearHistory()
		frequencyTable.Clear()
		SaveFrequency()
	})
	clearButton.ConnectClicked(func(checked bool) {
		resetQuery()
	})

	dictsButton.ConnectClicked(func(checked bool) {
		if dictManager == nil {
			dictManager = NewDictManager(app, window)
		}
		if dictManager.Dialog.Exec() == dialogAccepted {
			SaveDictManagerDialog(dictManager)
			onQuery(entry.Text(), queryArgs, false)
		}
	})

	if !conf.HistoryDisable {
		err := LoadHistory()
		if err != nil {
			log.Println(err)
		} else {
			historyView.AddHistoryList(history)
		}
	}

	entry.ConnectKeyPressEvent(func(event *gui.QKeyEvent) {
		entry.KeyPressEventDefault(event)
		switch event.Text() {
		case "", "\b":
			return
		case "\x1b":
			// Escape, is there a more elegant way?
			resetQuery()
			return
		}
		if conf.SearchOnType {
			text := entry.Text()
			if len(text) < conf.SearchOnTypeMinLength {
				return
			}
			onQuery(text, queryArgs, true)
		}
	})

	favoriteButton.ConnectClicked(func(checked bool) {
		if resultList.Active == nil {
			favoriteButton.SetChecked(false)
			return
		}
		term := resultList.Active.Terms()[0]
		if checked {
			favoritesWidget.AddFavorite(term)
		} else {
			favoritesWidget.RemoveFavorite(term)
		}
		if term == entry.Text() {
			queryFavoriteButton.SetChecked(checked)
		}
	})
	queryFavoriteButton.ConnectClicked(func(checked bool) {
		term := entry.Text()
		if term == "" {
			queryFavoriteButton.SetChecked(false)
			return
		}
		if checked {
			favoritesWidget.AddFavorite(term)
		} else {
			favoritesWidget.RemoveFavorite(term)
		}
		if resultList.Active != nil && term == resultList.Active.Terms()[0] {
			favoriteButton.SetChecked(checked)
		}
	})

	qs := getQSettings(window)
	restoreSplitterSizes(qs, mainSplitter, QS_mainSplitter)
	restoreMainWinGeometry(app, qs, window)
	window.ConnectResizeEvent(func(event *gui.QResizeEvent) {
		saveMainWinGeometry(qs, window)
	})
	window.ConnectMoveEvent(func(event *gui.QMoveEvent) {
		saveMainWinGeometry(qs, window)
	})
	restoreTableColumnsWidth(
		qs,
		frequencyTable.QTableWidget,
		QS_frequencyTable,
	)
	// frequencyTable.ConnectColumnResized does not work
	frequencyTable.HorizontalHeader().ConnectSectionResized(func(logicalIndex int, oldSize int, newSize int) {
		saveTableColumnsWidth(qs, frequencyTable.QTableWidget, QS_frequencyTable)
	})

	setupSplitterSizesSave(qs, mainSplitter, QS_mainSplitter)

	window.Show()
	app.Exec()
}
