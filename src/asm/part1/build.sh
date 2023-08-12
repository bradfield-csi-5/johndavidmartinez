#!/bin/bash

# Assemble
nasm -fmacho64 $1

# Link it
obj="$(echo $1 | cut -d "." -f1).o"
ld -lSystem -L/Library/Developer/CommandLineTools/SDKs/MacOSX.sdk/usr/lib $obj

# Run it
./a.out
