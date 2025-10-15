#!/bin/bash

echo statusicon:activate | socat - UNIX-CONNECT:/tmp/ayandict-$UID
