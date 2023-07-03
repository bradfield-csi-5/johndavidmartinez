#!/bin/bash

# Your favorite wiki reader. Ayee!
# show first senetence
# walrus first sentence of main
# walrus taxonomy and evolution
# first sentence of taxonomy and evolution
# walrus anatomy "tusks and dentition" first sentence like "the most prominent feature of living specieis is its long tusks
#cat file | awk -F '<h2>|</h2>' '{print $8}' | awk NF

#SETUP
article=${1^}
section=${2^}

function usage {
    echo "Usage: wiki [article] [section]"
    echo "Must provide a section and an article"
    echo "Example: wiki walrus anatomy"
}

if [[ -z "$article" || -z "$section" ]]; then
    usage
fi

if [ "$#" -gt 2 ]; then
    echo "Too many arguments"
    echo
    usage
fi

WIKI_URL="https://en.wikipedia.org/wiki/$article"
echo $WIKI_URL

SITE_TMP=$(mktemp)
curl $WIKI_URL > $SITE_TMP 2>/dev/null
echo "debug: $SITE_TMP"

# Load site into memory
SITE_RAW=$(cat $SITE_TMP)

echo $section
SECTION_RAW=${SITE_RAW#*${id=$section}}
echo $SECTION_RAW

#<h2><span class="mw-headline" id="Anatomy">Anatomy</span></h2>












