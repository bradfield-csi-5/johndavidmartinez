#!/bin/bash

NAMES=(Joe Jenny Sara Tony)
for N in ${NAMES[@]} ; do
    echo "My name is $N"
done

# loop on command output
for f in $( ls ); do
    echo "File is $f"
done

COUNT=4
while [ $COUNT -gt 0 ]; do
    echo "count: $COUNT"
    COUNT=$(($COUNT - 1))
done


# also until


COUNT=4
while [ $COUNT -gt 0 ]; do
    echo "iteration $COUNT"
    if [ $COUNT -eq 2 ]; then
        echo "Count is tew"
        break
    else
        COUNT=$(($COUNT - 1))
        continue
    fi
done



