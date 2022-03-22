#!/bin/bash

source './common.sh'

function test1_check_test_after_runtime() {
    local err=0
    if [ -z "$(ps -A | grep ${1})" ]; then
        error_message "pid for registry not found"
        err=1
    fi

    if [ -z "$(ps -A | grep ${2})" ]; then
        error_message "pid for client not found"
        err=1
    fi

    if [ $err -eq 1 ]; then
        error_message "test1_check_test_after_runtime failed"
        TEST_FAILED=1
    fi
}

function test1_cleanup() {
    kill ${2} &> ${1}/${TEST_LOG_FILE}
    sleep 0.5
    kill ${3} &> ${1}/${TEST_LOG_FILE}
    rm -rf ${1}
}

function check_network_files  {
    network=${2}
    jwt=${3}
    nkey=${4}

    if [ -f "${1}/creds/${network}.creds" ]; then
        if [ -z $(cat ${1}/creds/${network}.creds | grep "${jwt}") ]; then
            error_message "token not found in creds file"
            TEST_FAILED=1
        fi
        if [ -z $(cat ${1}/creds/${network}.creds | grep "${nkey}") ]; then
            error_message "nkey seed not found in creds file"
            TEST_FAILED=1
        fi

    else
        error_message "creds ${network}.creds file not found"
    fi
    ok_message "check_network_files"
}

function test1_do {
    echo -e "\e[7mRunning test1...\e[27m"
    TMP_DIR=$(mktemp -d)
    # TMP_DIR=/tmp/sidecar
    prepare_test ${TMP_DIR}

    jwt1="myjwt1"
    nkey1="mynkey1"
    accountPublicKey1="myaccountPublicKey1"
    add_credsfile ${TMP_DIR}/nats-credentials/..data/secret ${TMP_DIR}/nats-credentials/mynetwork ${jwt1} ${nkey1} ${accountPublicKey1}
    add_credsfile ${TMP_DIR}/nats-credentials/..data/edgefarm-sys ${TMP_DIR}/nats-credentials/edgefarm-sys ${jwt1} ${nkey1} ${accountPublicKey1}

    rm nats-leafnode-sidecar-registry-test1 nats-leafnode-sidecar-client-test1 >& /dev/null
    go build -o ${TMP_DIR}/nats-leafnode-sidecar-client-test1 ../../cmd/client/main.go

    #dlv debug ../../cmd/registry/main.go --listen=0.0.0.0:2345 --api-version=2 --output /tmp/__debug_bin --headless --build-flags="-mod=vendor" -- --natsconfig ${TMP_DIR}/config/nats.json --creds ${TMP_DIR}/creds --natsuri nats://127.0.0.1:${NATS_PORT} --state ${TMP_DIR}/state.json & # for debug
    #read  -n 1 -p "Enter for continue to end the test" mainmenuinput # for debug
    go build -o ${TMP_DIR}/nats-leafnode-sidecar-registry-test1 ../../cmd/registry/main.go
    ${TMP_DIR}/nats-leafnode-sidecar-registry-test1 --natsconfig ${TMP_DIR}/config/nats.json --creds ${TMP_DIR}/creds --natsuri nats://127.0.0.1:${NATS_PORT} --state ${TMP_DIR}/state.json & # &> ${TMP_DIR}/${TEST_LOG_FILE} &
    registry_pid=$!
    sleep 1

    #read  -n 1 -p "Enter for continue to start client" mainmenuinput # for debug
    ${TMP_DIR}/nats-leafnode-sidecar-client-test1 --creds ${TMP_DIR}/nats-credentials --natsuri nats://127.0.0.1:${NATS_PORT} --component mycomponent & # &> ${TMP_DIR}/${TEST_LOG_FILE} &
    client_pid=$!
    sleep 1

    cat ${TMP_DIR}/config/nats.json
    assert_eq `cat ${TMP_DIR}/config/nats.json | jq -r '.http'` "8222" "http port not equal"
    [[ $? -eq 0 ]] || TEST_FAILED=1
    assert_eq `cat ${TMP_DIR}/config/nats.json | jq -r '.leafnodes.remotes[0].url'` "tls://connect.ngs.global:7422" "remote url not equal"
    [[ $? -eq 0 ]] || TEST_FAILED=1
    assert_contain `cat ${TMP_DIR}/config/nats.json | jq -r '.leafnodes.remotes[0].credentials'` "creds/edgefarm-sys.creds" "remote credentials not equal"
    [[ $? -eq 0 ]] || TEST_FAILED=1
    assert_eq `cat ${TMP_DIR}/config/nats.json | jq -r '.leafnodes.remotes[1].url'` "tls://connect.ngs.global:7422" "remote url not equal"
    [[ $? -eq 0 ]] || TEST_FAILED=1
    assert_contain `cat ${TMP_DIR}/config/nats.json | jq -r '.leafnodes.remotes[1].credentials'` "creds/mynetwork.creds" "remote credentials not equal"
    [[ $? -eq 0 ]] || TEST_FAILED=1
    check_network_files ${TMP_DIR} mynetwork ${jwt1} ${nkey1}
    sleep 1
    #read  -n 1 -p "Enter for continue to add second credential" mainmenuinput # for debug
    jwt2="myjwt2"
    nkey2="mynkey2"

    ## Add another credentials file
    add_credsfile ${TMP_DIR}/nats-credentials/..data/secret2 ${TMP_DIR}/nats-credentials/mynetwork2 ${jwt2} ${nkey2} ${accountPublicKey1}
    sleep 1
    #read  -n 1 -p "Enter for continue to registering second credential" mainmenuinput # for debug
    cat ${TMP_DIR}/config/nats.json
    assert_eq `cat ${TMP_DIR}/config/nats.json | jq -r '.http'` "8222" "http port not equal"
    [[ $? -eq 0 ]] || TEST_FAILED=1
    assert_eq `cat ${TMP_DIR}/config/nats.json | jq -r '.leafnodes.remotes[0].url'` "tls://connect.ngs.global:7422" "remote url not equal"
    [[ $? -eq 0 ]] || TEST_FAILED=1
    assert_contain `cat ${TMP_DIR}/config/nats.json | jq -r '.leafnodes.remotes[0].credentials'` "creds/edgefarm-sys.creds" "remote credentials not equal"
    [[ $? -eq 0 ]] || TEST_FAILED=1
    assert_eq `cat ${TMP_DIR}/config/nats.json | jq -r '.leafnodes.remotes[1].url'` "tls://connect.ngs.global:7422" "remote url not equal"
    [[ $? -eq 0 ]] || TEST_FAILED=1
    assert_contain `cat ${TMP_DIR}/config/nats.json | jq -r '.leafnodes.remotes[1].credentials'` "creds/mynetwork.creds" "remote credentials not equal"
    [[ $? -eq 0 ]] || TEST_FAILED=1
    assert_eq `cat ${TMP_DIR}/config/nats.json | jq -r '.leafnodes.remotes[2].url'` "tls://connect.ngs.global:7422" "remote url not equal"
    [[ $? -eq 0 ]] || TEST_FAILED=1
    assert_contain `cat ${TMP_DIR}/config/nats.json | jq -r '.leafnodes.remotes[2].credentials'` "creds/mynetwork2.creds" "remote credentials not equal"
    [[ $? -eq 0 ]] || TEST_FAILED=1
    check_network_files ${TMP_DIR} mynetwork ${jwy1} ${nkey1}
    check_network_files ${TMP_DIR} mynetwork2 ${jwt2} ${nkey2}

    ## Remove first credentials file
    #read  -n 1 -p "Enter for continue to end the test" mainmenuinput # for debug
    rm ${TMP_DIR}/nats-credentials/mynetwork
    sleep 1
    cat ${TMP_DIR}/config/nats.json
    assert_eq `cat ${TMP_DIR}/config/nats.json | jq -r '.http'` "8222" "http port not equal"
    [[ $? -eq 0 ]] || TEST_FAILED=1
    assert_eq `cat ${TMP_DIR}/config/nats.json | jq -r '.leafnodes.remotes[0].url'` "tls://connect.ngs.global:7422" "remote url not equal"
    [[ $? -eq 0 ]] || TEST_FAILED=1
    assert_contain `cat ${TMP_DIR}/config/nats.json | jq -r '.leafnodes.remotes[0].credentials'` "creds/edgefarm-sys.creds" "remote credentials not equal"
    [[ $? -eq 0 ]] || TEST_FAILED=1
    assert_eq `cat ${TMP_DIR}/config/nats.json | jq -r '.leafnodes.remotes[1].url'` "tls://connect.ngs.global:7422" "remote url not equal"
    [[ $? -eq 0 ]] || TEST_FAILED=1
    assert_contain `cat ${TMP_DIR}/config/nats.json | jq -r '.leafnodes.remotes[1].credentials'` "creds/mynetwork2.creds" "remote credentials not equal"
    [[ $? -eq 0 ]] || TEST_FAILED=1
    check_network_files ${TMP_DIR} mynetwork2 ${jwt2} ${nkey2}
    #read  -n 1 -p "Enter for continue to end the test" mainmenuinput # for debug
    test1_check_test_after_runtime $client_pid $registry_pid
    test_status $TEST_FAILED "Test1" ${TMP_DIR}
    test1_cleanup ${TMP_DIR} $client_pid $registry_pid
}
