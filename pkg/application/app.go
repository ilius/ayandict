package application

import (
	"fmt"
	"os"

	// "github.com/ilius/qt/webengine"

	"github.com/ilius/ayandict/pkg/config"
	"github.com/ilius/ayandict/pkg/dictmgr"
	"github.com/ilius/ayandict/pkg/favorites"
	"github.com/ilius/ayandict/pkg/frequency"
	"github.com/ilius/ayandict/pkg/iface"
	"github.com/ilius/ayandict/pkg/qerr"
	"github.com/ilius/ayandict/pkg/server"
	"github.com/ilius/ayandict/pkg/settings"
	"github.com/ilius/qt/core"
	"github.com/ilius/qt/gui"
	"github.com/ilius/qt/widgets"
)

const (
	QS_mainSplitter   = "main_splitter"
	QS_frequencyTable = "frequencytable"
)

// basePx is %66 of the font size in pixels,
// I'm using it for spacing between widgets
// kinda like "em" in html, but probably not exactly the same
var basePx = float32(10)

func Run() {
	app := &Application{
		QApplication:   widgets.NewQApplication(len(os.Args), os.Args),
		window:         widgets.NewQMainWindow(nil, 0),
		allTextWidgets: []iface.HasSetFont{},
	}
	qerr.ShowQtError = true

	if cacheDir == "" {
		qerr.Error(cacheDir)
	}
	{
		err := os.MkdirAll(cacheDir, 0o755)
		if err != nil {
			qerr.Error(err)
		}
	}

	app.Run()
}

type Application struct {
	*widgets.QApplication

	window *widgets.QMainWindow

	dictManager *dictmgr.DictManager

	allTextWidgets []iface.HasSetFont

	headerLabel *HeaderLabel
}

func (app *Application) Run() {
	if !LoadConfig() {
		conf = config.Default()
	}
	if len(conf.LocalServerPorts) == 0 {
		panic("config local_server_ports is empty")
	}
	client.Timeout = conf.LocalClientTimeout

	if ok, _ := findLocalServer(conf.LocalServerPorts); ok {
		qerr.Error("Another instance is running")
		return
	}
	go server.StartServer(conf.LocalServerPorts[0])

	app.LoadUserStyle()
	dictmgr.InitDicts(conf)

	basePx = float32(fontPixelSize(
		app.Font(),
		app.PrimaryScreen().PhysicalDotsPerInch(),
	) * 0.66)

	basePxI := int(basePx)
	basePxHalf := int(basePx / 2)

	frequencyTable = frequency.NewFrequencyView(
		frequencyFilePath(),
		conf.MostFrequentMaxSize,
	)

	// icon := gui.NewQIcon5("./img/icon.png")

	window := app.window
	window.SetWindowTitle("AyanDict")
	window.Resize2(600, 400)

	entry := widgets.NewQLineEdit(nil)
	entry.SetPlaceholderText("Type search query and press Enter")
	// to reduce inner margins:
	entry.SetTextMargins(0, -3, 0, -3)

	queryModeCombo := widgets.NewQComboBox(nil)
	queryModeCombo.AddItems([]string{
		"Fuzzy",
		"Start with",
		"Regex",
		"Glob",
	})

	okButton := widgets.NewQPushButton2(" OK ", nil)

	queryFavoriteButton := NewPNGIconTextButton("", "favorite.png")
	queryFavoriteButton.SetCheckable(true)
	queryFavoriteButton.SetToolTip("Add this query to favorites")

	// favoriteButtonVBox := widgets.NewQVBoxLayout()
	favoriteButton := NewPNGIconTextButton("", "favorite.png")
	favoriteButton.SetCheckable(true)
	favoriteButton.SetToolTip("Add this term to favorites")
	favoriteButton.Hide()
	// favoriteButtonVBox.AddWidget(favoriteButton, 0, core.Qt__AlignBottom)

	okButton.ConnectResizeEvent(func(event *gui.QResizeEvent) {
		h := event.Size().Height()
		if h > 100 {
			return
		}
		queryFavoriteButton.SetFixedSize2(h, h)
		favoriteButton.SetFixedSize2(h, h)
	})

	queryLabel := widgets.NewQLabel2("Query:", nil, 0)
	queryBox := widgets.NewQFrame(nil, 0)
	queryBoxLayout := widgets.NewQHBoxLayout2(queryBox)
	queryBoxLayout.SetContentsMargins(basePxHalf, basePxHalf, basePxHalf, 0)
	queryBoxLayout.SetSpacing(basePxI)
	queryBoxLayout.AddWidget(queryLabel, 0, 0)
	queryBoxLayout.AddWidget(entry, 0, 0)
	queryBoxLayout.AddWidget(queryModeCombo, 0, 0)
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
	headerBoxLayout.AddSpacing(basePxHalf)
	headerBoxLayout.AddWidget(headerLabel, 1, 0)
	// headerBoxLayout.AddLayout(favoriteButtonVBox, 0)
	headerBoxLayout.AddWidget(favoriteButton, 0, core.Qt__AlignRight)
	headerBoxLayout.AddSpacing(int(basePx * 1.5))
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
	buttonBox.SetSpacing(basePxHalf)

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
	leftMainLayout.AddSpacing(basePxHalf)
	leftMainLayout.AddWidget(headerBox, 0, 0)
	leftMainLayout.AddSpacing(basePxHalf)
	leftMainLayout.AddWidget(articleView, 0, 0)
	leftMainLayout.AddSpacing(basePxHalf)
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
		favoriteButton.Show()
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
		Entry:       entry,
		ModeCombo:   queryModeCombo,
	}

	doQuery := func(query string) {
		onQuery(query, queryArgs, false)
		entry.SetText(query)
	}
	headerLabel.doQuery = doQuery
	articleView.doQuery = doQuery
	historyView.doQuery = doQuery

	rightPanel := widgets.NewQTabWidget(nil)
	rightPanel.AddTab(activityWidget, " Activity ")
	rightPanel.AddTab(miscBox, " Misc ")

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

	app.allTextWidgets = []iface.HasSetFont{
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
	app.ReloadFont()

	resetQuery := func() {
		queryArgs.ResetQuery()
		favoriteButton.Hide()
		queryFavoriteButton.SetChecked(false)
	}

	entry.ConnectReturnPressed(func() {
		onQuery(entry.Text(), queryArgs, false)
	})
	okButton.ConnectClicked(func(bool) {
		onQuery(entry.Text(), queryArgs, false)
	})
	queryModeCombo.ConnectCurrentIndexChanged(func(i int) {
		text := entry.Text()
		if text == "" {
			return
		}
		onQuery(text, queryArgs, false)
	})
	aboutButton.ConnectClicked(func(bool) {
		aboutClicked(window)
	})

	setupKeyPressEvent := func(widget KeyPressIface) {
		widget.ConnectKeyPressEvent(func(event *gui.QKeyEvent) {
			// log.Printf("KeyPressEvent: %T", widget)
			switch event.Text() {
			case " ":
				entry.SetFocus(core.Qt__ShortcutFocusReason)
				return
			case "+", "=": // core.Qt__Key_Plus
				articleView.ZoomIn(1)
				return
			case "-": // core.Qt__Key_Minus
				articleView.ZoomOut(1)
				return
			case "\x1b": // Escape
				resetQuery()
				return
			}
			widget.KeyPressEventDefault(event)
		})
	}

	for _, widget := range []KeyPressIface{
		resultList,
		articleView,
		historyView,
	} {
		setupKeyPressEvent(widget)
	}

	articleView.SetupCustomHandlers()
	historyView.SetupCustomHandlers()

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
		app.ReloadConfig()
		onQuery(entry.Text(), queryArgs, false)
	})
	reloadStyleButton.ConnectClicked(func(checked bool) {
		app.LoadUserStyle()
		onQuery(entry.Text(), queryArgs, false)
	})
	saveHistoryButton.ConnectClicked(func(checked bool) {
		SaveHistory()
		frequencyTable.SaveNoError()
	})
	clearHistoryButton.ConnectClicked(func(checked bool) {
		historyView.ClearHistory()
		frequencyTable.Clear()
		frequencyTable.SaveNoError()
	})
	saveFavoritesButton.ConnectClicked(func(checked bool) {
		favoritesWidget.Save()
	})
	clearButton.ConnectClicked(func(checked bool) {
		resetQuery()
	})

	dictsButton.ConnectClicked(func(checked bool) {
		if app.runDictManager() {
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
		favoritesWidget.SetFavorite(term, checked)
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
		favoritesWidget.SetFavorite(term, checked)
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

func (app *Application) runDictManager() bool {
	if app.dictManager == nil {
		app.dictManager = dictmgr.NewDictManager(app.QApplication, app.window, conf)
		app.allTextWidgets = append(
			app.allTextWidgets,
			app.dictManager.TextWidgets...,
		)
	}
	return app.dictManager.Run()
}
