#!/usr/bin/env bash

./build-keep-client-docker-img.sh

RESULT="$(./run-keep-client.sh)"
echo "RESULT: $RESULT"
SHA2_REGEX="[A-Fa-f0-9]{64}"
regex="$SHA2_REGEX"
if [[ $RESULT =~ $regex ]]; then
   echo "Pass"
else
   echo "Fail"
fi
