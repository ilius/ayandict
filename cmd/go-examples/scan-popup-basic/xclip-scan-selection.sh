#!/bin/bash

DIR=$(dirname $0)
cd $DIR
./scan-popup-basic $(xclip -o)
