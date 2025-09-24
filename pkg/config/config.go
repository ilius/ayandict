package config

import (
	"log/slog"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/BurntSushi/toml"
)

const fileName = "config.toml"

var mutex sync.Mutex

var PrivateMode = false

func init() {
	dir := GetConfigDir()
	err := os.MkdirAll(dir, 0o755)
	if err != nil {
		slog.Error("error in MkdirAll: " + err.Error())
	}
}

type LoggingConfig struct {
	NoColor bool   `toml:"no_color" doc:"Disable log colors"`
	Level   string `toml:"level" doc:"Log level"`
}

type MiscButtonsConfig struct {
	SaveHistory    bool `toml:"save_history" doc:"Show Save History button"`
	ClearHistory   bool `toml:"clear_history" doc:"Show Clear History button"`
	SaveFavorites  bool `toml:"save_favorites" doc:"Show Save Favorites button"`
	ReloadDicts    bool `toml:"reload_dicts" doc:"Show Reload Dicts button"`
	CloseDicts     bool `toml:"close_dicts" doc:"Show Close Dicts button"`
	ReloadStyle    bool `toml:"reload_style" doc:"Show Reload Style button"`
	RandomEntry    bool `toml:"random_entry" doc:"Show Random Entry button"`
	RandomFavorite bool `toml:"random_favorite" doc:"Show Random Favorite button"`
}

type Config struct {
	Logging LoggingConfig `toml:"logging" doc:"Logging config"`

	MiscButtons MiscButtonsConfig `toml:"misc_buttons" doc:"Misc buttons visibility"`

	DirectoryList []string `toml:"directory_list" doc:"List of dictionary directory paths (absolute or relative to home)"`

	Style string `toml:"style" doc:"Path to application stylesheet file (.qss)"`

	ArticleStyle string `toml:"article_style" doc:"Path to article stylesheet file (.css)"`

	FontFamily string `toml:"font_family" doc:"Application font family"`
	FontSize   int    `toml:"font_size" doc:"Application font size"`

	SearchOnType          bool `toml:"search_on_type" doc:"Enable/disable search-on-type"`
	SearchOnTypeMinLength int  `toml:"search_on_type_min_length" doc:"Minimum query length for search-on-type"`
	SearchOnTypeOnRegex   bool `toml:"search_on_type_on_regex" doc:"Enable/disable search-on-type in Regex mode"`

	HeaderTemplate string `toml:"header_template" doc:"HTML template for header (dict name + entry terms)"`
	HeaderWordWrap bool   `toml:"header_word_wrap" doc:"Enable word-wrapping for header (dict name + entry terms)"`

	HistoryDisable  bool `toml:"history_disable" doc:"Disable history"`
	HistoryAutoSave bool `toml:"history_auto_save" doc:"Auto-save history on every new record"`
	HistoryMaxSize  int  `toml:"history_max_size" doc:"Maximum size for history"`

	MostFrequentDisable  bool `toml:"most_frequent_disable" doc:"Disable keeping Most Frequent queries"`
	MostFrequentAutoSave bool `toml:"most_frequent_auto_save" doc:"Auto-save Most Frequent queries"`
	MostFrequentMaxSize  int  `toml:"most_frequent_max_size" doc:"Maximum size for Most Frequent queries"`

	FavoritesAutoSave bool `toml:"favorites_auto_save" doc:"Auto-save Favorites on every new record"`

	MaxResultsTotal int `toml:"max_results_total" doc:"Maximum number of search results"`

	Audio bool `toml:"audio" doc:"Enable audio in article"`

	AudioMPV bool `toml:"audio_mpv" doc:"Use ‘mpv‘ command for playing audio"`

	AudioDownloadTimeout time.Duration `toml:"audio_download_timeout" doc:"Timeout for downloading audio files"`

	AudioAutoPlay int `toml:"audio_auto_play" doc:"Number of audio file to auto-play, set ‘0‘ to disable."`

	AudioAutoPlayWaitBetween time.Duration `toml:"audio_auto_play_wait_between" doc:"Wait time between multiple audio files on auto-play"`

	AudioVolume int `toml:"audio_volume" doc:"Volume for playing audio, 0 to 100 (% multiplied by dict-specofic volume)"`

	EmbedExternalStylesheet bool `toml:"embed_external_stylesheet" doc:"Embed external stylesheet/css in article"`

	ResourceHttpDownloadTimeout time.Duration `toml:"resource_http_download_timeout" doc:"Timeout for downloading http/https resources in article"`

	ColorMapping map[string]string `toml:"color_mapping" doc:"Mapping for colors used in article"`

	PopupStyleStr string `toml:"popup_style_str" doc:"Stylesheet (text) for 'Loading' popup"`

	ArticleZoomFactor float64 `toml:"article_zoom_factor" doc:"Zoom factor for article with mouse wheel or keyboard"`

	ArticleArrowKeys bool `toml:"article_arrow_keys" doc:"Use arrow keys to scroll through article (when focused)"`

	ReduceMinimumWindowWidth bool `toml:"reduce_minimum_window_width" doc:"Use smaller buttons to reduce minimum width of window"`

	LocalServerPorts []string `toml:"local_server_ports" doc:"Ports for local server. Server runs on first port; Client tries all"`

	LocalClientTimeout time.Duration `toml:"local_client_timeout" doc:"Timeout for local web client"`

	WebEnable bool `toml:"web_enable" doc:"Set true/false and restart to enable/disable web service & web app"`
	WebExpose bool `toml:"web_expose" doc:"Expose web service & web app to outside (otherwise only available to 127.0.0.1)"`

	WebSearchOnType          bool `toml:"web_search_on_type" doc:"Web: Enable/disable search-on-type"`
	WebSearchOnTypeMinLength int  `toml:"web_search_on_type_min_length" doc:"Web: Minimum query length for search-on-type"`
	WebSearchOnTypeOnRegex   bool `toml:"web_search_on_type_on_regex" doc:"Web: Enable/disable search-on-type in Regex mode"`

	WebShowPoweredBy bool `toml:"web_show_powered_by" doc:"Show 'Powered By ...' footer in web."`

	SearchWorkerCount int `toml:"search_worker_count" doc:"The number of workers / goroutines used for search"`

	SearchTimeout time.Duration `toml:"search_timeout" doc:"Timeout for search on each dictionary. Only works if ‘search_worker_count > 1‘"`
}

const defaultHeaderTemplate = `<b><font color='#55f'>{{.DictName}}</font></b>
<font color='#777'> [Score: %{{.Score}}]</font>
{{if .ShowTerms }}
<div dir="ltr" style="font-size: xx-large;font-weight:bold;">
{{ index .Terms 0 }}
</div>
{{range slice .Terms 1}}
<span dir="ltr" style="font-size: large;font-weight:bold;">
	<span style="color:#ff0000;font-weight:bold;"> │ </span>
	{{ . }}
</span>
{{end}}
{{end}}`

func Default() *Config {
	return &Config{
		Logging: LoggingConfig{
			NoColor: false,
			Level:   "info",
		},
		MiscButtons: MiscButtonsConfig{
			SaveHistory:    true,
			ClearHistory:   true,
			SaveFavorites:  true,
			ReloadDicts:    true,
			CloseDicts:     true,
			ReloadStyle:    true,
			RandomEntry:    true,
			RandomFavorite: true,
		},
		DirectoryList: []string{
			".stardict/dic",
		},

		Style: "",

		ArticleStyle: "",

		FontFamily: "",
		FontSize:   0,

		SearchOnType:          false,
		SearchOnTypeMinLength: 3,
		SearchOnTypeOnRegex:   false,

		HeaderTemplate: defaultHeaderTemplate,
		HeaderWordWrap: true,

		HistoryDisable:  false,
		HistoryAutoSave: true,
		HistoryMaxSize:  100,

		MostFrequentDisable:  false,
		MostFrequentAutoSave: true,
		MostFrequentMaxSize:  100,

		FavoritesAutoSave: true,

		MaxResultsTotal: 40,

		Audio: true,

		AudioMPV: false,

		AudioDownloadTimeout: 1000 * time.Millisecond,

		AudioAutoPlay: 1,

		AudioAutoPlayWaitBetween: 500 * time.Millisecond,

		AudioVolume: 70,

		EmbedExternalStylesheet: false,

		ResourceHttpDownloadTimeout: 2 * time.Second,

		ColorMapping: map[string]string{},

		PopupStyleStr: "border: 1px solid red; background-color: #333; color: white",

		ArticleZoomFactor: 1.1,

		ArticleArrowKeys: false,

		ReduceMinimumWindowWidth: false,

		LocalServerPorts: []string{
			"8357",
		},

		LocalClientTimeout: 100 * time.Millisecond,

		WebEnable: false,
		WebExpose: false,

		WebSearchOnType:          false,
		WebSearchOnTypeMinLength: 3,
		WebSearchOnTypeOnRegex:   false,

		WebShowPoweredBy: true,

		SearchWorkerCount: 8,

		SearchTimeout: 5 * time.Second,
	}
}

func Path() string {
	_path := os.Getenv("CONFIG_FILE")
	if _path != "" {
		absPath, err := filepath.Abs(_path)
		if err == nil {
			return absPath
		} else {
			slog.Error("bad CONFIG_FILE", "CONFIG_FILE", _path, "err", err)
		}
	}
	return filepath.Join(GetConfigDir(), fileName)
}

func GetConfigDir() string {
	if os.Getenv("CONFIG_FILE") != "" {
		return filepath.Dir(Path())
	}
	return platformConfigDir()
}

func loadFile() ([]byte, error) {
	mutex.Lock()
	defer mutex.Unlock()

	pathStr := Path()
	tomlBytes, err := os.ReadFile(pathStr)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	return tomlBytes, nil
}

func Load() (*Config, error) {
	tomlBytes, err := loadFile()
	if err != nil {
		return nil, err
	}
	if tomlBytes == nil {
		return Default(), nil
	}
	conf := Default()
	_, err = toml.Decode(string(tomlBytes), conf)
	if err != nil {
		return nil, err
	}
	return conf, nil
}

func MustLoad() *Config {
	conf, err := Load()
	if err != nil {
		panic(err)
	}
	return conf
}
