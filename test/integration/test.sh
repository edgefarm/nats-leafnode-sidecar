#!/bin/bash

source './assert.sh'
source './test_lib.sh'

function main {

    prepare

    # source './test1.sh'
    # test1_do
    # sleep 1
    source './test2.sh'
    test2_do
    sleep 1
    cleanup
    [ $TEST_FAILED -eq 0 ] && echo -e "\e[32m\nAll tests passed\e[39m" && exit 0
    [ $TEST_FAILED -eq 1 ] && echo -e "\e[31m\nSome tests failed\e[39m" && exit 1
}

main
