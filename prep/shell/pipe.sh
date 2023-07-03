#!/bin/bash

# print processors on system
P_COUNT=$(cat /proc/cpuinfo | grep processor | wc -l)
echo "Processor count: $P_COUNT"
