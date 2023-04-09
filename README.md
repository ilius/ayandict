# AyanDict

A simple and minimalistic cross-platform desktop dictionary application based on Qt framework and written in Go that uses StarDict dictionary format.

It is designed for desktop and it should run on every desktop operating system that Qt supports. It is tested on Linux and Windows, and it should run perfectly on Mac, FreeBSD and other modern Unix-like systems. I will upload binaries for Linux, Windows and Mac (and maybe FreeBSD).

StarDict is the only supported format for now, and by default, it reads all StarDict dictionaries in `~/.stardict/dic` folder. But you can change the folder or add more folders through [configuration](#configuration).

# Screenshots

<img src="https://raw.githubusercontent.com/wiki/ilius/ayandict/img/linux-light-wordnet.png" width="50%" height="50%"/>

Linux - light style (default)
______________________________________________________________________

<img src="https://raw.githubusercontent.com/wiki/ilius/ayandict/img/linux-dark-fa.png" width="50%" height="50%"/>

Linux - dark style (Breeze)
______________________________________________________________________

<img src="https://raw.githubusercontent.com/wiki/ilius/ayandict/img/windows-light-fa.png" width="50%" height="50%"/>

Windows - light style (default)
______________________________________________________________________

<img src="https://raw.githubusercontent.com/wiki/ilius/ayandict/img/windows-dark-wordnet.png" width="50%" height="50%"/>

Windows - dark style (Breeze)

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

