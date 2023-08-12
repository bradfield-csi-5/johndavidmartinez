#!/bin/bash

cc -g -Wall -arch x86_64 -c -o call_sum_to_n.o call_sum_to_n.c
nasm -g -f macho64 --prefix _ -o sum_to_n.o sum_to_n.asm
cc -g -Wall -arch x86_64 -o sum_to_n call_sum_to_n.o sum_to_n.o
