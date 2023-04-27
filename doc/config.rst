``directory_list``
------------------
List of dictionary directory paths (absolute or relative to home)

Default value: ``[".stardict/dic"]``

``sql_dict_list``
-----------------
List of SQL databases, only SQLite is currently supported

Default value: ``null``

``style``
---------
Path to application stylesheet file (.qss)

Default value: ``""``

``article_style``
-----------------
Path to article stylesheet file (.css)

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

``header_template``
-------------------
HTML template for header (dict name + entry terms)

Default value: ``"<b><font color='#55f'>{{.DictName}}</font></b>\n<font color='#777'> [Score: %{{.Score}}]</font>\n{{if .ShowTerms }}\n<div dir=\"ltr\" style=\"font-size: xx-large;font-weight:bold;\">\n{{ index .Terms 0 }}\n</div>\n{{range slice .Terms 1}}\n<span dir=\"ltr\" style=\"font-size: large;font-weight:bold;\">\n\t<span style=\"color:#ff0000;font-weight:bold;\"> | </span>\n\t{{ . }}\n</span>\n{{end}}\n{{end}}"``

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

``favorites_auto_save``
-----------------------
Auto-save Favorites on every new record

Default value: ``true``

``max_results_total``
---------------------
Maximum number of search results

Default value: ``40``

``audio``
---------
Enable audio in article

Default value: ``true``

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

Default value: ``"500ms"``

``embed_external_stylesheet``
-----------------------------
Embed external stylesheet/css in article

Default value: ``false``

``color_mapping``
-----------------
Mapping for colors used in article

Default value: ``{}``

``popup_style_str``
-------------------
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

``search_worker_count``
-----------------------
The number of workers / goroutines used for search

Default value: ``8``

``search_timeout``
------------------
Timeout for search on each dictionary. Only works if ``search_worker_count > 1``

Default value: ``"5s"``

