#!/bin/bash
# Capture selected text on macOS and send it to AyanDict
# via UNIX socket for Scan Popup

SOCKET="/tmp/ayandict-$UID"

# Ensure the socket exists before trying
if [[ ! -S "$SOCKET" ]]; then
    echo "Socket not found: $SOCKET" >&2
    exit 1
fi

# Save clipboard to restore later
OLD_CLIP=$(pbpaste 2>/dev/null || echo "")

# Simulate Cmd+C (copy current selection)
osascript -e 'tell application "System Events" to keystroke "c" using command down'

# Wait briefly for clipboard to update
sleep 0.15

# Read new clipboard
SEL=$(pbpaste)

# Restore old clipboard content (optional)
echo -n "$OLD_CLIP" | pbcopy

# Trim whitespace
SEL=$(echo "$SEL" | sed 's/^[[:space:]]*//;s/[[:space:]]*$//')

# If selection is non-empty, send to socket
if [[ -n "$SEL" ]]; then
    printf "%s" "scanpopup:$SEL" | nc -U "$SOCKET"
else
    echo "No text selected or empty clipboard" >&2
fi
