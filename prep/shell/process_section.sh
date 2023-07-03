#!/bin/bash

IFS=$'\n' read -d '' -r -a lines < sec.txt


# Starting at second line process all but last line
i=1
while [ $i -lt $((${#lines[@]} - 1)) ]; do
    # Current section is previous sections end
    lne=$(echo ${lines[i]} | cut -d '~' -f1)
    pi=$(( $i - 1 ))
    prev=${lines[$pi]}
    plne=$(echo $prev | cut -d '~' -f1)
    cutpln=$(echo $prev | cut -d '~' -f2)
    # new line with start and end
    echo "$plne~$lne~$cutpln"
    i=$(( $i + 1 ))
done

# For last line make limit something reasonably large
lastln=${lines[$i]}
lastlns=$(echo $lastln | cut -d '~' -f1)
lastlncut=$(echo $lastln | cut -d '~' -f2)
lastlne=$((lastlns + 1000))
echo "$lastlns~$lastlne~$lastlncut"






