# AyanDict

A simple cross-platform desktop dictionary application based on Qt framework and written in Go that uses StarDict dictionary format.

It is designed for desktop and it should run on every desktop operating system that Qt supports. It is tested on Linux and Windows, and it should run perfectly on Mac, FreeBSD and other modern Unix-like systems. I will upload binaries for Linux, Windows and Mac (and maybe FreeBSD).

StarDict is the only supported format for now, and by default, it reads all StarDict dictionaries in `~/.stardict/dic` folder. But you can change the folder or add more folders through [configuration](#configuration).

# Installation

If you don't have Go language on your system, you can check [Releases](https://github.com/ilius/ayandict/releases) and download the latest binary for your platform if available.

If you have Go, you can compile and install the latest code with

```sh
go install github.com/ilius/ayandict/v2@latest
```

Or clone the repository, `cd` to it and run `go build`, which will create the binary (`ayandict.exe` or `ayandict`) in this directory.

It's good to know that the binary / executable file is completely portable, so you can copy it anywhere you want and run it from there (although on Unix the storage must support executable files).

# Screenshots

<img src="https://raw.githubusercontent.com/wiki/ilius/ayandict/img/v20-linux-light-wordnet.png" width="70%" height="70%"/>

Linux - light style (default)

______________________________________________________________________

<img src="https://raw.githubusercontent.com/wiki/ilius/ayandict/img/v20-linux-dark-wordnet.png" width="70%" height="70%"/>

Linux - dark style (Breeze) + Favorites

______________________________________________________________________

<img src="https://raw.githubusercontent.com/wiki/ilius/ayandict/img/v20-windows-light-wordnet.png" width="70%" height="70%"/>

Windows - light style (default)

______________________________________________________________________

<img src="https://raw.githubusercontent.com/wiki/ilius/ayandict/img/v20-windows-dark-wordnet.png" width="70%" height="70%"/>

Windows - dark style (Breeze)

______________________________________________________________________

<img src="https://raw.githubusercontent.com/wiki/ilius/ayandict/img/v20-linux-dark-frequent-wordnet.png" width="70%" height="70%"/>

Most Frequent queries

______________________________________________________________________

<img src="https://raw.githubusercontent.com/wiki/ilius/ayandict/img/v20-linux-dark-dict-manager.png" width="70%" height="70%"/>

Dictionaries dialog

______________________________________________________________________

<img src="https://raw.githubusercontent.com/wiki/ilius/ayandict/img/v20-linux-dark-misc-empty.png" width="70%" height="70%"/>

Misc tab

______________________________________________________________________

# Configuration

To change configuration (which includes most user settings), you have to edit the config file (we do not have GUI for it, and no plan to add it, sorry!).

After you run the program, you can click on "Config" button (as seen in screenshots) and it will open the `config.toml` file in your default text editor (for TOML files).

If `config.toml` does not exist, it will be created and filled with default config.

After you modify `config.toml`, you can click on "Reload" button (next to "Config" button) and it will apply the changes.

The full path for `config.toml` file:

- Linux: `~/.config/ayandict/config.toml`

  - If `$XDG_CONFIG_HOME` is set: `$XDG_CONFIG_HOME/ayandict/config.toml`

- Windows: `C:\Users\USERNAME\AppData\Roaming\AyanDict\config.toml`

  - More accurately: `%APPDATA%\AyanDict\config.toml`

- Mac: `~/Library/Preferences/AyanDict/config.toml`

Here is a [list of all config parameters](./doc/config.rst).

# Dictionaries

As you see in screenshots, there is a button called "Dictionaries". It opens a dialog and lets you disable, enable and change order of dictionaries.

Each dictionary has a "Symbol" which by default is the first letter of its name in curly brackets (for example `[W]` for WordNet). This symbol is shown in the list of results that is in the left side of window, as seen in screenshots. It is meant to show you which dictionary it comes from at first glance. You can change this symbol through "Dictionaries" dialog. Symbol can be empty, or be as long as you want (though it is 3 characters by default).

# Convert other Dictionary formats

You can use [PyGlossary](https://github.com/ilius/pyglossary) to convert various other formats to StarDict format and use them for this application. A [list of supported formats](https://github.com/ilius/pyglossary#supported-formats) is provided, and if you click on each format's link, it will lead you to more information about it.

# Download Dictionaries

There are tons of web pages that let you download various usable dictionaries, but here is a list I collected (feel free to open a pull request for more):

- [kaikki.org](https://kaikki.org/dictionary/index.html)
- [library.kiwix.org](https://library.kiwix.org/)
- [freedict.org](https://freedict.org/downloads/) and [@freedict/fd-dictionaries](https://github.com/freedict/fd-dictionaries)
- [@itkach/slob/wiki/Dictionaries](https://github.com/itkach/slob/wiki/Dictionaries)
- [goldendict.org](http://goldendict.org/dictionaries.php)
- [huzheng.org](http://www.huzheng.org/stardict/)
- My repositories: [@ilius/dict](https://github.com/ilius/dict) and [Persian Aryanpour in FreeDict](https://github.com/ilius/aryanpour-tei)
- [tuxor1337.frama.io](https://tuxor1337.frama.io/firedict/dictionaries.html)
- [XDXF on SourceForge](https://sourceforge.net/projects/xdxf/files/)
- [GoldenDict on SourceForge](https://sourceforge.net/projects/goldendict/files/dictionaries/)
- [kdr2.com](https://kdr2.com/resource/stardict.html)

# Keyboard bindings/shortcuts

- **Escape**: clear the input query and results
- **Space**: (while query entry is not focused) change keyboard focus to query entry
- **`+`** or **`=`**: Zoom in (article/definition/translation)
- **`-`**: Zoom out (article/definition/translation)

# Search Algorithm

The default search is fuzzy, and it is based on similarity scores that are calculated from [Levenshtein distance](https://en.wikipedia.org/wiki/Levenshtein_distance). We also split entry terms into words, for example if you type "language" (or with with a few misspelled letters, like "languge"), it first shows "language", and then terms like "language learning", but may also show terms like "sign language".

If you specifically want terms with "language" as the second word, you can type "\* language". We do not support pattern matching (yet), and you can only use `*` alone (not as part of a pattern).

Anything with at least %70 similarity score is listed (for example "languge" is %87 similar to "language"). But we have a limit of how many results are displayed, and by default it's 40 results. You can change this with config parameter [`max_results_total`](./doc/config.rst#max_results_total).

This works pretty well in most cases, but the only catch is that first letter of your query must match the first letter of one of your target words. For example if you type "symmetry", it will not match term "asymmetry" even though they are close enough (high similarity score), because their first letter is different.

But we also have 3 other search modes added in v2.0.0:

- Start with, shows all terms that start with given string
- Regex (regular expression), for example `.*symm.*`
- Glob, for example `*symm*`
