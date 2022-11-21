#!/bin/bash

START_COMMAND="go run cmd/server/main.go -port 8080"

# start store server
${START_COMMAND} &

# run e2e tests
go test -v ./test/e2e/...

# stop store server
_PID=`ps -ef | grep "${START_COMMAND}" | grep -v grep | awk '{ print $2 }'`
kill -15 ${_PID}
