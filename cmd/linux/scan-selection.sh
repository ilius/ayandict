#!/usr/bin/env bash
set -o errexit -o pipefail

# Get clipboard or primary text regardless of Wayland/X11
clip_get() {
    if [ -n "$WAYLAND_DISPLAY" ]; then
        if ! command -v wl-paste >/dev/null 2>&1; then
            echo "wl-paste command not found. Make sure wl-clipboard package is installed" >&2
            return 1
        fi
        wl-paste --primary
    elif [ -n "$DISPLAY" ]; then
        if ! command -v xclip >/dev/null 2>&1; then
            echo "xclip command not found. Make sure xclip package is installed" >&2
            return 1
        fi
        xclip -selection primary -o
    else
        echo "Unknown display server." >&2
        return 1
    fi
}

set -x

QUERY=$(clip_get)
SOCKET="/tmp/ayandict-$UID"
if command -v socat >/dev/null 2>&1 ; then
    printf 'scanpopup:%s' "$QUERY" | socat - "UNIX-CONNECT:$SOCKET"
else
    printf 'scanpopup:%s' "$QUERY" | nc -U "$SOCKET" -q 2
fi
