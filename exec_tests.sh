#!/bin/sh

#
# Execute verbose testing procedures for Go code and colorize the output
#


mage vtest \
    | sed ''/PASS/s//$(printf "\033[0;32mPASS\033[0m")/'' \
    | sed ''/FAIL/s//$(printf "\033[0;31mFAIL\033[0m")/'' \
    | sed ''/SKIP/s//$(printf "\033[0;35mSKIP\033[0m")/''
