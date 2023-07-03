#!/bin/bash

# regex
# Exercise doesn't exist. Made something up.


read -p "Enter email: " email

# Anyone even do this anymore
# keep it simple (don't make the user mad)
if [[ "$email" =~ ^.+@.+$ ]]; then
    echo "valid email"
else
    echo "invalid email"
fi
