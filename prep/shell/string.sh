#!/bin/bash

STRING="This is a string"
echo ${#STRING}


SUBSTRING="is"
expr index "$STRING" "$SUBSTRING"

STRING="extract the dog from this string"
POS=12
LEN=3
echo ${STRING:$POS:$LEN}


echo ${STRING:4}

echo "------\n"
#extraction
DATARECORD="last=Clifford,first=Johnny Boy,state=CA"
COMMA1=`expr index "$DATARECORD" ','`  # 14 position of first comma
CHOP1FIELD=${DATARECORD:$COMMA1}
COMMA2=`expr index "$CHOP1FIELD" ','`
LENGTH=`expr $COMMA2 - 6 - 1`
FIRSTNAME=${CHOP1FIELD:6:$LENGTH}
echo $FIRSTNAME

#Substring replace
STRING="to be or not to be"
echo ${STRING[@]/be/eat}
# Replace all
echo ${STRING[@]//be/eat}
# Replace occurance of substring if at beginning of $STRING
echo ${STRING[@]/#to be/eat now}
# Replace if at end
echo ${STRING[@]/%be/eat}
# Replace with shell output
echo ${STRING[@]/%be/be on $(date +%Y-%m-%d)}


##https://stackoverflow.com/questions/21688553/bash-expr-index-command
## s="Info.out.2014-02-08:INFO|SID:sXfzRjbmKbwX7jyaW1sog7n|Browser[Mozilla/5.0 (Windows NT 6.1; WOW64; rv:26.0) Gecko/20100101 Firefox/26.0]"
#
## strip everything after the first instance of 'Mozilla'
#prefix=${s%%Mozilla*}
#
## count number of characters in the string
#index=${#prefix}
#
## ...and show the result...
#echo "$index"
