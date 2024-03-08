package application

import (
	"fmt"
	"os"

	// "github.com/ilius/qt/webengine"

	"github.com/ilius/ayandict/v2/pkg/appinfo"
	"github.com/ilius/ayandict/v2/pkg/config"
	"github.com/ilius/ayandict/v2/pkg/dictmgr"
	"github.com/ilius/ayandict/v2/pkg/dictmgr/qdictmgr"
	"github.com/ilius/ayandict/v2/pkg/favorites"
	"github.com/ilius/ayandict/v2/pkg/frequency"
	"github.com/ilius/ayandict/v2/pkg/qerr"
	"github.com/ilius/ayandict/v2/pkg/qsettings"
	"github.com/ilius/ayandict/v2/pkg/qutils"
	"github.com/ilius/ayandict/v2/pkg/server"
	common "github.com/ilius/go-dict-commons"
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
		allTextWidgets: []qutils.HasSetFont{},
	}
	qerr.ShowMessage = showErrorMessage

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

	dictManager *qdictmgr.DictManager

	allTextWidgets []qutils.HasSetFont

	queryArgs       *QueryArgs
	headerLabel     *HeaderLabel
	articleView     *ArticleView
	resultList      *ResultListWidget
	historyView     *HistoryView
	entry           *widgets.QLineEdit
	queryModeCombo  *widgets.QComboBox
	favoritesWidget *favorites.FavoritesWidget

	favoriteButton      *widgets.QPushButton
	queryFavoriteButton *widgets.QPushButton
	reloadDictsButton   *widgets.QPushButton
	closeDictsButton    *widgets.QPushButton
	openConfigButton    *widgets.QPushButton
	reloadConfigButton  *widgets.QPushButton
	reloadStyleButton   *widgets.QPushButton
	saveHistoryButton   *widgets.QPushButton
	randomEntryButton   *widgets.QPushButton
	clearHistoryButton  *widgets.QPushButton
	saveFavoritesButton *widgets.QPushButton
	clearButton         *widgets.QPushButton
	dictsButton         *widgets.QPushButton
}

func (app *Application) init() {
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
	qdictmgr.InitDicts(conf, true)
}

func (app *Application) doQuery(query string) {
	onQuery(query, app.queryArgs, false)
	app.entry.SetText(query)
}

func (app *Application) runDictManager() bool {
	if app.dictManager == nil {
		app.dictManager = qdictmgr.NewDictManager(app.QApplication, app.window, conf)
		app.allTextWidgets = append(
			app.allTextWidgets,
			app.dictManager.TextWidgets...,
		)
	}
	return app.dictManager.Run()
}

func (app *Application) resetQuery() {
	app.entry.SetText("")
	app.resultList.Clear()
	app.headerLabel.SetText("")
	app.articleView.SetHtml("")
	app.favoriteButton.Hide()
	app.queryFavoriteButton.SetChecked(false)
}

func (app *Application) postQuery(query string) {
	if query == "" {
		app.queryFavoriteButton.SetChecked(false)
		return
	}
	app.queryFavoriteButton.SetChecked(app.favoritesWidget.HasFavorite(query))
}

func (app *Application) setupKeyPressEvent(widget KeyPressIface) {
	widget.ConnectKeyPressEvent(func(event *gui.QKeyEvent) {
		// log.Printf("KeyPressEvent: %T", widget)
		switch event.Text() {
		case " ":
			app.entry.SetFocus(core.Qt__ShortcutFocusReason)
			return
		case "+", "=": // core.Qt__Key_Plus
			app.articleView.ZoomIn()
			return
		case "-": // core.Qt__Key_Minus
			app.articleView.ZoomOut()
			return
		case "\x1b": // Escape
			app.resetQuery()
			return
		}
		widget.KeyPressEventDefault(event)
	})
}

// TODO: break down
func (app *Application) Run() {
	app.init()

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
	window.SetWindowTitle(appinfo.APP_DESC)
	window.Resize2(600, 400)

	entry := widgets.NewQLineEdit(nil)
	app.entry = entry
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
	app.queryModeCombo = queryModeCombo

	okButton := widgets.NewQPushButton2(" OK ", nil)

	queryFavoriteButton := NewPNGIconTextButton("", "favorite.png")
	app.queryFavoriteButton = queryFavoriteButton
	queryFavoriteButton.SetCheckable(true)
	queryFavoriteButton.SetToolTip("Add this query to favorites")

	// favoriteButtonVBox := widgets.NewQVBoxLayout()
	favoriteButton := NewPNGIconTextButton("", "favorite.png")
	app.favoriteButton = favoriteButton
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
	app.articleView = articleView

	historyView := NewHistoryView()
	app.historyView = historyView

	frequencyTable.SetHorizontalHeaderItem(
		0,
		widgets.NewQTableWidgetItem2("Query", 0),
	)
	frequencyTable.SetHorizontalHeaderItem(
		1,
		widgets.NewQTableWidgetItem2("Count", 0),
	)
	if !conf.MostFrequentDisable {
		err := frequencyTable.LoadFromFile(frequencyFilePath())
		if err != nil {
			qerr.Error(err)
		}
	}
	// TODO: save the width of 2 columns

	favoritesWidget := favorites.NewFavoritesWidget(conf)
	app.favoritesWidget = favoritesWidget
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
	app.saveHistoryButton = saveHistoryButton
	miscLayout.AddWidget(saveHistoryButton, 0, 0)

	clearHistoryButton := widgets.NewQPushButton2("Clear History", nil)
	app.clearHistoryButton = clearHistoryButton
	miscLayout.AddWidget(clearHistoryButton, 0, 0)

	saveFavoritesButton := widgets.NewQPushButton2("Save Favorites", nil)
	app.saveFavoritesButton = saveFavoritesButton
	miscLayout.AddWidget(saveFavoritesButton, 0, 0)

	reloadDictsButton := widgets.NewQPushButton2("Reload Dicts", nil)
	app.reloadDictsButton = reloadDictsButton
	miscLayout.AddWidget(reloadDictsButton, 0, 0)

	closeDictsButton := widgets.NewQPushButton2("Close Dicts", nil)
	app.closeDictsButton = closeDictsButton
	miscLayout.AddWidget(closeDictsButton, 0, 0)

	reloadStyleButton := widgets.NewQPushButton2("Reload Style", nil)
	app.reloadStyleButton = reloadStyleButton
	miscLayout.AddWidget(reloadStyleButton, 0, 0)

	randomEntryButton := widgets.NewQPushButton2("Random Entry", nil)
	app.randomEntryButton = randomEntryButton
	miscLayout.AddWidget(randomEntryButton, 0, 0)

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
	app.dictsButton = dictsButton
	buttonBox.AddWidget(dictsButton, 0, core.Qt__AlignLeft)

	aboutButtonLabel := "About"
	if conf.ReduceMinimumWindowWidth {
		aboutButtonLabel = "\u200c"
	}
	aboutButton := newIconTextButton(aboutButtonLabel, widgets.QStyle__SP_MessageBoxInformation)
	buttonBox.AddWidget(aboutButton, 0, core.Qt__AlignLeft)

	buttonBox.AddStretch(1)

	openConfigButton := NewPNGIconTextButton("Config", "preferences-system-22.png")
	app.openConfigButton = openConfigButton
	buttonBox.AddWidget(openConfigButton, 0, 0)

	reloadConfigButton := newIconTextButton("Reload", widgets.QStyle__SP_BrowserReload)
	app.reloadConfigButton = reloadConfigButton
	buttonBox.AddWidget(reloadConfigButton, 0, 0)

	buttonBox.AddStretch(1)

	clearButton := widgets.NewQPushButton2("Clear", nil)
	app.clearButton = clearButton
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
	app.resultList = resultList
	leftPanelLayout.AddWidget(resultList, 0, 0)

	queryArgs := &QueryArgs{
		ArticleView: articleView,
		ResultList:  resultList,
		HeaderLabel: headerLabel,
		HistoryView: historyView,
		PostQuery:   app.postQuery,
		Entry:       entry,
		ModeCombo:   queryModeCombo,
	}
	app.queryArgs = queryArgs

	headerLabel.doQuery = app.doQuery
	articleView.doQuery = app.doQuery
	historyView.doQuery = app.doQuery

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

	app.allTextWidgets = []qutils.HasSetFont{
		queryLabel,
		entry,
		queryModeCombo,
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
		randomEntryButton,
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

	entry.ConnectReturnPressed(func() {
		onQuery(entry.Text(), queryArgs, false)
	})
	okButton.ConnectClicked(func(bool) {
		onQuery(entry.Text(), queryArgs, false)
	})
	aboutButton.ConnectClicked(func(bool) {
		aboutClicked(window)
	})

	for _, widget := range []KeyPressIface{
		resultList,
		articleView,
		historyView,
	} {
		app.setupKeyPressEvent(widget)
	}

	// --------------------------------------------------
	// setting up handlers
	app.setupHandlers()

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
		if !conf.SearchOnType {
			return
		}
		// log.Printf("event text=%#v, key=%x, modifiers=%x", event.Text(), event.Key(), event.Modifiers())
		if event.Key() >= 0x01000000 {
			return
		}
		if event.Modifiers() > core.Qt__ShiftModifier {
			return
		}
		text := entry.Text()
		if len(text) < conf.SearchOnTypeMinLength {
			return
		}
		onQuery(text, queryArgs, true)
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

	qs := qsettings.GetQSettings(window)

	queryModeCombo.ConnectCurrentIndexChanged(func(i int) {
		text := entry.Text()
		if text != "" {
			onQuery(text, queryArgs, false)
		}
		go qsettings.SaveSearchSettings(qs, queryModeCombo)
	})

	qsettings.RestoreSplitterSizes(qs, mainSplitter, QS_mainSplitter)
	qsettings.RestoreMainWinGeometry(app.QApplication, qs, window)
	qsettings.SetupMainWinGeometrySave(qs, window)

	qsettings.RestoreTableColumnsWidth(
		qs,
		frequencyTable.QTableWidget,
		QS_frequencyTable,
	)
	// frequencyTable.ConnectColumnResized does not work
	frequencyTable.HorizontalHeader().ConnectSectionResized(func(logicalIndex int, oldSize int, newSize int) {
		qsettings.SaveTableColumnsWidth(qs, frequencyTable.QTableWidget, QS_frequencyTable)
	})

	qsettings.SetupSplitterSizesSave(qs, mainSplitter, QS_mainSplitter)

	qsettings.RestoreSearchSettings(qs, queryModeCombo)

	window.Show()
	app.Exec()
}

func (app *Application) setupHandlers() {
	app.articleView.SetupCustomHandlers()
	app.historyView.SetupCustomHandlers()

	entry := app.entry
	queryArgs := app.queryArgs

	frequencyTable.ConnectItemActivated(func(item *widgets.QTableWidgetItem) {
		key := frequencyTable.Keys[item.Row()]
		app.doQuery(key)
		newRow := frequencyTable.KeyMap[key]
		// item.Column() panics!
		frequencyTable.SetCurrentCell(newRow, 0)
	})
	app.favoritesWidget.ConnectItemActivated(func(item *widgets.QListWidgetItem) {
		app.doQuery(item.Text())
	})

	app.reloadDictsButton.ConnectClicked(func(checked bool) {
		qdictmgr.InitDicts(conf, true)
		app.dictManager = nil
		onQuery(entry.Text(), queryArgs, false)
	})
	app.closeDictsButton.ConnectClicked(func(checked bool) {
		dictmgr.CloseDicts()
	})
	app.openConfigButton.ConnectClicked(func(checked bool) {
		OpenConfig()
	})
	app.reloadConfigButton.ConnectClicked(func(checked bool) {
		app.ReloadConfig()
		onQuery(entry.Text(), queryArgs, false)
	})
	app.reloadStyleButton.ConnectClicked(func(checked bool) {
		app.LoadUserStyle()
		onQuery(entry.Text(), queryArgs, false)
	})
	app.saveHistoryButton.ConnectClicked(func(checked bool) {
		SaveHistory()
		frequencyTable.SaveNoError()
	})
	app.randomEntryButton.ConnectClicked(func(checked bool) {
		res := dictmgr.RandomEntry(conf, resultFlags)
		if res == nil {
			return
		}
		query := res.F_Terms[0]
		entry.SetText(query)
		queryArgs.ResultList.SetResults([]common.SearchResultIface{res})
		queryArgs.AddHistoryAndFrequency(query)
		app.postQuery(query)
	})
	app.clearHistoryButton.ConnectClicked(func(checked bool) {
		app.historyView.ClearHistory()
		frequencyTable.Clear()
		frequencyTable.SaveNoError()
	})
	app.saveFavoritesButton.ConnectClicked(func(checked bool) {
		err := app.favoritesWidget.Save()
		if err != nil {
			qerr.Error(err)
		}
	})
	app.clearButton.ConnectClicked(func(checked bool) {
		app.resetQuery()
	})
	app.dictsButton.ConnectClicked(func(checked bool) {
		if app.runDictManager() {
			onQuery(entry.Text(), queryArgs, false)
		}
	})
}
