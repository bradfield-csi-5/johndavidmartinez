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
subsection=${3^}

# Modes of reading
mode="subsection"
if [ -z "$subsection" ]; then
    mode="section"
fi

if [ -z "$section" ]; then
    mode="main"
fi

echo "debug: article, $article"
echo "debug: section, $section"
echo "debug: subsection, $subsection"
echo "debug: mode, $mode"

function usage {
    echo "Usage: wiki [article] [section] (optional) [subsection] (optional)"
    echo "Must provide a section and an article"
    echo "Example: wiki walrus anatomy"
}

# Get Wikipedia site with curl and store in tempfile
WIKI_URL="https://en.wikipedia.org/wiki/$article"
SITE_TMP=$(mktemp)
curl $WIKI_URL > $SITE_TMP 2>/dev/null

function determine_section_endline {
    OUT_TMP=$(mktemp)
    IFS=$'\n' read -d '' -r -a lines < $LN_SECTIONS_START_TMP

    # Starting at second line process all but last line
    local i=1
    while [ $i -lt $((${#lines[@]} - 1)) ]; do
        # Current section is previous sections end
        local lne=$(echo ${lines[i]} | cut -d '~' -f1)
        local pi=$(( $i - 1 ))
        local prev=${lines[$pi]}
        local plne=$(echo $prev | cut -d '~' -f1)
        cutpln=$(echo $prev | cut -d '~' -f2)
        # new line with start and end
        echo "$plne~$lne~$cutpln" >> $OUT_TMP
        local i=$(( $i + 1 ))
    done
    
    # For last line make limit something reasonably large
    local lastln=${lines[$i]}
    local lastlns=$(echo $lastln | cut -d '~' -f1)
    local lastlncut=$(echo $lastln | cut -d '~' -f2)
    local lastlne=$((lastlns + 1000))
    echo "$lastlns~$lastlne~$lastlncut" >> $OUT_TMP
    echo $OUT_TMP
}
LN_SECTIONS_TMP=$(mktemp)
cat $SITE_TMP | awk '/<h2>/ { print NR "~" $0 }' > $LN_SECTIONS_TMP
if [[ ! -z $section ]]; then
    LN_SECTIONS_START_AND_END_TMP=$(determine_section_endline)
    sec=$(cat $LN_SECTIONS_START_AND_END_TMP | grep $section)
    sec_lns=$(echo $sec | cut -d "~" -f1)
    sec_lne=$(echo $sec | cut -d "~" -f2)
    LN_SUBSECTION_TMP=$(mktemp)
    cat $SITE_TMP | awk -v lns=$sec_lns lne=$sec_lne '/<h3>/&&NR<lne&&NR>lns { print $0 }' > $LN_SUBSECTION_TMP
fi
# Now we have
# LN_SECTIONS_TMP
# which each section and its start line
# LN_SUBSECTION_TMP
# the subsections of the current section and its start line

# Set main
# main line end is the line of the first section
# main line's summary is the first paragraph found after its start line

# 3 branches
# main
# section
# subsection

function print_section_summary {
    section_ln=$(echo $ln_sections | grep "$section")
    echo "debug: section_ln, $section_ln"
    section_ln_start=$(echo $ln_sections | grep $section | cut -d '~' -f1)
    echo "debug: section_ln_start, $section_ln_start"
    section_sentence=$(cat $SITE_TMP | awk -v lns=$section_ln_start '/<p>/&&NR>lns { print $0 }' | head -n 1)
    echo $section_sentence
}

function print_main {
    mainln_end=$(cat $LN_SECTIONS_TMP | head -n1 | cut -d '~' -f1)
    main_sentence=$(cat $SITE_TMP | awk -v lne=$mainln_end '/<p>/&&NR<lne { print $0 }')
    echo
    echo "Summary:"
    echo $main_sentence
    echo
    echo "Sections:"
    cat $LN_SECTIONS_TMP
}

function print_section {
    sectionlns=$(cat $LN_SECTION_TMP | grep $section | cut -d '~' -f1)
    section_sentence=$(cat $SITE_TMP | awk -v lns=$sectionlns '/<p>&&NR>lns { print $0 }')
    echo
    echo "Summary:"
    echo $section_sentence
    echo
    echo "SubSections:"
    cat $LN_SUBSECTIONS_TMP
}

function print_subsection {
    echo "junk"
}

case $mode in
    "main")
        print_main
        ;;
    "section")
        print_section
        ;;
    "subsection")
        print_subsection
        ;;
esac
