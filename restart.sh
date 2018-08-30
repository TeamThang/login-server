#!/bin/sh

if [ $1 -eq 1 ]
then
    export GOPATH=`pwd`
    go build -o login_server src/server/main.go
    echo compile login_server finish !
fi

ps -ef | grep login_server | grep yzy | grep -v grep | awk -F' ' '{print $2}' | xargs kill -9
sleep 2s
nohup ./login_server &
echo restart login_server finish !