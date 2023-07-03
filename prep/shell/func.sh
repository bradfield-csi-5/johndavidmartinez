#!/bin/bash

#function adder {
#    echo "$(($1 + $2))"
#}
#
#ZIP=1
#ZIP=$(adder $ZIP 1)
#ZIP=$(adder $ZIP 1)
#ZIP=$(adder $ZIP 1)
#echo $ZIP

function ENGLISH_CAL {
    case $2 in
        "plus")
            echo $(($1 + $3))
            ;;
        "minus")
            echo $(($1 - $3))
            ;;
        "times")
            echo $(($1 * $3))
            ;;
    esac
}

VAL=$(ENGLISH_CAL 2 plus 2) #4
VAL=$(ENGLISH_CAL $VAL times 2) #8
VAL=$(ENGLISH_CAL $VAL minus 1) #7
echo $VAL
