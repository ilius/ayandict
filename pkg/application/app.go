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

func (app *Application) activityComboChanged(index int) {
	switch index {
	case 0:
		app.historyView.Show()
		frequencyTable.Hide()
		app.favoritesWidget.Hide()
	case 1:
		app.historyView.Hide()
		frequencyTable.Show()
		app.favoritesWidget.Hide()
	case 2:
		app.historyView.Hide()
		frequencyTable.Hide()
		app.favoritesWidget.Show()
	}
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

	app.entry = widgets.NewQLineEdit(nil)
	app.entry.SetPlaceholderText("Type search query and press Enter")
	// to reduce inner margins:
	app.entry.SetTextMargins(0, -3, 0, -3)

	app.queryModeCombo = widgets.NewQComboBox(nil)
	app.queryModeCombo.AddItems([]string{
		"Fuzzy",
		"Start with",
		"Regex",
		"Glob",
	})

	okButton := widgets.NewQPushButton2(" OK ", nil)

	app.queryFavoriteButton = NewPNGIconTextButton("", "favorite.png")
	app.queryFavoriteButton.SetCheckable(true)
	app.queryFavoriteButton.SetToolTip("Add this query to favorites")

	// favoriteButtonVBox := widgets.NewQVBoxLayout()
	app.favoriteButton = NewPNGIconTextButton("", "favorite.png")
	app.favoriteButton.SetCheckable(true)
	app.favoriteButton.SetToolTip("Add this term to favorites")
	app.favoriteButton.Hide()
	// favoriteButtonVBox.AddWidget(favoriteButton, 0, core.Qt__AlignBottom)

	okButton.ConnectResizeEvent(func(event *gui.QResizeEvent) {
		h := event.Size().Height()
		if h > 100 {
			return
		}
		app.queryFavoriteButton.SetFixedSize2(h, h)
		app.favoriteButton.SetFixedSize2(h, h)
	})

	queryLabel := widgets.NewQLabel2("Query:", nil, 0)
	queryBox := widgets.NewQFrame(nil, 0)
	queryBoxLayout := widgets.NewQHBoxLayout2(queryBox)
	queryBoxLayout.SetContentsMargins(basePxHalf, basePxHalf, basePxHalf, 0)
	queryBoxLayout.SetSpacing(basePxI)
	queryBoxLayout.AddWidget(queryLabel, 0, 0)
	queryBoxLayout.AddWidget(app.entry, 0, 0)
	queryBoxLayout.AddWidget(app.queryModeCombo, 0, 0)
	queryBoxLayout.AddWidget(app.queryFavoriteButton, 0, 0)
	queryBoxLayout.AddWidget(okButton, 0, 0)

	app.headerLabel = CreateHeaderLabel(app)
	app.headerLabel.SetAlignment(core.Qt__AlignLeft)

	headerBox := widgets.NewQWidget(nil, 0)
	headerBox.SetSizePolicy2(widgets.QSizePolicy__Preferred, widgets.QSizePolicy__Minimum)
	headerBoxLayout := widgets.NewQHBoxLayout2(headerBox)
	// headerBoxLayout.SetSizeConstraint(widgets.QLayout__SetMinimumSize)
	headerBoxLayout.SetContentsMargins(0, 0, 0, 0)
	headerBoxLayout.AddSpacing(basePxHalf)
	headerBoxLayout.AddWidget(app.headerLabel, 1, 0)
	// headerBoxLayout.AddLayout(favoriteButtonVBox, 0)
	headerBoxLayout.AddWidget(app.favoriteButton, 0, core.Qt__AlignRight)
	headerBoxLayout.AddSpacing(int(basePx * 1.5))
	headerBox.SetSizePolicy2(expanding, widgets.QSizePolicy__Minimum)

	app.articleView = NewArticleView(app)

	app.historyView = NewHistoryView()

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

	app.favoritesWidget = favorites.NewFavoritesWidget(conf)
	{
		err := app.favoritesWidget.Load()
		if err != nil {
			// conf.FavoritesAutoSave = false
			fmt.Println(err)
		}
	}

	miscBox := widgets.NewQFrame(nil, 0)
	miscLayout := widgets.NewQVBoxLayout2(miscBox)
	miscLayout.SetContentsMargins(0, 0, 0, 0)

	app.saveHistoryButton = widgets.NewQPushButton2("Save History", nil)
	miscLayout.AddWidget(app.saveHistoryButton, 0, 0)

	app.clearHistoryButton = widgets.NewQPushButton2("Clear History", nil)
	miscLayout.AddWidget(app.clearHistoryButton, 0, 0)

	app.saveFavoritesButton = widgets.NewQPushButton2("Save Favorites", nil)
	miscLayout.AddWidget(app.saveFavoritesButton, 0, 0)

	app.reloadDictsButton = widgets.NewQPushButton2("Reload Dicts", nil)
	miscLayout.AddWidget(app.reloadDictsButton, 0, 0)

	app.closeDictsButton = widgets.NewQPushButton2("Close Dicts", nil)
	miscLayout.AddWidget(app.closeDictsButton, 0, 0)

	app.reloadStyleButton = widgets.NewQPushButton2("Reload Style", nil)
	miscLayout.AddWidget(app.reloadStyleButton, 0, 0)

	app.randomEntryButton = widgets.NewQPushButton2("Random Entry", nil)
	miscLayout.AddWidget(app.randomEntryButton, 0, 0)

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
	app.dictsButton = newIconTextButton(dictsButtonLabel, widgets.QStyle__SP_FileDialogDetailedView)
	buttonBox.AddWidget(app.dictsButton, 0, core.Qt__AlignLeft)

	aboutButtonLabel := "About"
	if conf.ReduceMinimumWindowWidth {
		aboutButtonLabel = "\u200c"
	}
	aboutButton := newIconTextButton(aboutButtonLabel, widgets.QStyle__SP_MessageBoxInformation)
	buttonBox.AddWidget(aboutButton, 0, core.Qt__AlignLeft)

	buttonBox.AddStretch(1)

	app.openConfigButton = NewPNGIconTextButton("Config", "preferences-system-22.png")
	buttonBox.AddWidget(app.openConfigButton, 0, 0)

	app.reloadConfigButton = newIconTextButton("Reload", widgets.QStyle__SP_BrowserReload)
	buttonBox.AddWidget(app.reloadConfigButton, 0, 0)

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
	leftMainLayout.AddWidget(app.articleView, 0, 0)
	leftMainLayout.AddSpacing(basePxHalf)
	leftMainLayout.AddLayout(buttonBox, 0)

	activityTypeCombo := widgets.NewQComboBox(nil)
	activityTypeCombo.AddItems([]string{
		"Recent",
		"Most Frequent",
		"Favorites",
	})

	frequencyTable.Hide()
	app.favoritesWidget.Hide()

	activityWidget := widgets.NewQWidget(nil, 0)
	activityLayout := widgets.NewQVBoxLayout2(activityWidget)
	activityLayout.SetContentsMargins(5, 5, 5, 5)
	activityLayout.AddWidget(activityTypeCombo, 0, 0)
	activityLayout.AddWidget(app.historyView, 0, 0)
	activityLayout.AddWidget(frequencyTable, 0, 0)
	activityLayout.AddWidget(app.favoritesWidget, 0, 0)

	activityTypeCombo.ConnectCurrentIndexChanged(app.activityComboChanged)

	onResultDisplay := func(terms []string) {
		app.favoriteButton.Show()
		app.favoriteButton.SetChecked(app.favoritesWidget.HasFavorite(terms[0]))
	}

	leftPanel := widgets.NewQWidget(nil, 0)
	leftPanelLayout := widgets.NewQVBoxLayout2(leftPanel)
	leftPanelLayout.AddWidget(widgets.NewQLabel2("Results", nil, 0), 0, 0)
	app.resultList = NewResultListWidget(
		app.articleView,
		app.headerLabel,
		onResultDisplay,
	)
	leftPanelLayout.AddWidget(app.resultList, 0, 0)

	app.queryArgs = &QueryArgs{
		ArticleView: app.articleView,
		ResultList:  app.resultList,
		HeaderLabel: app.headerLabel,
		HistoryView: app.historyView,
		PostQuery:   app.postQuery,
		Entry:       app.entry,
		ModeCombo:   app.queryModeCombo,
	}

	app.headerLabel.doQuery = app.doQuery
	app.articleView.doQuery = app.doQuery
	app.historyView.doQuery = app.doQuery

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
		app.entry,
		app.queryModeCombo,
		okButton,
		app.headerLabel,
		app.articleView,
		app.historyView,
		frequencyTable,
		app.favoritesWidget,
		app.saveHistoryButton,
		app.clearHistoryButton,
		app.saveFavoritesButton,
		app.reloadDictsButton,
		app.closeDictsButton,
		app.reloadStyleButton,
		app.randomEntryButton,
		app.dictsButton,
		aboutButton,
		app.openConfigButton,
		app.reloadConfigButton,
		clearButton,
		activityTypeCombo,
		app.resultList,
		rightPanel,
	}
	app.ReloadFont()

	okButton.ConnectClicked(func(bool) {
		onQuery(app.entry.Text(), app.queryArgs, false)
	})
	aboutButton.ConnectClicked(func(bool) {
		aboutClicked(window)
	})

	for _, widget := range []KeyPressIface{
		app.resultList,
		app.articleView,
		app.historyView,
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
			app.historyView.AddHistoryList(history)
		}
	}

	qs := qsettings.GetQSettings(window)
	app.setupSettings(qs, mainSplitter)

	window.Show()
	app.Exec()
}

func (app *Application) setupSettings(qs *core.QSettings, mainSplitter *widgets.QSplitter) {
	app.queryModeCombo.ConnectCurrentIndexChanged(func(i int) {
		text := app.entry.Text()
		if text != "" {
			onQuery(text, app.queryArgs, false)
		}
		go qsettings.SaveSearchSettings(qs, app.queryModeCombo)
	})

	qsettings.RestoreSplitterSizes(qs, mainSplitter, QS_mainSplitter)
	qsettings.RestoreMainWinGeometry(app.QApplication, qs, app.window)
	qsettings.SetupMainWinGeometrySave(qs, app.window)

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

	qsettings.RestoreSearchSettings(qs, app.queryModeCombo)
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
		onQuery(text, app.queryArgs, true)
	})

	app.favoriteButton.ConnectClicked(func(checked bool) {
		if app.resultList.Active == nil {
			app.favoriteButton.SetChecked(false)
			return
		}
		term := app.resultList.Active.Terms()[0]
		app.favoritesWidget.SetFavorite(term, checked)
		if term == entry.Text() {
			app.queryFavoriteButton.SetChecked(checked)
		}
	})
	app.queryFavoriteButton.ConnectClicked(func(checked bool) {
		term := entry.Text()
		if term == "" {
			app.queryFavoriteButton.SetChecked(false)
			return
		}
		app.favoritesWidget.SetFavorite(term, checked)
		if app.resultList.Active != nil && term == app.resultList.Active.Terms()[0] {
			app.favoriteButton.SetChecked(checked)
		}
	})
	app.entry.ConnectReturnPressed(func() {
		onQuery(entry.Text(), app.queryArgs, false)
	})
}
