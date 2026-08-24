#!/usr/bin/env bash
# Capture selected text on macOS and send it to AyanDict
# via UNIX socket for Scan Popup

set -o errexit -o pipefail -o nounset

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

# Wait for clipboard to update (up to 1 s)
for _ in {1..20}; do
    sleep 0.05
    SEL=$(pbpaste)
    [[ -n "$SEL" ]] && break
done

# Restore old clipboard if nothing else was copied manually
CURRENT=$(pbpaste)
[[ "$CURRENT" == "$SEL" ]] && printf "%s" "$OLD_CLIP" | pbcopy

# Trim whitespace
SEL=$(echo "$SEL" | sed 's/^[[:space:]]*//;s/[[:space:]]*$//')

# Check selection is non-empty
if [[ -z "$SEL" ]]; then
    echo "No text selected or empty clipboard" >&2
    exit 1
fi

# Send to AyanDict via socket (with timeout)
printf -- "%s" "scanpopup:$SEL" | nc -w 1 -U "$SOCKET"

echo "Sent to AyanDict: ${#SEL} chars"
