package application

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/ilius/ayandict/v3/pkg/activity"
	"github.com/ilius/ayandict/v3/pkg/appinfo"
	"github.com/ilius/ayandict/v3/pkg/application/frequency"
	"github.com/ilius/ayandict/v3/pkg/application/qfavorites"
	"github.com/ilius/ayandict/v3/pkg/config"
	"github.com/ilius/ayandict/v3/pkg/dictmgr"
	"github.com/ilius/ayandict/v3/pkg/dictmgr/qdictmgr"
	"github.com/ilius/ayandict/v3/pkg/logging"
	"github.com/ilius/ayandict/v3/pkg/qtcommon"
	"github.com/ilius/ayandict/v3/pkg/qtcommon/qerr"
	"github.com/ilius/ayandict/v3/pkg/qtcommon/qsettings"
	"github.com/ilius/ayandict/v3/pkg/server"
	common "github.com/ilius/go-dict-commons"
	qt "github.com/mappu/miqt/qt6"
)

const (
	QS_mainSplitter   = "main_splitter"
	QS_frequencyTable = "frequencytable"

	escape               = int(qt.Key_Escape)
	shortcutModifierMask = int(qt.ControlModifier) | int(qt.AltModifier) | int(qt.MetaModifier)
)

// basePx is %66 of the font size in pixels,
// I'm using it for spacing between widgets
// kinda like "em" in html, but probably not exactly the same
var basePx = float32(10)

type Application struct {
	*qt.QApplication

	window *qt.QMainWindow

	style *qt.QStyle

	qs *qt.QSettings

	bottomBoxStyleOpt *qt.QStyleOptionButton

	dictManager *qdictmgr.DictManager

	allTextWidgets []qtcommon.HasSetFont

	queryArgs       *QueryArgs
	headerLabel     *HeaderLabel
	articleView     *ArticleView
	resultList      *ResultListWidget
	historyView     *HistoryView
	entry           *qt.QLineEdit
	queryModeCombo  *qt.QComboBox
	favoritesWidget *qfavorites.FavoritesWidget

	favoriteButton       *FavoriteButton
	queryFavoriteButton  *FavoriteButton
	reloadDictsButton    *qt.QPushButton
	closeDictsButton     *qt.QPushButton
	openConfigButton     *qt.QPushButton
	reloadConfigButton   *qt.QPushButton
	reloadStyleButton    *qt.QPushButton
	saveHistoryButton    *qt.QPushButton
	randomEntryButton    *qt.QPushButton
	randomFavoriteButton *qt.QPushButton
	clearHistoryButton   *qt.QPushButton
	saveFavoritesButton  *qt.QPushButton
	clearButton          *qt.QPushButton
	dictsButton          *qt.QPushButton
	activityTypeCombo    *qt.QComboBox
}

func Run() {
	app := &Application{
		QApplication:   qt.NewQApplication(os.Args),
		window:         qt.NewQMainWindow(nil),
		allTextWidgets: []qtcommon.HasSetFont{},
	}
	qerr.ShowMessage = showErrorMessage
	app.style = qt.QApplication_Style()
	app.bottomBoxStyleOpt = qt.NewQStyleOptionButton()
	qt.QCoreApplication_SetApplicationName(appinfo.APP_DESC)

	if cacheDir == "" {
		slog.Error("cacheDir is empty")
	}
	{
		err := os.MkdirAll(cacheDir, 0o755)
		if err != nil {
			slog.Error("error in MkdirAll: " + err.Error())
		}
	}

	app.Run()
}

func (app *Application) init() {
	if !LoadConfig() {
		conf = config.Default()
	}
	if len(conf.LocalServerPorts) == 0 {
		panic("config local_server_ports is empty")
	}
	client.Timeout = conf.LocalClientTimeout

	logging.SetupLoggerAfterConfigLoad(false, conf)

	if ok, _ := findLocalServer(conf.LocalServerPorts); ok {
		slog.Error("another instance is running")
		return
	}
	go server.StartServer(conf.LocalServerPorts[0])

	app.LoadUserStyle()
	qdictmgr.InitDicts(conf, true)
}

func (app *Application) newIconTextButton(label string, pix qt.QStyle__StandardPixmap) *qt.QPushButton {
	return qt.NewQPushButton4(
		app.style.StandardIcon(
			pix,
			app.bottomBoxStyleOpt.QStyleOption,
			nil,
		),
		label,
	)
}

func (app *Application) doQuery(query string) {
	onQuery(query, app.queryArgs, false)
	app.entry.SetText(query)
}

func (app *Application) runDictManager() bool {
	if app.dictManager == nil {
		app.dictManager = qdictmgr.NewDictManager(app.QApplication, app.window.QWidget, conf)
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
	widget.OnKeyPressEvent(func(super func(event *qt.QKeyEvent), event *qt.QKeyEvent) {
		switch event.Key() {
		case int(qt.Key_Space): // " "
			app.entry.SetFocusWithReason(qt.ShortcutFocusReason)
			return
		case int(qt.Key_Plus), int(qt.Key_Equal): // "+", "="
			app.articleView.ZoomIn()
			return
		case int(qt.Key_Minus): // "-"
			app.articleView.ZoomOut()
			return
		case escape: // event.Text()="\x1b"
			app.resetQuery()
			return
		case int(qt.Key_F1):
			aboutClicked(app.window.QWidget)
			return
		}
		super(event)
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
	qsettings.SaveActivityMode(app.qs, app.activityTypeCombo)
}

// TODO: break down
func (app *Application) Run() {
	app.init()

	basePx = float32(fontPixelSize(
		qt.QApplication_Font(),
		qt.QGuiApplication_PrimaryScreen().PhysicalDotsPerInch(),
	) * 0.66)

	basePxI := int(basePx)
	basePxHalf := int(basePx / 2)

	activityStorage := activity.NewActivityStorage(conf, config.GetConfigDir())

	frequencyTable = frequency.NewFrequencyView(
		activityStorage,
		conf.MostFrequentMaxSize,
	)

	// icon := qt.NewQIcon5("./img/icon.png")

	window := app.window
	window.SetWindowTitle(appinfo.APP_DESC)
	window.Resize(600, 400)

	app.entry = qt.NewQLineEdit2()
	app.entry.SetPlaceholderText("Type search query and press Enter")
	// to reduce inner margins:
	app.entry.SetTextMargins(0, -3, 0, -3)

	app.queryModeCombo = qt.NewQComboBox2()
	app.queryModeCombo.AddItems([]string{
		"Fuzzy",
		"Start with",
		"Regex",
		"Glob",
	})

	okButton := qt.NewQPushButton3(" OK ")

	app.queryFavoriteButton = NewFavoriteButton(func(checked bool) {
		term := app.entry.Text()
		if term == "" {
			app.queryFavoriteButton.SetChecked(false)
			return
		}
		app.favoritesWidget.SetFavorite(term, checked)
		if app.resultList.Active != nil && term == app.resultList.Active.Terms()[0] {
			app.favoriteButton.SetChecked(checked)
		}
	})
	app.queryFavoriteButton.SetToolTips(
		"Add this query to favorites",
		"Remove this query from favorites",
	)

	// favoriteButtonVBox := qt.NewQVBoxLayout()
	app.favoriteButton = NewFavoriteButton(func(checked bool) {
		if app.resultList.Active == nil {
			app.favoriteButton.SetChecked(false)
			return
		}
		term := app.resultList.Active.Terms()[0]
		app.favoritesWidget.SetFavorite(term, checked)
		if term == app.entry.Text() {
			app.queryFavoriteButton.SetChecked(checked)
		}
	})

	app.favoriteButton.SetToolTips(
		"Add this term to favorites",
		"Remove this term from favorites",
	)
	app.favoriteButton.Hide()
	// favoriteButtonVBox.AddWidget(favoriteButton, 0, qt.AlignBottom)

	okButton.OnResizeEvent(func(_ func(*qt.QResizeEvent), event *qt.QResizeEvent) {
		h := event.Size().Height()
		if h > 100 {
			return
		}
		app.queryFavoriteButton.SetFixedSize2(h, h)
		app.favoriteButton.SetFixedSize2(h, h)
	})

	queryLabel := qt.NewQLabel3("Query:")
	queryBox := qt.NewQFrame(nil)
	queryBoxLayout := qt.NewQHBoxLayout(queryBox.QWidget)
	queryBoxLayout.SetContentsMargins(basePxHalf, basePxHalf, basePxHalf, 0)
	queryBoxLayout.SetSpacing(basePxI)
	queryBoxLayout.AddWidget(queryLabel.QWidget)
	queryBoxLayout.AddWidget(app.entry.QWidget)
	queryBoxLayout.AddWidget(app.queryModeCombo.QWidget)
	queryBoxLayout.AddWidget(app.queryFavoriteButton.QWidget)
	queryBoxLayout.AddWidget(okButton.QWidget)

	app.headerLabel = CreateHeaderLabel(app)
	app.headerLabel.SetAlignment(qt.AlignLeft)

	headerBox := qt.NewQWidget(nil)
	headerBox.SetSizePolicy2(qt.QSizePolicy__Preferred, qt.QSizePolicy__Minimum)
	headerBoxLayout := qt.NewQHBoxLayout(headerBox)
	// headerBoxLayout.SetSizeConstraint(qt.QLayout__SetMinimumSize)
	headerBoxLayout.SetContentsMargins(0, 0, 0, 0)
	headerBoxLayout.AddSpacing(basePxHalf)
	headerBoxLayout.AddWidget3(app.headerLabel.QWidget, 1, 0)
	// headerBoxLayout.AddLayout(favoriteButtonVBox, 0)
	headerBoxLayout.AddWidget3(app.favoriteButton.QWidget, 0, qt.AlignRight)
	headerBoxLayout.AddSpacing(int(basePx * 1.5))
	headerBox.SetSizePolicy2(expanding, qt.QSizePolicy__Minimum)

	app.articleView = NewArticleView(app)

	app.historyView = NewHistoryView(activityStorage, conf.HistoryMaxSize)
	if !conf.HistoryDisable {
		err := app.historyView.Load()
		if err != nil {
			slog.Error("error in loading history: " + err.Error())
		}
	}

	{
		item := qt.NewQTableWidgetItem2("Query")
		item.SetTextAlignment(0)
		frequencyTable.SetHorizontalHeaderItem(0, item)
	}
	{
		item := qt.NewQTableWidgetItem2("#")
		item.SetTextAlignment(0)
		frequencyTable.SetHorizontalHeaderItem(1, item)
	}
	if !conf.MostFrequentDisable {
		err := frequencyTable.Load()
		if err != nil {
			slog.Error("error in loading frequency table: " + err.Error())
		}
	}
	// TODO: save and restore the width of 2 columns

	app.favoritesWidget = qfavorites.NewFavoritesWidget(conf)
	{
		err := app.favoritesWidget.Load()
		if err != nil {
			// conf.FavoritesAutoSave = false
			fmt.Println(err)
		}
	}

	miscBox := qt.NewQFrame(nil)
	miscLayout := qt.NewQVBoxLayout(miscBox.QWidget)
	miscLayout.SetContentsMargins(0, 0, 0, 0)

	app.saveHistoryButton = qt.NewQPushButton3("Save History")
	miscLayout.AddWidget(app.saveHistoryButton.QWidget)

	app.clearHistoryButton = qt.NewQPushButton3("Clear History")
	miscLayout.AddWidget(app.clearHistoryButton.QWidget)

	app.saveFavoritesButton = qt.NewQPushButton3("Save Favorites")
	miscLayout.AddWidget(app.saveFavoritesButton.QWidget)

	app.reloadDictsButton = qt.NewQPushButton3("Reload Dicts")
	miscLayout.AddWidget(app.reloadDictsButton.QWidget)

	app.closeDictsButton = qt.NewQPushButton3("Close Dicts")
	miscLayout.AddWidget(app.closeDictsButton.QWidget)

	app.reloadStyleButton = qt.NewQPushButton3("Reload Style")
	miscLayout.AddWidget(app.reloadStyleButton.QWidget)

	app.randomEntryButton = qt.NewQPushButton3("Random Entry")
	miscLayout.AddWidget(app.randomEntryButton.QWidget)

	app.randomFavoriteButton = qt.NewQPushButton3("Random Favorite")
	miscLayout.AddWidget(app.randomFavoriteButton.QWidget)

	buttonBox := qt.NewQHBoxLayout2()
	buttonBox.SetContentsMargins(0, 0, 0, 0)
	buttonBox.SetSpacing(basePxHalf)

	dictsButtonLabel := "Dictionaries"
	if conf.ReduceMinimumWindowWidth {
		dictsButtonLabel = "Dicts"
	}
	app.dictsButton = app.newIconTextButton(dictsButtonLabel, qt.QStyle__SP_FileDialogDetailedView)
	buttonBox.AddWidget3(app.dictsButton.QWidget, 0, qt.AlignLeft)

	aboutButton := app.makeAboutButton(conf)
	buttonBox.AddWidget3(aboutButton.QWidget, 0, qt.AlignLeft)

	buttonBox.AddStretch()

	app.openConfigButton = NewPNGIconTextButton("Config", "preferences-system-22.png")
	buttonBox.AddWidget3(app.openConfigButton.QWidget, 0, 0)

	app.reloadConfigButton = app.newIconTextButton("Reload", qt.QStyle__SP_BrowserReload)
	buttonBox.AddWidget3(app.reloadConfigButton.QWidget, 0, 0)

	buttonBox.AddStretch()

	app.clearButton = qt.NewQPushButton3("Clear")
	buttonBox.AddWidget3(app.clearButton.QWidget, 0, qt.AlignRight)

	leftMainWidget := qt.NewQWidget(nil)
	leftMainLayout := qt.NewQVBoxLayout(leftMainWidget)
	leftMainLayout.SetContentsMargins(0, 0, 0, 0)
	leftMainLayout.SetSpacing(0)
	leftMainLayout.AddWidget3(queryBox.QWidget, 0, 0)
	leftMainLayout.AddSpacing(basePxHalf)
	leftMainLayout.AddWidget3(headerBox, 0, 0)
	leftMainLayout.AddSpacing(basePxHalf)
	leftMainLayout.AddWidget3(app.articleView.QWidget, 0, 0)
	leftMainLayout.AddSpacing(basePxHalf)
	leftMainLayout.AddLayout(buttonBox.Layout())

	activityTypeCombo := qt.NewQComboBox2()
	app.activityTypeCombo = activityTypeCombo
	activityTypeCombo.AddItems([]string{
		"Recent",
		"Most Frequent",
		"Favorites",
	})

	frequencyTable.Hide()
	app.favoritesWidget.Hide()

	activityWidget := qt.NewQWidget(nil)
	activityLayout := qt.NewQVBoxLayout(activityWidget)
	activityLayout.SetContentsMargins(5, 5, 5, 5)
	activityLayout.AddWidget(activityTypeCombo.QWidget)
	activityLayout.AddWidget(app.historyView.QWidget)
	activityLayout.AddWidget(frequencyTable.QWidget)
	activityLayout.AddWidget(app.favoritesWidget.QWidget)

	activityTypeCombo.OnCurrentIndexChanged(app.activityComboChanged)

	onResultDisplay := func(terms []string) {
		app.favoriteButton.Show()
		app.favoriteButton.SetChecked(app.favoritesWidget.HasFavorite(terms[0]))
	}

	leftPanel := qt.NewQWidget(nil)
	leftPanelLayout := qt.NewQVBoxLayout(leftPanel)
	leftPanelLayout.AddWidget(qt.NewQLabel3("Results").QWidget)
	app.resultList = NewResultListWidget(
		app.articleView,
		app.headerLabel,
		onResultDisplay,
	)
	leftPanelLayout.AddWidget(app.resultList.QWidget)

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

	rightPanel := qt.NewQTabWidget(nil)
	_ = rightPanel.AddTab(activityWidget, " Activity ")
	_ = rightPanel.AddTab(miscBox.QWidget, " Misc ")

	mainSplitter := qt.NewQSplitter(nil)
	mainSplitter.SetSizePolicy2(expanding, expanding)
	mainSplitter.AddWidget(leftPanel)
	mainSplitter.AddWidget(leftMainWidget)
	mainSplitter.AddWidget(rightPanel.QWidget)
	mainSplitter.SetStretchFactor(0, 1)
	mainSplitter.SetStretchFactor(1, 5)
	mainSplitter.SetStretchFactor(2, 1)

	window.SetCentralWidget(mainSplitter.QWidget)

	qt.QApplication_SetFont(ConfigFont())

	app.allTextWidgets = []qtcommon.HasSetFont{
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
		app.randomFavoriteButton,
		app.dictsButton,
		aboutButton,
		app.openConfigButton,
		app.reloadConfigButton,
		app.clearButton,
		activityTypeCombo,
		app.resultList,
		rightPanel,
	}
	app.ReloadFont()

	okButton.OnClicked(func() {
		onQuery(app.entry.Text(), app.queryArgs, false)
	})

	for _, widget := range []KeyPressIface{
		app.window,
		app.resultList.QListWidget,
		app.articleView,
		app.historyView.QListWidget,
	} {
		app.setupKeyPressEvent(widget)
	}

	// --------------------------------------------------
	// setting up handlers
	app.setupHandlers()

	qs := qsettings.GetQSettings(window.QObject)
	app.qs = qs
	app.setupSettings(qs, mainSplitter)
	qsettings.RestoreActivityMode(qs, activityTypeCombo)

	window.Show()
	_ = qt.QApplication_Exec()
}

func (app *Application) makeAboutButton(conf *config.Config) *qt.QPushButton {
	aboutButtonLabel := "About"
	if conf.ReduceMinimumWindowWidth {
		aboutButtonLabel = "\u200c"
	}
	aboutButton := app.newIconTextButton(aboutButtonLabel, qt.QStyle__SP_MessageBoxInformation)
	aboutButton.OnClicked(func() {
		aboutClicked(app.window.QWidget)
	})
	return aboutButton
}

func (app *Application) setupSettings(qs *qt.QSettings, mainSplitter *qt.QSplitter) {
	app.queryModeCombo.OnCurrentIndexChanged(func(i int) {
		text := app.entry.Text()
		if text != "" {
			onQuery(text, app.queryArgs, false)
		}
		go qsettings.SaveSearchSettings(qs, app.queryModeCombo)
	})

	qsettings.RestoreSplitterSizes(qs, mainSplitter, QS_mainSplitter)
	qsettings.RestoreMainWinGeometry(qs, app.window)
	qsettings.SetupMainWinGeometrySave(qs, app.window)

	qsettings.RestoreTableColumnsWidth(
		qs,
		frequencyTable.QTableWidget,
		QS_frequencyTable,
	)
	// frequencyTable.OnColumnResized does not work
	frequencyTable.HorizontalHeader().OnSectionResized(func(logicalIndex int, oldSize int, newSize int) {
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

	frequencyTable.OnItemActivated(func(item *qt.QTableWidgetItem) {
		key := frequencyTable.Keys[item.Row()]
		app.doQuery(key)
		newRow := frequencyTable.KeyMap[key]
		// item.Column() panics!
		frequencyTable.SetCurrentCell(newRow, 0)
	})
	app.favoritesWidget.OnItemActivated(func(item *qt.QListWidgetItem) {
		app.doQuery(item.Text())
	})

	app.reloadDictsButton.OnClicked(func() {
		qdictmgr.InitDicts(conf, true)
		app.dictManager = nil
		onQuery(entry.Text(), queryArgs, false)
	})
	app.closeDictsButton.OnClicked(func() {
		dictmgr.CloseDicts()
	})
	app.openConfigButton.OnClicked(func() {
		OpenConfig()
	})
	app.reloadConfigButton.OnClicked(func() {
		app.ReloadConfig()
		onQuery(entry.Text(), queryArgs, false)
	})
	app.reloadStyleButton.OnClicked(func() {
		app.LoadUserStyle()
		onQuery(entry.Text(), queryArgs, false)
	})
	app.saveHistoryButton.OnClicked(func() {
		app.historyView.Save()
		frequencyTable.SaveNoError()
	})
	app.randomEntryButton.OnClicked(func() {
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
	app.randomFavoriteButton.OnClicked(func() {
		term := app.favoritesWidget.Data.Random()
		if term == "" {
			// show "No Favorites" error?
			return
		}
		onQuery(term, queryArgs, false)
	})
	app.clearHistoryButton.OnClicked(func() {
		app.historyView.ClearHistory()
		frequencyTable.Clear()
		frequencyTable.SaveNoError()
	})
	app.saveFavoritesButton.OnClicked(func() {
		err := app.favoritesWidget.Save()
		if err != nil {
			slog.Error("error saving favorites: " + err.Error())
		}
	})
	app.clearButton.OnClicked(func() {
		app.resetQuery()
	})
	app.dictsButton.OnClicked(func() {
		if app.runDictManager() {
			onQuery(entry.Text(), queryArgs, false)
		}
	})
	entry.OnKeyPressEvent(func(super func(*qt.QKeyEvent), event *qt.QKeyEvent) {
		// slog.Info(
		// 	"entry: KeyPressEvent",
		// 	"text", fmt.Sprintf("%#v", event.Text()),
		// 	"key", event.Key(),
		// )
		key := event.Key()
		switch key {
		case escape: // event.Text()="\x1b"
			app.window.SetFocus()
			return
		case int(qt.Key_Return), int(qt.Key_Enter): // event.Text()="\r"
			onQuery(entry.Text(), app.queryArgs, false)
			return
		}

		super(event)

		// event.Modifiers(): qt.NoModifier, qt.ShiftModifier, KeypadModifier
		if conf.SearchOnType && key < escape {
			if int(event.Modifiers())&shortcutModifierMask == 0 {
				text := entry.Text()
				// slog.Debug("checking SearchOnType") // FIXME: panics
				if len(text) >= conf.SearchOnTypeMinLength {
					onQuery(text, app.queryArgs, true)
				}
				return
			}
		}
	})

	if config.PrivateMode {
		app.favoriteButton.SetDisabled(true)
		app.queryFavoriteButton.SetDisabled(true)
	}
	// slog.Error("test error", "s", "hello", "n", 2, "b", true)
}
