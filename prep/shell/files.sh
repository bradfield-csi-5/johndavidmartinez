#!/bin/bash

filename="sample.md"
if [ -e "$filename" ]; then
    echo "$filename exists as a file!"
else
    echo "$filename doesn\'t exist. Creating"
    touch $filename
fi

dir="secrets"
if [ -d $dir ]; then
    echo "$dir exists!"
else
    echo "$dir doesn't exist; creating"
    mkdir $dir
fi

if [ -r "$filename" ]; then
    echo "I can read $filename"
    cat $filename
else
    echo "can't read $filename"
fi
