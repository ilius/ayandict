#!/bin/bash

echo "Opening AyanDict main window with query: $@"
echo "mainquery:$@" | socat -t 5 - UNIX-CONNECT:/tmp/ayandict-$UID