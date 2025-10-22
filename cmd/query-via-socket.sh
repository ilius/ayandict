#!/bin/bash

echo "query:fuzzy:$@" | socat -t 5 - UNIX-CONNECT:/tmp/ayandict-$UID,crnl
echo