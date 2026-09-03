# Scan Popup on macOS

"Scan Popup" is a small popup window that shows dictionary results for a word, without opening the main window. This document describes how to trigger it on macOS with a keyboard shortcut, using the helper scripts in `cmd/mac/`.

## How it works

The keyboard-shortcut approach works over a local Unix socket. While AyanDict is running, it listens on `/tmp/ayandict-$UID` (where `$UID` is your user id) and accepts a `scanpopup:<query>` command, which opens a Scan Popup for the given query.

Using this requires:

- AyanDict to be running (the socket exists only while the program is running)
- The `scan_popup_api` config option to be enabled (default: `true`)
- `nc` (netcat), to connect to the Unix socket; the built-in macOS `nc` supports `-U` for Unix sockets
- The clipboard scripts use the built-in `pbpaste`/`pbcopy` commands, so no extra packages are needed

## Helper scripts

The repository provides these scripts in `cmd/mac/`:

- [scan-selection.sh](https://github.com/ilius/ayandict/blob/v3/cmd/mac/scan-selection.sh): captures the currently selected text (by simulating Cmd+C) and scans it
- [scan-clipboard.sh](https://github.com/ilius/ayandict/blob/v3/cmd/mac/scan-clipboard.sh): scans the text currently on the clipboard

Run them from the repository with their absolute path, e.g.:

```sh
/path/to/ayandict/cmd/mac/scan-selection.sh
/path/to/ayandict/cmd/mac/scan-clipboard.sh
```

These scripts are also bundled in the macOS `.dmg` inside the `Scripts/` folder, together with this document and the other docs in `Docs/`. The scripts are self-contained and can be run directly from the mounted DMG, but since a DMG mounts as read-only, copying `Scripts/` to a writable location (e.g. `/usr/local/bin/`) is recommended if you plan to use them with an Automator Quick Action (see below).

### Permissions

`scan-selection.sh` uses `osascript` to simulate the Cmd+C keystroke, so it needs **Accessibility** permission. Grant it in *System Settings → Privacy & Security → Accessibility* (the terminal app running the script must be enabled there).

## Setting up a keyboard shortcut

macOS has no built-in way to assign a global shortcut directly to a shell script, so create a Quick Action (Service) that runs the script, then bind a shortcut to it:

1. Open **Automator** and create a new **Quick Action**.
1. Set "Workflow receives" to *no input* in *any application*.
1. Add a **Run Shell Script** action (shell: `/bin/bash`) and enter the absolute path to the script, for example `/path/to/ayandict/cmd/mac/scan-selection.sh`.
1. Save it under a name like "AyanDict Scan Selection".
1. Open *System Settings → Keyboard → Keyboard Shortcuts → Services*, find the Quick Action under *General*, and assign it a keyboard shortcut.

For "scan selection" to work, keep the text selected when you press the shortcut (the script copies it via Cmd+C). For the clipboard variant, copy the text first.
