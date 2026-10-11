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
TEST_SERVER="192.168.139.53"
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

if [ "$SEND" == "on" ]; then
	echo "Sending 192.168.139.42"
	sshpass -p "a" scp /usr/local/fpanel/core/bin/$TOOLNAME root@192.168.139.42:/usr/local/fpanel/core/bin/
#	echo "Sending 192.168.139.53"
#	sshpass -p "a" scp /usr/local/fpanel/core/bin/$TOOLNAME root@192.168.139.53:/usr/local/fpanel/core/bin/
#	echo "Sending 192.168.139.62"
#	sshpass -p "a" scp /usr/local/fpanel/core/bin/$TOOLNAME root@192.168.139.62:/usr/local/fpanel/core/bin/
fi
