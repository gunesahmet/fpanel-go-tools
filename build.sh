#!/bin/bash
RELEASES=`lsb_release -is`
GOPATH="$HOME/go"
export GOROOT=/usr/local/src/go1.18
export GOPATH=/usr/local/src/sources/fpanel/tools
GO11=/usr/local/src/go1.18/bin/go
TOOLNAME=$1
SEND="off"
echo "Building $TOOLNAME"
[[ ! -z "$2" ]] && SEND="on"
$GO11 clean
if [ -z "$TOOLNAME" ]; then
  echo "Please choose a tool name"
  exit 0
fi
export GO111MODULE=off
#echo "Compiling $TOOLNAME"

$GO11 build -ldflags "-s -w" -gcflags="all=-N -l" -o $TOOLNAME $TOOLNAME.go
RETVAL=$?
if [ "$RETVAL" == "0" ]
    then
        if [ "$TOOLNAME" == "logger" ]; then
          cp -rP $TOOLNAME /usr/local/fpanel/core/sbin/$TOOLNAME
          cp -rP $TOOLNAME /root/panel/packages/panel-core/usr/local/fpanel/core/sbin/$TOOLNAME
        else
          cp -rP $TOOLNAME /usr/local/fpanel/core/bin/$TOOLNAME
          cp -rP $TOOLNAME /root/panel/packages/panel-core/usr/local/fpanel/core/bin/$TOOLNAME
        fi
        rm -rf $TOOLNAME
fi


