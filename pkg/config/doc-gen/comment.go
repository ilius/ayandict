package main

var commentMap = map[string]string{
	"logging":        "Logging config",
	"directory_list": "List of dictionary directory paths (absolute or relative to home)",

	"sql_dict_list": "List of SQL databases, only SQLite is currently supported",

	"style":         "Path to application stylesheet file (.qss)",
	"article_style": "Path to article stylesheet file (.css)",
	"font_family":   "Application font family",
	"font_size":     "Application font size",

	"search_on_type":            "Enable/disable search-on-type",
	"search_on_type_min_length": "Minimum query length for search-on-type",

	"header_template":  "HTML template for header (dict name + entry terms)",
	"header_word_wrap": "Enable word-wrapping for header (dict name + entry terms)",

	"history_disable":   "Disable history",
	"history_auto_save": "Auto-save history on every new record",
	"history_max_size":  "Maximum size for history",

	"most_frequent_disable":   "Disable keeping Most Frequent queries",
	"most_frequent_auto_save": "Auto-save Most Frequent queries",
	"most_frequent_max_size":  "Maximum size for Most Frequent queries",

	"favorites_auto_save": "Auto-save Favorites on every new record",

	"max_results_total": "Maximum number of search results",

	"audio": "Enable audio in article",

	"audio_mpv": "Use `mpv` command for playing audio",

	"audio_download_timeout": "Timeout for downloading audio files",

	"audio_auto_play": "Number of audio file to auto-play, set `0` to disable.",

	"audio_auto_play_wait_between": "Wait time between multiple audio files on auto-play",

	"embed_external_stylesheet": "Embed external stylesheet/css in article",

	"color_mapping": "Mapping for colors used in article",

	"popup_style_str": "Stylesheet (text) for 'Loading' popup",

	"article_zoom_factor": "Zoom factor for article with mouse wheel or keyboard",

	"article_arrow_keys": "Use arrow keys to scroll through article (when focused)",

	"reduce_minimum_window_width": "Use smaller buttons to reduce minimum width of window",

	"local_server_ports": "Ports for local server. Server runs on first port; Client tries all",

	"local_client_timeout": "Timeout for local web client",

	"web_enable": "Set true/false and restart to enable/disable web service & web app",
	"web_expose": "Expose web service & web app to outside (otherwise only available to 127.0.0.1)",

	"web_search_on_type":            "Web: Enable/disable search-on-type",
	"web_search_on_type_min_length": "Web: Minimum query length for search-on-type",

	"web_show_powered_by": "Show 'Powered By ...' footer in web.",

	"search_worker_count": "The number of workers / goroutines used for search",

	"search_timeout": "Timeout for search on each dictionary. Only works if `search_worker_count > 1`",
}
