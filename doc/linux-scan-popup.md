# Scan Popup on Linux

"Scan Popup" is a small popup window that shows dictionary results for a word, without opening the main window. This document describes how to trigger it on Linux with a keyboard shortcut, using the helper scripts in `cmd/linux/`.

## How it works

The keyboard-shortcut approach works over a local Unix socket. While AyanDict is running, it listens on `/tmp/ayandict-$UID` (where `$UID` is your user id) and accepts a `scanpopup:<query>` command, which opens a Scan Popup for the given query.

Using this requires:

- AyanDict to be running (the socket exists only while the program is running)
- The `scan_popup_api` config option to be enabled (default: `true`; it can also be toggled from the tray icon menu, "Scan via API")
- `socat` or `netcat` (`nc`), to connect to the Unix socket
- `wl-clipboard` (on Wayland) or `xclip` (on X11), to read the selection/clipboard

On Ubuntu, `sudo apt install wl-clipboard socat` covers the Wayland clipboard and the Unix socket connection.

## Helper scripts

The repository provides these scripts in `cmd/linux/`:

- [scan-selection.sh](/cmd/linux/scan-selection.sh): scans the currently selected text (the "primary selection")
- [scan-clipboard.sh](/cmd/linux/scan-clipboard.sh): scans the text currently on the clipboard
- [status-icon-activate.sh](/cmd/linux/status-icon-activate.sh): simulates a click on the tray icon (shows/hides the main window)

### When installed from a deb/rpm package

The scripts are installed to `/usr/lib/ayandict/`:

```sh
/usr/lib/ayandict/scan-selection.sh
/usr/lib/ayandict/scan-clipboard.sh
/usr/lib/ayandict/status-icon-activate.sh
```

### When running the raw binary and scripts from the repository

Use the scripts in `cmd/linux/`:

```sh
./cmd/linux/scan-selection.sh
./cmd/linux/scan-clipboard.sh
./cmd/linux/status-icon-activate.sh
```

For a desktop shortcut, use the absolute path to the script (e.g. `/path/to/ayandict/cmd/linux/scan-selection.sh`), since the desktop environment does not inherit your shell's working directory.

## Setting up a keyboard shortcut

Bind one of the scripts above to a key combination in your desktop environment:

- GNOME (and other gsettings-based desktops): *Settings → Keyboard → View and Customize Shortcuts → Custom Shortcuts*, then click the **+** button. Alternatively, on the command line:

  ```sh
  gsettings set org.gnome.settings-daemon.plugins.media-keys custom-keybindings "['/org/gnome/settings-daemon/plugins/media-keys/custom-keybindings/custom0/']"
  gsettings set org.gnome.settings-daemon.plugins.media-keys.custom-keybinding:/org/gnome/settings-daemon/plugins/media-keys/custom-keybindings/custom0/ name 'AyanDict: scan selection'
  gsettings set org.gnome.settings-daemon.plugins.media-keys.custom-keybinding:/org/gnome/settings-daemon/plugins/media-keys/custom-keybindings/custom0/ command '/usr/lib/ayandict/scan-selection.sh'
  gsettings set org.gnome.settings-daemon.plugins.media-keys.custom-keybinding:/org/gnome/settings-daemon/plugins/media-keys/custom-keybindings/custom0/ binding '<Ctrl><Alt>d'
  ```

  If the `custom-keybindings` list already contains entries, append the new path to the existing array instead of overwriting it.

- KDE Plasma: *System Settings → Shortcuts → Add New → Global Shortcut → Command/URL*, enter the script path as the command and set a trigger.

- XFCE: *Settings → Keyboard → Application Shortcuts → Add*, enter the script path and press the desired key combination.

For "scan selection" to work, keep the text selected when you press the shortcut, so the script can read the primary selection. For the clipboard variant, copy the text first.
