#!/usr/bin/env bash
# Add a global hotkey to run ~/bin/ayandict-scan-selection
# Supports: GNOME, Unity, Budgie, Cinnamon, MATE, KDE, XFCE, LXQt, LXDE, i3, Pantheon.
# Safe to re-run: existing bindings are not duplicated.

set -euo pipefail

CMD="$HOME/bin/ayandict-scan-selection"
SRC_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SRC_SCRIPT="$SRC_DIR/scan-selection.sh"
KEYBIND="Alt+Z"
ACTION_NAME="AyanDict Scan Selection"

# Ensure the target command exists
if [[ ! -x "$CMD" ]]; then
	if [[ -f "$SRC_SCRIPT" ]]; then
		echo "Installing $SRC_SCRIPT → $CMD"
		mkdir -p "$HOME/bin"
		cp "$SRC_SCRIPT" "$CMD"
		chmod +x "$CMD"
	else
		echo "Error: Neither $CMD nor $SRC_SCRIPT found."
		echo "Expected $SRC_SCRIPT to exist beside this script."
		exit 1
	fi
fi

DE="${XDG_CURRENT_DESKTOP:-${DESKTOP_SESSION:-unknown}}"
DE="$(echo "$DE" | tr '[:upper:]' '[:lower:]')"
echo "Detected desktop environment: $DE"

check_dep() {
	if ! command -v "$1" >/dev/null 2>&1; then
		echo "Missing dependency: $1"
		echo "Install it with your package manager, e.g.:"
		echo "  sudo apt install $1"
		echo "  or"
		echo "  sudo dnf install $1"
		echo
		return 1
	fi
	return 0
}

add_gsettings_binding() {
	local SCHEMA="$1"
	local BASE="$2"
	local CMD="$3"
	local KEYBIND="$4"
	local ACTION_NAME="$5"

	local EXISTING_PATHS
	EXISTING_PATHS=$(gsettings get "$SCHEMA" custom-keybindings)

	if gsettings list-recursively | grep -q "$CMD"; then
		echo "Hotkey already exists for $CMD in $SCHEMA."
		return
	fi

	local IDX
	IDX=$(($(echo "$EXISTING_PATHS" | grep -o "custom" | wc -l)))
	local NEW_PATH="${BASE}custom${IDX}/"
	local NEW_LIST
	NEW_LIST=$(echo "$EXISTING_PATHS" | sed "s/]$/, '${NEW_PATH}']/")

	gsettings set "$SCHEMA" custom-keybindings "$NEW_LIST"
	gsettings set "${SCHEMA}.custom-keybinding:${NEW_PATH}" name "$ACTION_NAME"
	gsettings set "${SCHEMA}.custom-keybinding:${NEW_PATH}" command "$CMD"
	gsettings set "${SCHEMA}.custom-keybinding:${NEW_PATH}" binding "$KEYBIND"

	echo "Added '$KEYBIND' → $CMD for $SCHEMA."
}

setup_gnome() {
	check_dep gsettings || exit 1
	local SCHEMA="org.gnome.settings-daemon.plugins.media-keys"
	local BASE="/org/gnome/settings-daemon/plugins/media-keys/custom-keybindings/"
	add_gsettings_binding "$SCHEMA" "$BASE" "$CMD" "$KEYBIND" "$ACTION_NAME"
}

setup_cinnamon() {
	check_dep gsettings || exit 1
	local SCHEMA="org.cinnamon.settings-daemon.plugins.media-keys"
	local BASE="/org/cinnamon/settings-daemon/plugins/media-keys/custom-keybindings/"
	add_gsettings_binding "$SCHEMA" "$BASE" "$CMD" "$KEYBIND" "$ACTION_NAME"
}

setup_mate() {
	check_dep gsettings || exit 1
	local SCHEMA="org.mate.settings-daemon.plugins.media-keys"
	local BASE="/org/mate/desktop/keybindings/"
	add_gsettings_binding "$SCHEMA" "$BASE" "$CMD" "$KEYBIND" "$ACTION_NAME"
}

setup_kde() {
	check_dep qdbus || true
	local CONFIG_DIR="$HOME/.config/kglobalshortcutsrc"
	if grep -q "$CMD" "$CONFIG_DIR" 2>/dev/null; then
		echo "Hotkey already set for KDE Plasma."
	else
		echo "[AyanDict]" >>"$CONFIG_DIR"
		echo "trigger=$KEYBIND" >>"$CONFIG_DIR"
		echo "command=$CMD" >>"$CONFIG_DIR"
		echo "name=$ACTION_NAME" >>"$CONFIG_DIR"
		echo "Hotkey '$KEYBIND' added to KDE Plasma config."
	fi
	echo "You may need to restart KDE or run:"
	echo "  qdbus org.kde.kglobalaccel /component/AyanDict reconfigure"
}

setup_xfce() {
	check_dep xfconf-query || exit 1
	local PROP="/commands/custom/$KEYBIND"
	if xfconf-query -c xfce4-keyboard-shortcuts -p "$PROP" >/dev/null 2>&1; then
		echo "Hotkey already exists for XFCE."
	else
		xfconf-query -c xfce4-keyboard-shortcuts -p "$PROP" -n -t string -s "$CMD"
		echo "Hotkey '$KEYBIND' set for XFCE."
	fi
}

setup_lxqt() {
	check_dep xmlstarlet || exit 1
	local CFG="$HOME/.config/openbox/lxqt-rc.xml"
	if [[ -f "$CFG" ]]; then
		if grep -q "$CMD" "$CFG"; then
			echo "Hotkey already exists for LXQt."
		else
			xmlstarlet ed -L \
				-s "/lxqt_config/globalkeyshortcuts" -t elem -n "keybind" \
				-i "/lxqt_config/globalkeyshortcuts/keybind[not(@key)]" -t attr -n "key" -v "$KEYBIND" \
				-s "/lxqt_config/globalkeyshortcuts/keybind[@key='$KEYBIND']" -t elem -n "command" -v "$CMD" \
				"$CFG"
			echo "Hotkey '$KEYBIND' set for LXQt."
		fi
	else
		echo "LXQt config not found: $CFG"
	fi
}

setup_lxde() {
	check_dep xmlstarlet || exit 1
	local CFG="$HOME/.config/openbox/lxde-rc.xml"
	if [[ -f "$CFG" ]]; then
		if grep -q "$CMD" "$CFG"; then
			echo "Hotkey already exists for LXDE."
		else
			xmlstarlet ed -L \
				-s "/lxde_config/globalkeyshortcuts" -t elem -n "keybind" \
				-i "/lxde_config/globalkeyshortcuts/keybind[not(@key)]" -t attr -n "key" -v "$KEYBIND" \
				-s "/lxde_config/globalkeyshortcuts/keybind[@key='$KEYBIND']" -t elem -n "command" -v "$CMD" \
				"$CFG"
			echo "Hotkey '$KEYBIND' set for LXDE."
		fi
	else
		echo "LXDE config not found: $CFG"
	fi
}

setup_i3() {
	local CFG="$HOME/.config/i3/config"
	local LINE="bindsym $KEYBIND exec $CMD"
	if grep -q "$CMD" "$CFG" 2>/dev/null; then
		echo "Hotkey already exists for i3."
	else
		echo "$LINE" >>"$CFG"
		echo "Hotkey '$KEYBIND' added to i3 config."
	fi
}

setup_pantheon() {
	echo "Pantheon does not support programmatic shortcut setup."
	echo "Please set '$KEYBIND' for '$CMD' manually in System Settings → Keyboard → Shortcuts."
}

# Dispatcher
case "$DE" in
	*gnome* | *unity* | *budgie*) setup_gnome ;;
	*cinnamon*) setup_cinnamon ;;
	*mate*) setup_mate ;;
	*plasma* | *kde*) setup_kde ;;
	*xfce*) setup_xfce ;;
	*lxqt*) setup_lxqt ;;
	*lxde*) setup_lxde ;;
	*i3*) setup_i3 ;;
	*pantheon*) setup_pantheon ;;
	*)
		echo "Unsupported or undetected desktop: $DE"
		echo "You can manually assign '$KEYBIND' to '$CMD' in your desktop’s keyboard settings."
		;;
esac

echo "Done."
