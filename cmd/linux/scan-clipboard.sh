#!/bin/bash
set -o errexit -o pipefail

# Get clipboard or primary text regardless of Wayland/X11
clip_get() {
    if [ -n "$WAYLAND_DISPLAY" ]; then
        if ! command -v wl-paste >/dev/null 2>&1; then
            echo "wl-paste command not found. Make sure wl-clipboard package is installed" >&2
            return 1
        fi
        wl-paste
    elif [ -n "$DISPLAY" ]; then
        if ! command -v xclip >/dev/null 2>&1; then
            echo "xclip command not found. Make sure xclip package is installed" >&2
            return 1
        fi
        xclip -selection clipboard -o
    else
        echo "Unknown display server." >&2
        return 1
    fi
}

set -x

if which socat 2>/dev/null ; then
    echo scanpopup:$(clip_get) | socat - UNIX-CONNECT:/tmp/ayandict-$UID
else
    echo scanpopup:$(clip_get) | nc -U /tmp/ayandict-$UID -q 2
fi
