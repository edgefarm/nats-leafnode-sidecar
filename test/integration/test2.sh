#!/bin/bash

source './common.sh'

# uncomment this if you want to step manually through the test
# MANUAL_STEPPING="1"

function test2_check_test_after_runtime() {
    local err=0
    if [ -z "$(ps -A | grep ${1})" ]; then
        error_message "pid for client1 not found"
        err=1
    fi

    if [ -z "$(ps -A | grep ${2})" ]; then
        error_message "pid for client2 not found"
        err=1
    fi

    if [ -z "$(ps -A | grep ${3})" ]; then
        error_message "pid for registry not found"
        err=1
    fi

    if [ $err -eq 1 ]; then
        error_message "test2_check_test_after_runtime failed"
        TEST_FAILED=1
    fi
}

function test2_cleanup_clients {
    kill ${2} &> ${1}/${TEST_LOG_FILE}
    kill ${3} &> ${1}/${TEST_LOG_FILE}
}

function test2_cleanup_registry {
    kill ${2} &> ${1}/${TEST_LOG_FILE}
}

function test2_cleanup_tmpdir {
    rm -rf ${1}
}

function test2_do {
    echo -e "\e[7mRunning test2...\e[27m\n"
    TMP_DIR=$(mktemp -d)
    # TMP_DIR=/tmp/sidecar
    mkdir -p ${TMP_DIR}/client1
    mkdir -p ${TMP_DIR}/client2
    prepare_nats_config ${TMP_DIR}
    prepare_volume ${TMP_DIR}/client1
    prepare_volume ${TMP_DIR}/client2

    jwt1="myjwt1"
    nkey1="mynkey1"
    add_credsfile ${TMP_DIR}/client1/nats-credentials/..data/secret ${TMP_DIR}/client1/nats-credentials/mynetwork ${jwt1} ${nkey1}
    add_credsfile ${TMP_DIR}/client2/nats-credentials/..data/secret ${TMP_DIR}/client2/nats-credentials/mynetwork ${jwt1} ${nkey1}


    rm nats-leafnode-sidecar-registry nats-leafnode-sidecar-client >& /dev/null
    go build -o ${TMP_DIR}/nats-leafnode-sidecar-client ../../cmd/client/main.go

    go build -o ${TMP_DIR}/nats-leafnode-sidecar-registry ../../cmd/registry/main.go
    ${TMP_DIR}/nats-leafnode-sidecar-registry --natsconfig ${TMP_DIR}/config/nats.json --creds ${TMP_DIR}/creds --natsuri nats://127.0.0.1:${NATS_PORT} --state ${TMP_DIR}/state.json & # &> ${TMP_DIR}/${TEST_LOG_FILE} &
    registry_pid=$!
    sleep 1

    echo -e "\n\e[7mStarting client 1 with 'mynetwork'\e[27m\n"
    ${TMP_DIR}/nats-leafnode-sidecar-client --creds ${TMP_DIR}/client1/nats-credentials --natsuri nats://127.0.0.1:${NATS_PORT} --component client1 & # &> ${TMP_DIR}/${TEST_LOG_FILE} &
    client1_pid=$!
    sleep 1

    echo -e "\n\e[7mStarting client 2 with 'mynetwork'\e[27m\n"
    # dlv debug ../../cmd/client/main.go --listen=0.0.0.0:2345 --api-version=2 --output /tmp/__debug_bin --headless --build-flags="-mod=vendor" --  --creds ${TMP_DIR}/client2/nats-credentials --natsuri nats://127.0.0.1:${NATS_PORT} --component client2 & # &> ${TMP_DIR}/${TEST_LOG_FILE} &
    ${TMP_DIR}/nats-leafnode-sidecar-client --creds ${TMP_DIR}/client2/nats-credentials --natsuri nats://127.0.0.1:${NATS_PORT} --component client2 & # &> ${TMP_DIR}/${TEST_LOG_FILE} &
    client2_pid=$!
    sleep 1
    echo client1 pid: ${client1_pid}
    echo client2 pid: ${client2_pid}


    cat ${TMP_DIR}/config/nats.json
    assert_eq `cat ${TMP_DIR}/config/nats.json | jq -r '.http'` "8222" "http port not equal"
    [[ $? -eq 0 ]] || TEST_FAILED=1
    assert_eq `cat ${TMP_DIR}/config/nats.json | jq -r '.leafnodes.remotes[0].url'` "tls://connect.ngs.global:7422" "remote url not equal"
    [[ $? -eq 0 ]] || TEST_FAILED=1
    assert_contain `cat ${TMP_DIR}/config/nats.json | jq -r '.leafnodes.remotes[0].credentials'` "creds/mynetwork.creds" "remote credentials not equal"
    [[ $? -eq 0 ]] || TEST_FAILED=1
    check_network_files ${TMP_DIR} mynetwork ${jwt1} ${nkey1}
    sleep 1


    echo -e "\n\e[7mAdding 'mynetwork2' to client2\e[27m\n"
    user_input
    jwt2="myjwt2"
    nkey2="mynkey2"
    add_credsfile ${TMP_DIR}/client2/nats-credentials/..data/secret2 ${TMP_DIR}/client2/nats-credentials/mynetwork2 ${jwt2} ${nkey2}
    sleep 1
    cat ${TMP_DIR}/config/nats.json
    assert_eq `cat ${TMP_DIR}/config/nats.json | jq -r '.http'` "8222" "http port not equal"
    [[ $? -eq 0 ]] || TEST_FAILED=1
    assert_eq `cat ${TMP_DIR}/config/nats.json | jq -r '.leafnodes.remotes[0].url'` "tls://connect.ngs.global:7422" "remote url not equal"
    [[ $? -eq 0 ]] || TEST_FAILED=1
    assert_contain `cat ${TMP_DIR}/config/nats.json | jq -r '.leafnodes.remotes[0].credentials'` "creds/mynetwork.creds" "remote credentials not equal"
    [[ $? -eq 0 ]] || TEST_FAILED=1
    assert_eq `cat ${TMP_DIR}/config/nats.json | jq -r '.leafnodes.remotes[1].url'` "tls://connect.ngs.global:7422" "remote url not equal"
    [[ $? -eq 0 ]] || TEST_FAILED=1
    assert_contain `cat ${TMP_DIR}/config/nats.json | jq -r '.leafnodes.remotes[1].credentials'` "creds/mynetwork2.creds" "remote credentials not equal"
    [[ $? -eq 0 ]] || TEST_FAILED=1
    check_network_files ${TMP_DIR} mynetwork ${jwy1} ${nkey1}
    check_network_files ${TMP_DIR} mynetwork2 ${jwt2} ${nkey2}

    echo -e "\n\e[7mRemoving 'mynetwork' for client2\e[27m\n"
    user_input
    rm ${TMP_DIR}/client2/nats-credentials/mynetwork
    sleep 1
    cat ${TMP_DIR}/config/nats.json
    assert_eq `cat ${TMP_DIR}/config/nats.json | jq -r '.http'` "8222" "http port not equal"
    [[ $? -eq 0 ]] || TEST_FAILED=1
    assert_eq `cat ${TMP_DIR}/config/nats.json | jq -r '.leafnodes.remotes[0].url'` "tls://connect.ngs.global:7422" "remote url not equal"
    [[ $? -eq 0 ]] || TEST_FAILED=1
    assert_contain `cat ${TMP_DIR}/config/nats.json | jq -r '.leafnodes.remotes[0].credentials'` "creds/mynetwork.creds" "remote credentials not equal"
    [[ $? -eq 0 ]] || TEST_FAILED=1
    assert_eq `cat ${TMP_DIR}/config/nats.json | jq -r '.leafnodes.remotes[1].url'` "tls://connect.ngs.global:7422" "remote url not equal"
    [[ $? -eq 0 ]] || TEST_FAILED=1
    assert_contain `cat ${TMP_DIR}/config/nats.json | jq -r '.leafnodes.remotes[1].credentials'` "creds/mynetwork2.creds" "remote credentials not equal"
    [[ $? -eq 0 ]] || TEST_FAILED=1
    check_network_files ${TMP_DIR} mynetwork ${jwy1} ${nkey1}
    check_network_files ${TMP_DIR} mynetwork2 ${jwt2} ${nkey2}

    echo -e "\n\e[7mRemoving 'mynetwork' client1\e[27m\n"
    user_input
    rm ${TMP_DIR}/client1/nats-credentials/mynetwork
    sleep 1
    cat ${TMP_DIR}/config/nats.json
    assert_eq `cat ${TMP_DIR}/config/nats.json | jq -r '.http'` "8222" "http port not equal"
    [[ $? -eq 0 ]] || TEST_FAILED=1
    assert_eq `cat ${TMP_DIR}/config/nats.json | jq -r '.leafnodes.remotes[0].url'` "tls://connect.ngs.global:7422" "remote url not equal"
    [[ $? -eq 0 ]] || TEST_FAILED=1
    assert_contain `cat ${TMP_DIR}/config/nats.json | jq -r '.leafnodes.remotes[0].credentials'` "creds/mynetwork2.creds" "remote credentials not equal"
    [[ $? -eq 0 ]] || TEST_FAILED=1
    check_network_files ${TMP_DIR} mynetwork2 ${jwt2} ${nkey2}

    user_input

    test2_check_test_after_runtime $client1_pid $client2_pid $registry_pid
    test_status $TEST_FAILED "test2" ${TMP_DIR}

    test2_cleanup_clients ${TMP_DIR} $client1_pid $client2_pid
    sleep 1
    test2_cleanup_registry ${TMP_DIR} $registry_pid
    sleep 1
    cat ${TMP_DIR}/config/nats.json
    assert_contain `cat ${TMP_DIR}/config/nats.json | jq -r '.leafnodes.remotes'` "[]" "remote credentials not equal"

    test2_cleanup_tmpdir ${TMP_DIR}
}
