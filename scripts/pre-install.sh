#!/bin/bash

HOST="127.0.0.1"
PORT="8545"
URL=$HOST:$PORT

isRunning=$(curl -L {$URL} -o /dev/null -w '%{http_code}\n' -s)
if [[ $isRunning == 200 || $isRunning == 403 ]]
then
    # Service is online
    echo "Geth is running"
    exit 0
else
    # Service is offline or not working correctly
    echo "Geth is not running or not working correctly"
    exit 1
fi
