#!/bin/bash
GOPATH="$HOME/go"
export GOROOT=/usr/local/src/go
export GOPATH=$GOPATH:/usr/local/src/sources/fpanel/tools
GO11=/usr/local/src/go/bin/go
TOOLNAME=$1
$GO11 clean
if [ -z "$TOOLNAME" ]; then
  echo "Please choose a tool name"
  exit 0
fi
env GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc-win32 CGO_ENABLED=1 \
CGO_CFLAGS="-I/usr/include -Wno-error" $GO11 build -o $TOOLNAME $TOOLNAME.go
