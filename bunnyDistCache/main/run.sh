#!/bin/bash
trap "rm main;kill 0" EXIT

go build -o main
./main &

sleep 2
echo ">>> start test"
curl "http://localhost:9999/api?key=Tom" &
curl "http://localhost:9999/api?key=Jack" &
curl "http://localhost:9999/api?key=Sam" &
curl "http://localhost:9999/api?key=Tom" &
curl "http://localhost:9999/api?key=Jack" &
curl "http://localhost:9999/api?key=Sam" &
curl "http://localhost:9999/api?key=Tom" &
curl "http://localhost:9999/api?key=Jack" &
curl "http://localhost:9999/api?key=Sam" &

wait