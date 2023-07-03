#!/bin/bash

PRICE_PER_APPLE=5
echo "The price of an Apple today is \$HK $PRICE_PER_APPLE"

MyFirstLetters=ABC
echo "The first 10 letters in the alphabet are ${MyFirstLetters}DEFGHIJ"

greeting='Hello     world!'
echo $greeting" now with spaces: $greeting"

whitespaceignored='Why      is the whitespace ignored'
echo $whitespaceignored

LS_RESULT=$(ls)
echo $LS_RESULT
