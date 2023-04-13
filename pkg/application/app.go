package application

import (
	"fmt"
	"os"
	"time"

	// "github.com/therecipe/qt/webengine"

	"github.com/ilius/ayandict/pkg/common"
	"github.com/ilius/ayandict/pkg/config"
	"github.com/ilius/ayandict/pkg/dictmgr"
	"github.com/ilius/ayandict/pkg/favorites"
	"github.com/ilius/ayandict/pkg/frequency"
	"github.com/ilius/ayandict/pkg/qerr"
	"github.com/ilius/ayandict/pkg/settings"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

const (
	QS_mainSplitter   = "main_splitter"
	QS_frequencyTable = "frequencytable"
)

func Run() {
	app := &Application{
		QApplication:   widgets.NewQApplication(len(os.Args), os.Args),
		allTextWidgets: []common.HasSetFont{},
	}
	qerr.ShowQtError = true
	app.Run()
}

type Application struct {
	*widgets.QApplication

	dictManager *dictmgr.DictManager

	allTextWidgets []common.HasSetFont

	headerLabel *HeaderLabel
}

func (app *Application) Run() {
	if !LoadConfig() {
		conf = config.Default()
	}
	if len(conf.LocalServerPorts) == 0 {
		panic("config local_server_ports is empty")
	}
	if conf.LocalClientTimeout != "" {
		timeout, err := time.ParseDuration(conf.LocalClientTimeout)
		if err != nil {
			qerr.Errorf("bad local_client_timeout=%v", conf.LocalClientTimeout)
		} else if timeout > 0 {
			client.Timeout = timeout
		}
	}

	if isSingleInstanceRunning(APP_NAME, conf.LocalServerPorts) {
		qerr.Error("Another instance is running")
		return
	}
	go startSingleInstanceServer(APP_NAME, conf.LocalServerPorts[0])

	LoadUserStyle(app)
	dictmgr.InitDicts(conf)

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
	queryFavoriteButton.SetToolTip("Add this query to favorites")

	// favoriteButtonVBox := widgets.NewQVBoxLayout()
	favoriteButton := NewPNGIconTextButton("", "favorite.png")
	favoriteButton.SetCheckable(true)
	favoriteButton.SetToolTip("Add this term to favorites")
	// favoriteButtonVBox.AddWidget(favoriteButton, 0, core.Qt__AlignBottom)

	queryLabel := widgets.NewQLabel2("Query:", nil, 0)
	queryBox := widgets.NewQFrame(nil, 0)
	queryBoxLayout := widgets.NewQHBoxLayout2(queryBox)
	queryBoxLayout.SetContentsMargins(5, 5, 5, 0)
	queryBoxLayout.SetSpacing(10)
	queryBoxLayout.AddWidget(queryLabel, 0, 0)
	queryBoxLayout.AddWidget(entry, 0, 0)
	queryBoxLayout.AddWidget(queryFavoriteButton, 0, 0)
	queryBoxLayout.AddWidget(okButton, 0, 0)

	headerLabel := CreateHeaderLabel(app)
	app.headerLabel = headerLabel
	headerLabel.SetAlignment(core.Qt__AlignLeft)

	headerBox := widgets.NewQWidget(nil, 0)
	headerBox.SetSizePolicy2(widgets.QSizePolicy__Preferred, widgets.QSizePolicy__Minimum)
	headerBoxLayout := widgets.NewQHBoxLayout2(headerBox)
	// headerBoxLayout.SetSizeConstraint(widgets.QLayout__SetMinimumSize)
	headerBoxLayout.SetContentsMargins(0, 0, 0, 0)
	// headerBoxLayout.AddLayout(favoriteButtonVBox, 0)
	headerBoxLayout.AddWidget(favoriteButton, 0, core.Qt__AlignCenter)
	headerBoxLayout.AddSpacing(10)
	headerBoxLayout.AddWidget(headerLabel, 1, 0)
	// it is very important that last argument ^ above is 0
	// otherwise label will not expand (while in layout and in wrap mode)
	// don't ask me why!
	headerBox.SetSizePolicy2(expanding, widgets.QSizePolicy__Minimum)

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

	saveFavoritesButton := widgets.NewQPushButton2("Save Favorites", nil)
	miscLayout.AddWidget(saveFavoritesButton, 0, 0)

	reloadDictsButton := widgets.NewQPushButton2("Reload Dicts", nil)
	miscLayout.AddWidget(reloadDictsButton, 0, 0)
	closeDictsButton := widgets.NewQPushButton2("Close Dicts", nil)
	miscLayout.AddWidget(closeDictsButton, 0, 0)
	reloadStyleButton := widgets.NewQPushButton2("Reload Style", nil)
	miscLayout.AddWidget(reloadStyleButton, 0, 0)

	buttonBox := widgets.NewQHBoxLayout()
	buttonBox.SetContentsMargins(0, 0, 0, 0)
	buttonBox.SetSpacing(5)

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

	dictsButtonLabel := "Dictionaries"
	if conf.ReduceMinimumWindowWidth {
		dictsButtonLabel = "Dicts"
	}
	dictsButton := newIconTextButton(dictsButtonLabel, widgets.QStyle__SP_FileDialogDetailedView)
	buttonBox.AddWidget(dictsButton, 0, core.Qt__AlignLeft)

	aboutButtonLabel := "About"
	if conf.ReduceMinimumWindowWidth {
		aboutButtonLabel = "\u200c"
	}
	aboutButton := newIconTextButton(aboutButtonLabel, widgets.QStyle__SP_MessageBoxInformation)
	buttonBox.AddWidget(aboutButton, 0, core.Qt__AlignLeft)

	buttonBox.AddStretch(1)

	openConfigButton := NewPNGIconTextButton("Config", "preferences-system-22.png")
	buttonBox.AddWidget(openConfigButton, 0, 0)
	reloadConfigButton := newIconTextButton("Reload", widgets.QStyle__SP_BrowserReload)
	buttonBox.AddWidget(reloadConfigButton, 0, 0)

	buttonBox.AddStretch(1)

	clearButton := widgets.NewQPushButton2("Clear", nil)
	buttonBox.AddWidget(clearButton, 0, core.Qt__AlignRight)

	leftMainWidget := widgets.NewQWidget(nil, 0)
	leftMainLayout := widgets.NewQVBoxLayout2(leftMainWidget)
	leftMainLayout.SetContentsMargins(0, 0, 0, 0)
	leftMainLayout.SetSpacing(0)
	leftMainLayout.AddWidget(queryBox, 0, 0)
	leftMainLayout.AddSpacing(5)
	leftMainLayout.AddWidget(headerBox, 0, 0)
	leftMainLayout.AddSpacing(5)
	leftMainLayout.AddWidget(articleView, 0, 0)
	leftMainLayout.AddSpacing(5)
	leftMainLayout.AddLayout(buttonBox, 0)

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

	app.allTextWidgets = []common.HasSetFont{
		queryLabel,
		entry,
		okButton,
		headerLabel,
		articleView,
		historyView,
		frequencyTable,
		favoritesWidget,
		saveHistoryButton,
		clearHistoryButton,
		saveFavoritesButton,
		reloadDictsButton,
		closeDictsButton,
		reloadStyleButton,
		dictsButton,
		aboutButton,
		openConfigButton,
		reloadConfigButton,
		clearButton,
		activityTypeCombo,
		resultList,
		rightPanel,
	}

	resetQuery := func() {
		entry.SetText("")
		resultList.Clear()
		articleView.SetHtml("")
		headerLabel.SetText("")
		favoriteButton.SetChecked(false)
		queryFavoriteButton.SetChecked(false)
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

	commonKeyPressEvent := func(event *gui.QKeyEvent) bool {
		switch event.Text() {
		case " ":
			entry.SetFocus(core.Qt__ShortcutFocusReason)
			return true
		case "+", "=": // core.Qt__Key_Plus
			articleView.ZoomIn(1)
			return true
		case "-": // core.Qt__Key_Minus
			articleView.ZoomOut(1)
			return true
		case "\x1b": // Escape
			resetQuery()
			return true
		}
		return false
	}

	resultList.ConnectKeyPressEvent(func(event *gui.QKeyEvent) {
		if commonKeyPressEvent(event) {
			return
		}
		resultList.KeyPressEventDefault(event)
	})

	articleView.SetupCustomHandlers()
	articleView.ConnectKeyPressEvent(func(event *gui.QKeyEvent) {
		if commonKeyPressEvent(event) {
			return
		}
		switch event.Key() {
		case int(core.Qt__Key_Up):
			if conf.ArticleArrowKeys {
				articleView.VerticalScrollBar().TriggerAction(widgets.QAbstractSlider__SliderSingleStepSub)
				return
			}
		case int(core.Qt__Key_Down):
			if conf.ArticleArrowKeys {
				articleView.VerticalScrollBar().TriggerAction(widgets.QAbstractSlider__SliderSingleStepAdd)
				return
			}
		}
		articleView.KeyPressEventDefault(event)
	})

	historyView.SetupCustomHandlers()
	historyView.ConnectKeyPressEvent(func(event *gui.QKeyEvent) {
		if commonKeyPressEvent(event) {
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
		dictmgr.InitDicts(conf)
		app.dictManager = nil
		onQuery(entry.Text(), queryArgs, false)
	})
	closeDictsButton.ConnectClicked(func(checked bool) {
		dictmgr.CloseDicts()
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
	saveFavoritesButton.ConnectClicked(func(checked bool) {
		favoritesWidget.Save()
	})
	clearButton.ConnectClicked(func(checked bool) {
		resetQuery()
	})

	dictsButton.ConnectClicked(func(checked bool) {
		if app.dictManager == nil {
			app.dictManager = dictmgr.NewDictManager(app.QApplication, window, conf)
			app.allTextWidgets = append(
				app.allTextWidgets,
				app.dictManager.TextWidgets...,
			)
		}
		if app.dictManager.Run() {
			onQuery(entry.Text(), queryArgs, false)
		}
	})

	if !conf.HistoryDisable {
		err := LoadHistory()
		if err != nil {
			qerr.Error(err)
		} else {
			historyView.AddHistoryList(history)
		}
	}

	entry.ConnectKeyPressEvent(func(event *gui.QKeyEvent) {
		entry.KeyPressEventDefault(event)
		switch event.Text() {
		case "", "\b":
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

	qs := settings.GetQSettings(window)
	settings.RestoreSplitterSizes(qs, mainSplitter, QS_mainSplitter)
	settings.RestoreMainWinGeometry(app.QApplication, qs, window)
	settings.SetupMainWinGeometrySave(qs, window)

	settings.RestoreTableColumnsWidth(
		qs,
		frequencyTable.QTableWidget,
		QS_frequencyTable,
	)
	// frequencyTable.ConnectColumnResized does not work
	frequencyTable.HorizontalHeader().ConnectSectionResized(func(logicalIndex int, oldSize int, newSize int) {
		settings.SaveTableColumnsWidth(qs, frequencyTable.QTableWidget, QS_frequencyTable)
	})

	settings.SetupSplitterSizesSave(qs, mainSplitter, QS_mainSplitter)

	window.Show()
	app.Exec()
}
