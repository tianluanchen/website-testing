#!/usr/bin/env bash

#  You need to install nodejs and nodemon for it to work properly.
nodemon --exec "go run main.go --open=false" --ext "go,html,js,css,json" --delay 1s
