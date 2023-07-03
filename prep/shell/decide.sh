#!/bin/bash

NAME="John"
if [ "$NAME" = "John" ]; then
    echo "True - my name is indeed John"
fi

NAME="Bill"
if [ "$NAME" = "John" ]; then
    echo "Ya name John"
else
    echo "False"
    echo "You must have mistake me for $NAME"
fi

NAME="George"
if [ "$NAME" = "John" ]; then
    echo "John Lennon"
elif [ "$NAME" = "George" ]; then
    echo "George Harrison"
else
    echo "This leaves us with Paul and Ringo"
fi

#numerics
#comparison    Evaluated to true when
#$a -lt $b    $a < $b
#$a -gt $b    $a > $b
#$a -le $b    $a <= $b
#$a -ge $b    $a >= $b
#$a -eq $b    $a is equal to $b
#$a -ne $b    $a is not equal to $b

a=1
b=10
if [ $a -lt $b ]; then
    echo "ya a less than be"
fi

##strcmp
#comparison    Evaluated to true when
#"$a" = "$b"     $a is the same as $b
#"$a" == "$b"    $a is the same as $b
#"$a" != "$b"    $a is different from $b
#-z "$a"         $a is empty

A=""
if [ -z $A ]; then
    echo "A is nada"
fi


if [[ -z $A && $a -lt $b ]]; then
    echo "A is nada and a less than b"
fi

val=3
case $val in
    "1")
        echo "val is 1"
        ;;
    "2")
        echo "val is 2"
        ;;
    "3")
        echo "val is 3"
        ;;
    "4")
        echo "val is 4"
        ;;
esac




