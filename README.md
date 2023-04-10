# AyanDict

A simple and minimalistic cross-platform desktop dictionary application based on Qt framework and written in Go that uses StarDict dictionary format.

It is designed for desktop and it should run on every desktop operating system that Qt supports. It is tested on Linux and Windows, and it should run perfectly on Mac, FreeBSD and other modern Unix-like systems. I will upload binaries for Linux, Windows and Mac (and maybe FreeBSD).

StarDict is the only supported format for now, and by default, it reads all StarDict dictionaries in `~/.stardict/dic` folder. But you can change the folder or add more folders through [configuration](#configuration).

# Installation

If you don't have Go langauge on your system, you can check [Releases](https://github.com/ilius/ayandict/releases) and download the latest binary for your platform if available.

If you have Go, you can compile install the latest code with

```sh
go install github.com/ilius/ayandict@latest
```

Or clone the reposotory, `cd` to it and run `go build`, which will create the binary (`ayandict.exe` or `ayandict`) in this directory.

# Screenshots

<img src="https://raw.githubusercontent.com/wiki/ilius/ayandict/img/linux-light-wordnet.png" width="70%" height="70%"/>

Linux - light style (default)

______________________________________________________________________

<img src="https://raw.githubusercontent.com/wiki/ilius/ayandict/img/linux-dark-fa.png" width="70%" height="70%"/>

Linux - dark style (Breeze)

______________________________________________________________________

<img src="https://raw.githubusercontent.com/wiki/ilius/ayandict/img/windows-light-fa.png" width="70%" height="70%"/>

Windows - light style (default)

______________________________________________________________________

<img src="https://raw.githubusercontent.com/wiki/ilius/ayandict/img/windows-dark-wordnet.png" width="70%" height="70%"/>

Windows - dark style (Breeze)

______________________________________________________________________

<img src="https://raw.githubusercontent.com/wiki/ilius/ayandict/img/linux-light-wordnet-frequent.png" width="70%" height="70%"/>
<img src="https://raw.githubusercontent.com/wiki/ilius/ayandict/img/linux-light-wordnet-favorites.png" width="70%" height="70%"/>

Most Frequent queries and Favorites

______________________________________________________________________

<img src="https://raw.githubusercontent.com/wiki/ilius/ayandict/img/linux-light-dicts.png" width="70%" height="70%"/>

Dictionaries dialog

______________________________________________________________________

<img src="https://raw.githubusercontent.com/wiki/ilius/ayandict/img/linux-light-misc-empty.png" width="70%" height="70%"/>

Misc tab

______________________________________________________________________

# Configuration

To change configuration (which includes most user settings), you have to edit the config file (we do not have GUI for it, and no plan to add it, sorry!).

After you run the program, you can click on "Config" button (as seen in screenshots) and it will open the `config.toml` file in your default text editor (for TOML files).

If `config.toml` does not exist, it will be created and filled with default config.

After you modify `config.toml`, you can click on "Reload" button (next to "Config" button) and it will apply the changes (with exception of font probably).

The full path for `config.toml` file:

- Linux: `~/.config/ayandict/config.toml`

  - If `$XDG_CONFIG_HOME` is set: `$XDG_CONFIG_HOME/ayandict/config.toml`

- Windows: `C:\Users\USERNAME\AppData\Roaming\AyanDict\config.toml`

  - More accurately: `%APPDATA%\AyanDict\config.toml`

- Mac: `~/Library/Preferences/AyanDict/config.toml`

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
- My repos: [@ilius/dict](https://github.com/ilius/dict) and [Persian Aryanpour in FreeDict](https://github.com/ilius/aryanpour-tei)
- [tuxor1337.frama.io](https://tuxor1337.frama.io/firedict/dictionaries.html)
- [XDXF on SourceForge](https://sourceforge.net/projects/xdxf/files/)
- [GoldenDict on SourceForge](https://sourceforge.net/projects/goldendict/files/dictionaries/)
- [kdr2.com](https://kdr2.com/resource/stardict.html)

# Keyboard bindings/shortcuts

- **Escape**: clear the input query and results
- **Space**: (while not typing in query entry) change keyboard focus to query entry
- **`+`** or **`=`**: Zoom in (article/definition/translation)
- **`-`**: Zoom out (article/definition/translation)
