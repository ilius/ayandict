package config

import (
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/ilius/ayandict/v2/pkg/qerr"
)

const fileName = "config.toml"

var mutex sync.Mutex

func init() {
	dir := GetConfigDir()
	err := os.MkdirAll(dir, 0o755)
	if err != nil {
		qerr.Error(err)
	}
}

type Config struct {
	DirectoryList []string `toml:"directory_list"`

	SqlDictList []string `toml:"sql_dict_list"`

	Style string `toml:"style"`

	ArticleStyle string `toml:"article_style"`

	FontFamily string `toml:"font_family"`
	FontSize   int    `toml:"font_size"`

	SearchOnType          bool `toml:"search_on_type"`
	SearchOnTypeMinLength int  `toml:"search_on_type_min_length"`

	HeaderTemplate string `toml:"header_template"`
	HeaderWordWrap bool   `toml:"header_word_wrap"`

	HistoryDisable  bool `toml:"history_disable"`
	HistoryAutoSave bool `toml:"history_auto_save"`
	HistoryMaxSize  int  `toml:"history_max_size"`

	MostFrequentDisable  bool `toml:"most_frequent_disable"`
	MostFrequentAutoSave bool `toml:"most_frequent_auto_save"`
	MostFrequentMaxSize  int  `toml:"most_frequent_max_size"`

	FavoritesAutoSave bool `toml:"favorites_auto_save"`

	MaxResultsTotal int `toml:"max_results_total"`

	Audio bool `toml:"audio"`

	AudioMPV bool `toml:"audio_mpv"`

	AudioDownloadTimeout time.Duration `toml:"audio_download_timeout"`

	AudioAutoPlay int `toml:"audio_auto_play"`

	AudioAutoPlayWaitBetween time.Duration `toml:"audio_auto_play_wait_between"`

	EmbedExternalStylesheet bool `toml:"embed_external_stylesheet"`

	ColorMapping map[string]string `toml:"color_mapping"`

	PopupStyleStr string `toml:"popup_style_str"`

	ArticleZoomFactor float64 `toml:"article_zoom_factor"`

	ArticleArrowKeys bool `toml:"article_arrow_keys"`

	ReduceMinimumWindowWidth bool `toml:"reduce_minimum_window_width"`

	LocalServerPorts []string `toml:"local_server_ports"`

	LocalClientTimeout time.Duration `toml:"local_client_timeout"`

	WebEnable bool `toml:"web_enable"`
	WebExpose bool `toml:"web_expose"`

	WebShowPoweredBy bool `toml:"web_show_powered_by"`

	SearchWorkerCount int `toml:"search_worker_count"`

	SearchTimeout time.Duration `toml:"search_timeout"`
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
		DirectoryList: []string{
			".stardict/dic",
		},

		Style: "",

		ArticleStyle: "",

		FontFamily: "",
		FontSize:   0,

		SearchOnType:          false,
		SearchOnTypeMinLength: 3,

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

		EmbedExternalStylesheet: false,

		ColorMapping: map[string]string{},

		PopupStyleStr: "border: 1px solid red; background-color: #333; color: white",

		ArticleZoomFactor: 1.1,

		ArticleArrowKeys: false,

		ReduceMinimumWindowWidth: false,

		LocalServerPorts: []string{
			"8357",
		},

		LocalClientTimeout: 100 * time.Millisecond,

		WebEnable:        false,
		WebExpose:        false,
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
			log.Printf("bad CONFIG_FILE=%#v, error: %v", _path, err)
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
