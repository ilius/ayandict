#!/usr/bin/env bash
# Capture copied text on macOS and send it to AyanDict
# via UNIX socket for Scan Popup

set -o errexit -o pipefail -o nounset

SOCKET="/tmp/ayandict-$UID"

# Ensure the socket exists before trying
if [[ ! -S "$SOCKET" ]]; then
    echo "Socket not found: $SOCKET" >&2
    exit 1
fi

# Read clipboard
SEL=$(pbpaste)

# Trim whitespace
SEL=$(echo "$SEL" | sed 's/^[[:space:]]*//;s/[[:space:]]*$//')

# Check selection is non-empty
if [[ -z "$SEL" ]]; then
    echo "No text selected or empty clipboard" >&2
    exit 1
fi

# Send to AyanDict via socket
printf -- "%s" "scanpopup:$SEL" | nc -w 1 -U "$SOCKET"

echo "Sent to AyanDict: ${#SEL} chars"
