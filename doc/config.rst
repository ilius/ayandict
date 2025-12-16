``directory_list``
------------------
List of dictionary directory paths (absolute or relative to home)

Default value: ``[".stardict/dic"]``

``style``
---------
Path to application stylesheet file (.qss)

Default value: ``""``

``article_style``
-----------------
Path to article stylesheet file (.css)

Default value: ``""``

``soundex_words_file``
----------------------
Soundex words file

Default value: ``""``

``font_family``
---------------
Application font family

Default value: ``""``

``font_size``
-------------
Application font size

Default value: ``0``

``search_on_type``
------------------
Enable/disable search-on-type

Default value: ``false``

``search_on_type_min_length``
-----------------------------
Minimum query length for search-on-type

Default value: ``3``

``search_on_type_on_regex``
---------------------------
Enable/disable search-on-type in Regex mode

Default value: ``false``

``header_template``
-------------------
HTML template for header (dict name + entry terms)

Default value: ``"<b><font color='#55f'>{{.DictName}}</font></b>\n<font color='#777'> [Score: %{{.Score}}]</font>\n{{if .ShowTerms }}\n<span dir=\"ltr\" style=\"font-size: xx-large;font-weight:bold;\">\n{{ index .Terms 0 }}\n</span>\n{{range slice .Terms 1}}\n<span dir=\"ltr\" style=\"font-size: large;font-weight:bold;\">\n\t<span style=\"color:#ff0000;font-weight:bold;\"> │ </span>\n\t{{ . }}\n</span>\n{{end}}\n{{end}}"``

``header_word_wrap``
--------------------
Enable word-wrapping for header (dict name + entry terms)

Default value: ``true``

``history_disable``
-------------------
Disable history

Default value: ``false``

``history_auto_save``
---------------------
Auto-save history on every new record

Default value: ``true``

``history_max_size``
--------------------
Maximum size for history

Default value: ``100``

``most_frequent_disable``
-------------------------
Disable keeping Most Frequent queries

Default value: ``false``

``most_frequent_auto_save``
---------------------------
Auto-save Most Frequent queries

Default value: ``true``

``most_frequent_max_size``
--------------------------
Maximum size for Most Frequent queries

Default value: ``100``

``favorite_button_image``
-------------------------
Image file name for favorite button. Try favorite-blue-64.png for color blindness

Default value: ``"favorite-green-64.png"``

``favorite_button_hue``
-----------------------
Color hue for favorite button in scan popup (120=green, 240=blue, 0=red)

Default value: ``120``

``favorites_auto_save``
-----------------------
Auto-save Favorites on every new record

Default value: ``true``

``max_results_total``
---------------------
Maximum number of search results

Default value: ``40``

``random_favorite_search_mode``
-------------------------------
Search mode for Random Favorite

Default value: ``"wordMatch"``

``start_hidden``
----------------
Hide main window on startup (if tray icon is available)

Default value: ``false``

``desktop_widget``
------------------
Desktop Widget: enable

Default value: ``true``

``desktop_widget_click_time``
-----------------------------
Dektop Widget: max click time in millisecond

Default value: ``100``

``scan_popup_clipboard``
------------------------
Scan Popup: activate on copy to clipboard

Default value: ``false``

``scan_popup_selection``
------------------------
Scan Popup: activate on selection

Default value: ``false``

``scan_popup_api``
------------------
Scan Popup: activate via API (socket/pipe)

Default value: ``true``

``scan_popup_mode``
-------------------
Scan Popup: search mode

Default value: ``"fuzzy"``

``scan_popup_min_score``
------------------------
Scan Popup: minimum score (0 to 100)

Default value: ``0``

``scan_popup_width``
--------------------
Scan Popup: window width

Default value: ``700``

``scan_popup_height``
---------------------
Scan Popup: window height

Default value: ``400``

``scan_popup_max_count``
------------------------
Scan Popup: max number of windows

Default value: ``3``

``scan_popup_history``
----------------------
Scan Popup: add to history

Default value: ``true``

``scan_popup_header_icons``
---------------------------
Scan Popup: use icons for header buttons

Default value: ``true``

``scan_popup_header_template``
------------------------------
Scan Popup: HTML template for header (dict name and score)

Default value: ``"<span style='font-weight:200;'>{{.DictName}} [Score: %{{.Score}}]</span>"``

``scan_popup_terms_style``
--------------------------
Scan Popup: CSS style string for terms

Default value: ``"style=\"font-size:large;font-weight:bold;\""``

``scan_popup_font_size_factor``
-------------------------------
Scan Popup: font size factor (relative to app)

Default value: ``0.8``

``scan_popup_bypass_window_manager``
------------------------------------
Scan Popup: bypass window manager

Default value: ``true``

``random_favorite_popup_interval_seconds``
------------------------------------------
Show a random favorite term popup every N seconds (0 to disable)

Default value: ``0``

``audio``
---------
Enable audio in article

Default value: ``true``

``audio_mpv``
-------------
Use ``mpv`` command for playing audio

Default value: ``false``

``audio_download_timeout``
--------------------------
Timeout for downloading audio files

Default value: ``"1s"``

``audio_auto_play``
-------------------
Number of audio file to auto-play, set ``0`` to disable.

Default value: ``1``

``audio_auto_play_wait_between``
--------------------------------
Wait time between multiple audio files on auto-play

Default value: ``"800ms"``

``audio_auto_play_min_socre``
-----------------------------
Minimum score for audio auto-play

Default value: ``100``

``audio_volume``
----------------
Volume for playing audio, 0 to 100 (% multiplied by dict-specofic volume)

Default value: ``70``

``embed_external_stylesheet``
-----------------------------
Embed external stylesheet/css in article

Default value: ``false``

``resource_http_download_timeout``
----------------------------------
Timeout for downloading http/https resources in article

Default value: ``"2s"``

``color_mapping``
-----------------
Mapping for colors used in article

Default value: ``{}``

``loading_popup_style_str``
---------------------------
Stylesheet (text) for 'Loading' popup

Default value: ``"border: 1px solid red; background-color: #333; color: white"``

``article_zoom_factor``
-----------------------
Zoom factor for article with mouse wheel or keyboard

Default value: ``1.1``

``article_arrow_keys``
----------------------
Use arrow keys to scroll through article (when focused)

Default value: ``false``

``misc_buttons_vertical_padding``
---------------------------------
Misc buttons vertical padding

Default value: ``5``

``reduce_minimum_window_width``
-------------------------------
Use smaller buttons to reduce minimum width of window

Default value: ``false``

``local_server_ports``
----------------------
Ports for local server. Server runs on first port; Client tries all

Default value: ``["8357"]``

``local_client_timeout``
------------------------
Timeout for local web client

Default value: ``"100ms"``

``web_enable``
--------------
Set true/false and restart to enable/disable web service & web app

Default value: ``false``

``web_expose``
--------------
Expose web service & web app to outside (otherwise only available to 127.0.0.1)

Default value: ``false``

``web_search_on_type``
----------------------
Web: Enable/disable search-on-type

Default value: ``false``

``web_search_on_type_min_length``
---------------------------------
Web: Minimum query length for search-on-type

Default value: ``3``

``web_search_on_type_on_regex``
-------------------------------
Web: Enable/disable search-on-type in Regex mode

Default value: ``false``

``web_show_powered_by``
-----------------------
Show 'Powered By ...' footer in web.

Default value: ``true``

``search_worker_count``
-----------------------
The number of workers / goroutines used for search

Default value: ``8``

``search_timeout``
------------------
Timeout for search on each dictionary. Only works if ``search_worker_count > 1``

Default value: ``"5s"``

``logging.no_color``
--------------------
Disable log colors

Default value: ``false``

``logging.level``
-----------------
Log level

Default value: ``"info"``

``misc_buttons.save_history``
-----------------------------
Show Save History button

Default value: ``true``

``misc_buttons.clear_history``
------------------------------
Show Clear History button

Default value: ``true``

``misc_buttons.save_favorites``
-------------------------------
Show Save Favorites button

Default value: ``true``

``misc_buttons.close_dicts``
----------------------------
Show Close Dicts button

Default value: ``true``

``misc_buttons.random_entry``
-----------------------------
Show Random Entry button

Default value: ``true``

``misc_buttons.random_favorite``
--------------------------------
Show Random Favorite button

Default value: ``true``

