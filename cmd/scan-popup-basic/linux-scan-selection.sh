# Get clipboard or primary text regardless of Wayland/X11
clip_get() {
    # usage: clip_get [clipboard|primary]
    local sel=${1:-clipboard}

    if [ -n "$WAYLAND_DISPLAY" ] && command -v wl-paste >/dev/null 2>&1; then
        if [ "$sel" = "primary" ]; then
            wl-paste --primary
        else
            wl-paste
        fi
    elif [ -n "$DISPLAY" ] && command -v xclip >/dev/null 2>&1; then
        if [ "$sel" = "primary" ]; then
            xclip -selection primary -o
        else
            xclip -selection clipboard -o
        fi
    else
        echo "No clipboard utility found or unknown display server." >&2
        return 1
    fi
}


DIR=$(dirname $0)
cd $DIR
./scan-popup-basic $(clip_get primary)