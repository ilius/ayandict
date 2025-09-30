package application

import (
	"fmt"
	"log/slog"

	"github.com/ilius/ayandict/v3/pkg/activity"
	"github.com/ilius/ayandict/v3/pkg/appinfo"
	"github.com/ilius/ayandict/v3/pkg/application/frequency"
	"github.com/ilius/ayandict/v3/pkg/application/qfavorites"
	"github.com/ilius/ayandict/v3/pkg/config"
	"github.com/ilius/ayandict/v3/pkg/dictmgr/qdictmgr"
	"github.com/ilius/ayandict/v3/pkg/logging"
	"github.com/ilius/ayandict/v3/pkg/qtcommon"
	"github.com/ilius/ayandict/v3/pkg/qtcommon/qsettings"
	"github.com/ilius/ayandict/v3/pkg/server"
	qt "github.com/mappu/miqt/qt6"
)

var searchModes = []string{
	"Fuzzy",
	"Start with",
	"Regex",
	"Glob",
	"Word Match",
}

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
	searchModeCombo *qt.QComboBox
	favoritesWidget *qfavorites.FavoritesWidget
	frequencyTable  *frequency.FrequencyTable

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
	app.queryArgs.ResultsLabel.SetText("Results")
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

func (app *Application) onResultDisplay(terms []string) {
	app.favoriteButton.Show()
	app.favoriteButton.SetChecked(app.favoritesWidget.HasFavorite(terms[0]))
}

// TODO: break down
func (app *Application) Run() {
	app.init()

	basePx := app.baseFontPixelSize()

	basePxI := int(basePx)
	basePxHalf := int(basePx / 2)

	activityStorage := activity.NewActivityStorage(conf, config.GetConfigDir())

	frequencyTable := frequency.NewFrequencyView(
		activityStorage,
		conf.MostFrequentMaxSize,
	)
	app.frequencyTable = frequencyTable

	// icon := qt.NewQIcon5("./img/icon.png")

	window := app.window
	window.SetWindowTitle(appinfo.APP_DESC)
	window.Resize(600, 400)

	entry := qt.NewQLineEdit2()
	app.entry = entry
	entry.SetPlaceholderText("Type search query and press Enter")
	entry.SetTextMargins(0, -3, 0, -3) // to reduce inner margins

	searchModeCombo := qt.NewQComboBox2()
	app.searchModeCombo = searchModeCombo
	app.searchModeCombo.AddItems(searchModes)

	okButton := qt.NewQPushButton3(" OK ")

	app.queryFavoriteButton = NewFavoriteButton(app.queryFavoriteButtonClicked)
	app.queryFavoriteButton.SetToolTips(
		"Add this query to favorites",
		"Remove this query from favorites",
	)

	// favoriteButtonVBox := qt.NewQVBoxLayout()
	app.favoriteButton = NewFavoriteButton(app.favoriteButtonClicked)

	app.favoriteButton.SetToolTips(
		"Add this term to favorites",
		"Remove this term from favorites",
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
	queryBoxLayout.AddWidget(app.queryFavoriteButton.QWidget)
	queryBoxLayout.AddWidget(okButton.QWidget)

	headerLabel := CreateHeaderLabel(app)
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
	headerBoxLayout.AddWidget3(app.favoriteButton.QWidget, 0, qt.AlignRight)
	headerBoxLayout.AddSpacing(int(basePx * 1.5))
	headerBox.SetSizePolicy2(expanding, qt.QSizePolicy__Minimum)

	articleView := NewArticleView(app)
	app.articleView = articleView

	historyView := NewHistoryView(activityStorage, conf.HistoryMaxSize)
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

	app.updateMiscButtonsVisibility()
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
	activityLayout.AddWidget(historyView.QWidget)
	activityLayout.AddWidget(frequencyTable.QWidget)
	activityLayout.AddWidget(app.favoritesWidget.QWidget)

	activityTypeCombo.OnCurrentIndexChanged(app.activityComboChanged)

	leftPanel := qt.NewQWidget(nil)
	leftPanelLayout := qt.NewQVBoxLayout(leftPanel)
	resultsLabel := qt.NewQLabel3("Results")
	leftPanelLayout.AddWidget(resultsLabel.QWidget)
	resultList := NewResultListWidget(
		articleView,
		headerLabel,
		app.onResultDisplay,
	)
	app.resultList = resultList
	leftPanelLayout.AddWidget(app.resultList.QWidget)

	app.queryArgs = &QueryArgs{
		ArticleView:    articleView,
		ResultsLabel:   resultsLabel,
		ResultList:     resultList,
		HeaderLabel:    headerLabel,
		HistoryView:    historyView,
		PostQuery:      app.postQuery,
		Entry:          entry,
		ModeCombo:      searchModeCombo,
		FrequencyTable: frequencyTable,
	}

	headerLabel.doQuery = app.doQuery
	articleView.doQuery = app.doQuery
	historyView.doQuery = app.doQuery

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
		// local:
		queryLabel,
		okButton,
		aboutButton,
		rightPanel,
		// fields:
		app.frequencyTable,
		app.activityTypeCombo,
		app.entry,
		app.searchModeCombo,
		app.headerLabel,
		app.articleView,
		app.historyView,
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
		app.openConfigButton,
		app.reloadConfigButton,
		app.clearButton,
		app.resultList,
	}
	app.ReloadFont()

	okButton.OnClicked(func() {
		onQuery(entry.Text(), app.queryArgs, false)
	})

	app.setupKeyPressEvent(app.window)
	app.setupKeyPressEvent(app.resultList.QListWidget)
	app.setupKeyPressEvent(app.articleView)
	app.setupKeyPressEvent(app.historyView.QListWidget)

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

func (app *Application) setupSettings(qs *qt.QSettings, mainSplitter *qt.QSplitter) {
	app.searchModeCombo.OnCurrentIndexChanged(func(i int) {
		text := app.entry.Text()
		if text != "" {
			onQuery(text, app.queryArgs, false)
		}
		go qsettings.SaveSearchSettings(qs, app.searchModeCombo)
	})

	qsettings.RestoreSplitterSizes(qs, mainSplitter, QS_mainSplitter)
	qsettings.RestoreMainWinGeometry(qs, app.window)
	qsettings.SetupMainWinGeometrySave(qs, app.window)

	frequencyTable := app.frequencyTable
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

	qsettings.RestoreSearchSettings(qs, app.searchModeCombo)
}

func (app *Application) updateMiscButtonsVisibility() {
	app.saveHistoryButton.SetVisible(conf.MiscButtons.SaveHistory)
	app.clearHistoryButton.SetVisible(conf.MiscButtons.ClearHistory)
	app.saveFavoritesButton.SetVisible(conf.MiscButtons.SaveFavorites)
	app.reloadDictsButton.SetVisible(conf.MiscButtons.ReloadDicts)
	app.closeDictsButton.SetVisible(conf.MiscButtons.CloseDicts)
	app.reloadStyleButton.SetVisible(conf.MiscButtons.ReloadStyle)
	app.randomEntryButton.SetVisible(conf.MiscButtons.RandomEntry)
	app.randomFavoriteButton.SetVisible(conf.MiscButtons.RandomFavorite)
}

func (app *Application) updateMiscButtonsPadding() {
	vpadding := conf.MiscButtonsVerticalPadding
	stylesheet := fmt.Sprintf("padding-top: %dpx; padding-bottom: %dpx;", vpadding, vpadding)

	app.saveHistoryButton.SetStyleSheet(stylesheet)
	app.clearHistoryButton.SetStyleSheet(stylesheet)
	app.saveFavoritesButton.SetStyleSheet(stylesheet)
	app.reloadDictsButton.SetStyleSheet(stylesheet)
	app.closeDictsButton.SetStyleSheet(stylesheet)
	app.reloadStyleButton.SetStyleSheet(stylesheet)
	app.randomEntryButton.SetStyleSheet(stylesheet)
	app.randomFavoriteButton.SetStyleSheet(stylesheet)
}
