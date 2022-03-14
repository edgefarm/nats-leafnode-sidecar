#!/bin/bash

source './assert.sh'
source './test_lib.sh'

function main {

    prepare

    source './test1.sh'
    test1_do
    sleep 1
    source './test2.sh'
    test2_do
    sleep 1
    cleanup
}

main
