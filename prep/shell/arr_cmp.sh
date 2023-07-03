#!/bin/bash

arr=(23 45 34 1 2 3)
#echo ${arr[2]}
#
#echo ${arr[@]}
#
#echo ${#arr[@]}



a=(3 5 8 85 10 6)
b=(6 5 4 12 85)
c=(85 14 7 5 7)

i=0
k=0
j=0

i=0
while [ $i -lt ${#a[@]} ]; do
    k=0
    while [ $k -lt ${#b[@]} ]; do
        j=0
        while [ $j -lt ${#c[@]} ]; do
            #echo "cmp ${a[i]} ${b[k]} ${c[j]}"
            if [[ ${a[i]} -eq ${b[k]} && ${b[k]} -eq ${c[j]} ]]; then
                echo "common ${a[i]}"
            fi
            j=$((j+1))
        done
        k=$((k+1))
    done
    i=$((i + 1))
done

