#!/bin/bash

# doesn't exist lets make something up
# resource to expand if needed: https://www.shellscript.sh/examples/getopts/
# https://stackoverflow.com/questions/16483119/an-example-of-how-to-use-getopts-in-bash


# maybe some example usage of getops

function usage {
    echo "used it wrong bud"
}
# s & p arguments. z boolean flag
while getopts ":s:p:c" o; do
    case "${o}" in
        s)
            s=${OPTARG}
            ;;
        p)
            p=${OPTARG}
            ;;
        c)
            c="true"
            ;;
        *)
            usage
            ;;
    esac
done


echo "s = ${s}"
echo "p = ${p}"

if [ "$c" = "true" ]; then
    echo "c set"
else
    echo "c not set"
fi
