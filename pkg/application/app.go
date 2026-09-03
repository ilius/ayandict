package application

import (
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"sync/atomic"
	"time"

	"github.com/ilius/ayandict/v3/pkg/about"
	"github.com/ilius/ayandict/v3/pkg/activity"
	"github.com/ilius/ayandict/v3/pkg/appinfo"
	"github.com/ilius/ayandict/v3/pkg/articleview"
	"github.com/ilius/ayandict/v3/pkg/audiocache"
	"github.com/ilius/ayandict/v3/pkg/config"
	"github.com/ilius/ayandict/v3/pkg/dictmgr"
	"github.com/ilius/ayandict/v3/pkg/dictmgr/qdictmgr"
	"github.com/ilius/ayandict/v3/pkg/favoritebutton"
	"github.com/ilius/ayandict/v3/pkg/frequencytable"
	"github.com/ilius/ayandict/v3/pkg/headerlabel"
	"github.com/ilius/ayandict/v3/pkg/logging"
	"github.com/ilius/ayandict/v3/pkg/qfavorites"
	"github.com/ilius/ayandict/v3/pkg/qlocalserver"
	"github.com/ilius/ayandict/v3/pkg/qtcommon"
	"github.com/ilius/ayandict/v3/pkg/qtcommon/qsettings"
	"github.com/ilius/ayandict/v3/pkg/qtutils"
	"github.com/ilius/ayandict/v3/pkg/resultlist"
	"github.com/ilius/ayandict/v3/pkg/resulttree"
	"github.com/ilius/ayandict/v3/pkg/utils"
	"github.com/ilius/ayandict/v3/pkg/webclient"
	"github.com/ilius/ayandict/v3/pkg/webserver"
	qt "github.com/mappu/miqt/qt6"
	"github.com/mappu/miqt/qt6/network"
)

const s_Soundex = "Soundex"

var (
	tipCtrl     = "Ctrl"
	tipCtrlPlus = "Ctrl﹢"
)

func init() {
	if runtime.GOOS == "darwin" {
		tipCtrl = "Command"
		tipCtrlPlus = "⌘ "
	}
}

var searchModes = []string{
	"Fuzzy",
	"Start with",
	"Regex",
	"Glob",
	"Word Match",
	// s_Soundex, added if config is set
}

var searchModeByDesc = map[string]dictmgr.SearchMode{
	"Fuzzy":      dictmgr.SearchModeFuzzy,
	"Start with": dictmgr.SearchModeStartWith,
	"Regex":      dictmgr.SearchModeRegex,
	"Glob":       dictmgr.SearchModeGlob,
	"Word Match": dictmgr.SearchModeWordMatch,
	s_Soundex:    dictmgr.SearchModeSoundex,
}

var exitSignalChan = make(chan os.Signal, 1)

type Application struct {
	*qt.QApplication

	audioCache *audiocache.AudioCache

	window *qt.QMainWindow
	icon   *qt.QIcon

	style *qt.QStyle

	mainWindowSettingsChan chan time.Time

	bottomBoxStyleOpt *qt.QStyleOptionButton

	dictManager *qdictmgr.DictManager

	queryArgs       *QueryArgs
	headerLabel     *headerlabel.HeaderLabel
	articleView     *articleview.ArticleView
	resultList      ResultsIface
	historyView     *HistoryView
	entry           *qt.QLineEdit
	searchModeCombo *qt.QComboBox
	favoritesWidget *qfavorites.FavoritesWidget
	frequencyTable  *frequencytable.FrequencyTable

	favoriteButton       FavoriteButtonInterface
	queryFavoriteButton  FavoriteButtonInterface
	openConfigButton     *qt.QPushButton
	reloadButton         *qt.QPushButton
	saveHistoryButton    *qt.QPushButton
	randomEntryButton    *qt.QPushButton
	randomFavoriteButton *qt.QPushButton
	clearHistoryButton   *qt.QPushButton
	clearButton          *qt.QPushButton
	dictsButton          *qt.QPushButton
	activityTypeCombo    *qt.QComboBox

	statusIconActions []*qt.QAction

	desktopWidget *DesktopWidget

	trayIcon          *qt.QSystemTrayIcon
	trayScanSelection *qt.QAction
	trayScanClipboard *qt.QAction
	trayScanAPI       *qt.QAction

	scanPopupCount atomic.Int32
}

func (app *Application) init() {
	if !LoadConfig() {
		conf = config.Default()
	}
	if len(conf.LocalServerPorts) == 0 {
		panic("config local_server_ports is empty")
	}
	webclient.Init(conf)
	app.mainWindowSettingsChan = make(chan time.Time, 100)
	app.audioCache = audiocache.NewAudioCache(conf)
}

func (app *Application) IsPopup() bool {
	return false
}

func (app *Application) Query(query string) {
	app.onQuery(query)
	app.entry.SetText(query)
}

func (app *Application) runDictManager() bool {
	if app.dictManager == nil {
		app.dictManager = qdictmgr.NewDictManager(app.QApplication, app.window.QWidget, conf)
	}
	return app.dictManager.Run()
}

func (app *Application) resetQuery() {
	app.entry.SetText("")
	app.queryArgs.ResultsLabel.SetText("Results")
	app.resultList.Clear()
	app.headerLabel.SetText("")
	app.articleView.SetHtml("")
	app.favoriteButton.Hide()
	app.queryFavoriteButton.SetChecked(false)
	app.historyView.ClearCursor()
}

func (app *Application) postQuery(query string) {
	if query == "" {
		app.queryFavoriteButton.SetChecked(false)
		return
	}
	app.queryFavoriteButton.SetChecked(app.favoritesWidget.HasFavorite(query))
}

func (app *Application) OnResultDisplay(terms []string) {
	app.favoriteButton.Show()
	app.favoriteButton.SetChecked(app.favoritesWidget.HasFavorite(terms[0]))
	app.favoriteButton.SetTerms(terms)
}

func (app *Application) onExit() {
	network.QLocalServer_RemoveServer(appinfo.LOCAL_SOCKET_NAME)
}

func (app *Application) Exit() {
	app.onExit()
	os.Exit(0)
}

// TODO: break down
func (app *Application) Run(query string) {
	slog.Info("Run", "query", query, "WebEnable", conf.WebEnable)
	app.init()

	logging.SetupLoggerAfterConfigLoad(
		os.Getenv("NO_COLOR") != "",
		conf,
	)
	slog.Debug(
		"Paths",
		"configDir", config.GetConfigDir(),
		"cacheDir", config.Paths.CacheDir(),
		"stateDir", config.Paths.StateDir(),
	)
	startSocketServer := func() bool {
		return qlocalserver.StartLocalSocketServer(
			conf,
			app.ShowWindowAndQuery,
			app.QueryPopup,
			app.statusIconActivate,
		)
	}
	if !startSocketServer() {
		if qlocalserver.PingLocalServer() {
			_ = qlocalserver.SendQueryToLocalServer(query)
			return
		}
		// No server responded to ping, so the socket file is stale
		// (left behind by a crashed previous process). Remove it and retry.
		slog.Warn(
			"socket exists but no server is running, removing stale socket (in /tmp)",
			"filename",
			appinfo.LOCAL_SOCKET_NAME,
		)
		network.QLocalServer_RemoveServer(appinfo.LOCAL_SOCKET_NAME)
		if !startSocketServer() {
			slog.Error(
				"another instance is running, or dead socket (in /tmp)",
				"filename",
				appinfo.LOCAL_SOCKET_NAME,
			)
			return
		}
	}
	if conf.WebEnable {
		if ok, _ := webclient.FindLocalWebServer(conf.LocalServerPorts); ok {
			slog.Error("another web instance is running")
			if query != "" {
				qlocalserver.SendQueryToLocalServer(query)
			}
			return
		}
		go webserver.StartServer(conf, conf.LocalServerPorts[0])
	}
	defer func() {
		if r := recover(); r != nil {
			app.onExit()
			panic(r)
		}
	}()
	app.OnAboutToQuit(app.onExit)
	{
		go func() {
			<-exitSignalChan
			app.Exit()
		}()
	}

	app.LoadUserStyle()
	qdictmgr.InitDicts(conf, true)

	basePx := app.baseFontPixelSize()

	basePxI := int(basePx)
	basePxHalf := int(basePx / 2)

	activityStorage := activity.NewActivityStorage(conf, config.GetConfigDir())

	frequencyTable := frequencytable.NewFrequencyView(
		activityStorage,
		conf.MostFrequentMaxSize,
	)
	app.frequencyTable = frequencyTable

	// icon := qt.NewQIcon5("./img/icon.png")

	window := app.window
	if config.PrivateMode {
		window.SetWindowTitle(appinfo.APP_DESC + " (private mode)")
	} else {
		window.SetWindowTitle(appinfo.APP_DESC)
	}
	window.OnCloseEvent(func(super func(*qt.QCloseEvent), event *qt.QCloseEvent) {
		if app.trayIcon != nil && app.trayIcon.IsVisible() {
			event.Ignore()
			window.Hide()
			return
		}
		super(event)
	})

	trayAvailable := qt.QSystemTrayIcon_IsSystemTrayAvailable()
	if !trayAvailable {
		if !conf.DesktopWidget {
			slog.Warn("system tray is not available, enabling desktop widget")
			conf.DesktopWidget = true
		}
	}
	{
		icon, err := loadPNGIcon(utils.IconPixName)
		if err != nil {
			slog.Error("failed to load window icon", "err", err)
			panic(err)
		}
		app.icon = icon
		window.SetWindowIcon(icon)
		if trayAvailable {
			app.setupTrayIcon(icon)
		}
	}

	qtutils.SetWinSize(window.QWidget, 600, 400)

	entry := qt.NewQLineEdit2()
	app.entry = entry
	entry.SetPlaceholderText("Type search query and press Enter")
	entry.SetTextMargins(0, -3, 0, -3) // to reduce inner margins

	searchModeCombo := qt.NewQComboBox2()
	app.searchModeCombo = searchModeCombo
	app.searchModeCombo.AddItems(searchModes)
	if conf.SoundexWordsFile != "" {
		app.searchModeCombo.AddItem(s_Soundex)
	}

	okButton := qt.NewQPushButton3(" OK ")

	app.queryFavoriteButton = favoritebutton.NewImageFavoriteButton(
		conf,
		app.queryFavoriteButtonClicked,
		app,
	)
	app.queryFavoriteButton.SetToolTips(
		"Add this query to favorites",
		"Remove this query from favorites",
	)

	// favoriteButtonVBox := qt.NewQVBoxLayout()
	app.favoriteButton = favoritebutton.NewImageFavoriteButton(
		conf,
		app.favoriteButtonClicked,
		app,
	)

	app.favoriteButton.SetToolTips(
		"Add this term to favorites (F)\nRight-click for multiple terms",
		"Remove this term from favorites (F)\nRight-click for multiple terms",
	)
	app.favoriteButton.Hide()
	// favoriteButtonVBox.AddWidget(favoriteButton, 0, qt.AlignBottom)

	okButton.OnResizeEvent(app.okButtonResized)

	queryLabel := qt.NewQLabel3("Query:")
	queryBox := qt.NewQFrame(nil)
	queryBoxLayout := qt.NewQHBoxLayout(queryBox.QWidget)
	queryBoxLayout.SetContentsMargins(basePxHalf, basePxHalf, basePxHalf, 0)
	queryBoxLayout.SetSpacing(basePxI)
	queryBoxLayout.AddWidget(queryLabel.QWidget)
	queryBoxLayout.AddWidget(entry.QWidget)
	queryBoxLayout.AddWidget(searchModeCombo.QWidget)
	queryBoxLayout.AddWidget(app.queryFavoriteButton.QWidget())
	queryBoxLayout.AddWidget(okButton.QWidget)

	headerLabel := headerlabel.NewHeaderLabel(conf, app, headerTpl)
	app.headerLabel = headerLabel
	app.headerLabel.SetAlignment(qt.AlignLeft)

	headerBox := qt.NewQWidget(nil)
	headerBox.SetSizePolicy2(qt.QSizePolicy__Preferred, qt.QSizePolicy__Minimum)
	headerBoxLayout := qt.NewQHBoxLayout(headerBox)
	// headerBoxLayout.SetSizeConstraint(qt.QLayout__SetMinimumSize)
	headerBoxLayout.SetContentsMargins(0, 0, 0, 0)
	headerBoxLayout.AddSpacing(basePxHalf)
	headerBoxLayout.AddWidget3(headerLabel.QWidget, 1, 0)
	// headerBoxLayout.AddLayout(favoriteButtonVBox, 0)
	headerBoxLayout.AddWidget3(app.favoriteButton.QWidget(), 0, qt.AlignRight)
	headerBoxLayout.AddSpacing(int(basePx * 1.5))
	headerBox.SetSizePolicy2(expanding, qt.QSizePolicy__Minimum)

	articleView := articleview.NewArticleView(conf, app)
	app.articleView = articleView

	historyView := NewHistoryView(
		activityStorage,
		conf.HistoryMaxSize,
		app.Query,
	)
	app.historyView = historyView
	if !conf.HistoryDisable {
		err := historyView.Load()
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
	app.favoritesWidget.SetFocusPolicy(qt.NoFocus)
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
	app.clearHistoryButton.SetToolTip(fmt.Sprintf("Shortcut: %sDelete", tipCtrlPlus))

	app.randomEntryButton = qt.NewQPushButton3("Random Entry")
	miscLayout.AddWidget(app.randomEntryButton.QWidget)

	app.randomFavoriteButton = qt.NewQPushButton3("Random Favorite")
	miscLayout.AddWidget(app.randomFavoriteButton.QWidget)

	app.updateMiscButtonsPadding()

	buttonBox := qt.NewQHBoxLayout2()
	buttonBox.SetContentsMargins(0, 0, 0, 0)
	buttonBox.SetSpacing(basePxHalf)

	dictsButtonLabel := "Dictionaries"
	if conf.ReduceMinimumWindowWidth {
		dictsButtonLabel = "Dicts"
	}
	app.dictsButton = app.newIconTextButton(dictsButtonLabel, qt.QStyle__SP_FileDialogDetailedView)
	buttonBox.AddWidget3(app.dictsButton.QWidget, 0, qt.AlignLeft)
	app.dictsButton.SetToolTip(fmt.Sprintf("Manage Dictionaries (%sD)", tipCtrlPlus))
	app.dictsButton.SetFocusPolicy(qt.NoFocus)

	aboutButton := app.makeAboutButton(conf)
	buttonBox.AddWidget3(aboutButton.QWidget, 0, qt.AlignLeft)
	aboutButton.SetToolTip("Show About window")
	aboutButton.SetFocusPolicy(qt.NoFocus)

	buttonBox.AddStretch()

	app.openConfigButton = NewPNGIconTextButton("Config", "preferences-system-22.png")
	buttonBox.AddWidget3(app.openConfigButton.QWidget, 0, 0)
	app.openConfigButton.SetToolTip("Open config file in your default editor")
	app.openConfigButton.SetFocusPolicy(qt.NoFocus)

	app.reloadButton = app.newIconTextButton("Reload", qt.QStyle__SP_BrowserReload)
	buttonBox.AddWidget3(app.reloadButton.QWidget, 0, 0)
	app.reloadButton.SetToolTip(
		"Reload config file" +
			"\nHold " + tipCtrl + " to also reload dictionaries and style" +
			"\nHold Shift to reload config and style",
	)
	app.reloadButton.SetFocusPolicy(qt.NoFocus)

	buttonBox.AddStretch()

	app.clearButton = qt.NewQPushButton3("Clear")
	app.clearButton.SetToolTip(fmt.Sprintf(
		"Clear query and results\nHold %s to also clear history",
		tipCtrl,
	))
	buttonBox.AddWidget3(app.clearButton.QWidget, 0, qt.AlignRight)
	app.clearButton.SetFocusPolicy(qt.NoFocus)

	leftMainWidget := qt.NewQWidget(nil)
	leftMainLayout := qt.NewQVBoxLayout(leftMainWidget)
	leftMainLayout.SetContentsMargins(0, 0, 0, 0)
	leftMainLayout.SetSpacing(0)
	leftMainLayout.AddWidget3(queryBox.QWidget, 0, 0)
	leftMainLayout.AddSpacing(basePxHalf)
	leftMainLayout.AddWidget3(headerBox, 0, 0)
	leftMainLayout.AddSpacing(basePxHalf)
	leftMainLayout.AddWidget3(app.articleView.Widget, 0, 0)
	leftMainLayout.AddSpacing(basePxHalf)
	leftMainLayout.AddLayout(buttonBox.Layout())

	activityTypeCombo := qt.NewQComboBox2()
	app.activityTypeCombo = activityTypeCombo
	activityTypeCombo.AddItems([]string{
		"History",
		"Most Frequent",
		"Favorites",
	})

	frequencyTable.Hide()
	app.favoritesWidget.Hide()

	activityWidget := qt.NewQWidget(nil)
	activityLayout := qt.NewQVBoxLayout(activityWidget)
	activityLayout.SetContentsMargins(5, 5, 5, 5)
	activityLayout.AddWidget(activityTypeCombo.QWidget)
	activityLayout.AddWidget(historyView.QWidget)
	activityLayout.AddWidget(frequencyTable.QWidget)
	activityLayout.AddWidget(app.favoritesWidget.QWidget)

	activityTypeCombo.OnCurrentIndexChanged(app.activityComboChanged)

	leftPanel := qt.NewQWidget(nil)
	leftPanelLayout := qt.NewQVBoxLayout(leftPanel)
	resultsLabel := qt.NewQLabel3("Results")
	leftPanelLayout.AddWidget(resultsLabel.QWidget)
	if conf.ResultTree {
		app.resultList = resulttree.NewResultTree(
			articleView,
			headerLabel,
			app,
		)
	} else {
		app.resultList = resultlist.NewResultList(
			articleView,
			headerLabel,
			app,
		)
	}
	leftPanelLayout.AddWidget(app.resultList.QWidget())

	app.queryArgs = &QueryArgs{
		ArticleView:    articleView,
		ResultsLabel:   resultsLabel,
		HeaderLabel:    headerLabel,
		HistoryView:    historyView,
		Entry:          entry,
		ModeCombo:      searchModeCombo,
		FrequencyTable: frequencyTable,
	}

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

	qt.QApplication_SetFont(qtcommon.ConfigFont(conf))

	app.reloadFont()

	okButton.OnClicked(func() {
		app.onQuery(entry.Text())
	})

	// --------------------------------------------------
	// setting up handlers
	app.setupHandlers()

	app.setupSettings(mainSplitter)

	app.setupScanPopup()

	if conf.DesktopWidget {
		app.setupDekstopWidget()
	}

	hidden := conf.StartHidden && app.trayIcon != nil && app.trayIcon.IsVisible()
	if query != "" {
		app.Query(query)
		hidden = false
	}
	if !hidden {
		window.Show()
	}

	if conf.RandomFavoritePopupIntervalSeconds > 0 {
		timer := qt.NewQTimer()
		timer.SetSingleShot(true)
		onClose := func() {
			timer.Start(1000 * conf.RandomFavoritePopupIntervalSeconds)
		}
		timer.OnTimeout(func() {
			app.randomFavoritePopupAuto(onClose)
		})
		onClose()
	}

	_ = qt.QApplication_Exec()
}

func (app *Application) AudioCache() *audiocache.AudioCache {
	return app.audioCache
}

func (app *Application) setupSettings(mainSplitter *qt.QSplitter) {
	app.searchModeCombo.OnCurrentIndexChanged(func(i int) {
		text := app.entry.Text()
		if text != "" {
			app.onQuery(text)
		}
		app.mainWindowSettingsChan <- time.Now()
	})

	qsettings.RestoreSplitterSizes(mainSplitter, QS_mainSplitter)
	app.restoreMainWindowSettings()
	app.setupMainWindowSettings()

	frequencyTable := app.frequencyTable
	qsettings.RestoreTableColumnsWidth(
		frequencyTable.QTableWidget,
		QS_frequencyTable,
	)
	// frequencyTable.OnColumnResized does not work
	frequencyTable.HorizontalHeader().OnSectionResized(func(logicalIndex int, oldSize int, newSize int) {
		qsettings.SaveTableColumnsWidth(frequencyTable.QTableWidget, QS_frequencyTable)
	})

	qsettings.SetupSplitterSizesSave(mainSplitter, QS_mainSplitter)
}

func (app *Application) updateMiscButtonsPadding() {
	vpadding := conf.MiscButtonsVerticalPadding
	stylesheet := fmt.Sprintf("padding-top: %dpx; padding-bottom: %dpx;", vpadding, vpadding)

	app.saveHistoryButton.SetStyleSheet(stylesheet)
	app.clearHistoryButton.SetStyleSheet(stylesheet)
	app.randomEntryButton.SetStyleSheet(stylesheet)
	app.randomFavoriteButton.SetStyleSheet(stylesheet)
}

func (app *Application) statusIconActivate() {
	window := app.window
	if window.IsActiveWindow() {
		window.Hide()
	} else {
		window.ShowNormal()
		window.ActivateWindow()
	}
}

func (app *Application) ShowAbout() {
	window := qt.NewQDialog(app.window.QWidget)
	about.ShowAbout(window.QWidget, app.icon)
	window.ShowNormal()
}
