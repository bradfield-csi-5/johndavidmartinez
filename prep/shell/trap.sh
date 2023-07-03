#!/bin/bash

echo "GOTCHYA"

trap "BOOM!" SIGINT SIGTERM
echo "RUNNING"
echo "GOTCHYA GOTCHYA"

while true
do
    sleep 60
done

# useful for shutdown hooks
