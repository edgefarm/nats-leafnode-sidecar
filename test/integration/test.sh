#!/bin/bash
source './assert.sh'
source './test_lib.sh'

function prepare_test1 {
    TEST_FAILED=0
    mkdir -p ${TMP_DIR}/nats-credentials

    cat <<EOF > ${1}/nats-credentials/secret
{"userAccountName":"myAccount","username":"myUser","password":"myPassword","creds":"-----BEGIN NATS USER JWT-----\nmyToken\n------END NATS USER JWT------\n\n************************* IMPORTANT *************************\nNKEY Seed printed below can be used to sign and prove identity.\nNKEYs are sensitive and should be treated as secrets.\n\n-----BEGIN USER NKEY SEED-----\nmyNkeySeed\n------END USER NKEY SEED------\n\n*************************************************************\n"}
EOF
    ln -s ${1}/nats-credentials/secret ${1}/nats-credentials/edgefarm.network-natsUserData

    mkdir ${TMP_DIR}/config && chmod 777 -R ${TMP_DIR}/config
    cat <<EOF > ${TMP_DIR}/config/nats.json
{
	"accounts": {
		"default": {
			"users": [{
				"user": "default",
				"password": ""
				}]
		}
	},
	"http": 8222,
	"leafnodes": {
		"remotes": []
	},
	"pid_file": "/var/run/nats.pid"
}
EOF
}

function check_test1_after_runtime() {
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
        error_message "check_test1_after_runtime failed"
        TEST_FAILED=1
    fi
}

function cleanup_test1() {
    kill ${2} &> ${1}/${TEST_LOG_FILE}
    sleep 0.5
    kill ${3} &> ${1}/${TEST_LOG_FILE}
    rm -rf ${1}
    TEST_FAILED=0
}

function check_accounts_file_test1 {
    if [ -f "${1}/creds/myAccount.creds" ]; then
        if [ -z $(cat ${1}/creds/myAccount.creds | grep "myToken") ]; then
            error_message "token not found in creds file"
            TEST_FAILED=1
        fi
        if [ -z $(cat ${1}/creds/myAccount.creds | grep "myNkeySeed") ]; then
            error_message "nkey seed not found in creds file"
            TEST_FAILED=1
        fi

    else
        error_message "creds file not found"
    fi
    ok_message "check_accounts_file_test1"
}

function test1 {
    echo -e "\e[7mRunning test1...\e[27m"
    TMP_DIR=$(mktemp -d)
    # TMP_DIR=/tmp/sidecar
    prepare_test1 ${TMP_DIR}

    rm nats-leafnode-sidecar-registry-test1 nats-leafnode-sidecar-client-test1 >& /dev/null
    go build -o ${TMP_DIR}/nats-leafnode-sidecar-registry-test1 ../../cmd/registry/main.go
    go build -o ${TMP_DIR}/nats-leafnode-sidecar-client-test1 ../../cmd/client/main.go

    ${TMP_DIR}/nats-leafnode-sidecar-registry-test1 --natsconfig ${TMP_DIR}/config/nats.json --creds ${TMP_DIR}/creds --natsuri nats://127.0.0.1:$NATS_PORT &> ${TMP_DIR}/${TEST_LOG_FILE} &
    registry_pid=$!
    sleep 1
    ${TMP_DIR}/nats-leafnode-sidecar-client-test1 --creds ${TMP_DIR}/nats-credentials --natsuri nats://127.0.0.1:$NATS_PORT &> ${TMP_DIR}/${TEST_LOG_FILE} &
    client_pid=$!
    sleep 1
    assert_eq `cat ${TMP_DIR}/config/nats.json | jq -r '.http'` "8222" "http port not equal"
    [[ $? -eq 0 ]] || TEST_FAILED=1
    assert_eq `cat ${TMP_DIR}/config/nats.json | jq -r '.accounts.myAccount.users[0].password'` myPassword "password not equal"
    [[ $? -eq 0 ]] || TEST_FAILED=1
    assert_eq `cat ${TMP_DIR}/config/nats.json | jq -r '.accounts.myAccount.users[0].user'` myUser "user not equal"
    [[ $? -eq 0 ]] || TEST_FAILED=1
    [ -z $(cat ${TMP_DIR}/config/nats.json | jq -r '.accounts.default.users[0].password') ] || TEST_FAILED=1
    assert_eq `cat ${TMP_DIR}/config/nats.json | jq -r '.accounts.default.users[0].user'` default "user not equal"
    [[ $? -eq 0 ]] || TEST_FAILED=1
    assert_eq `cat ${TMP_DIR}/config/nats.json | jq -r '.leafnodes.remotes[0].url'` "tls://connect.ngs.global:7422" "remote url not equal"
    [[ $? -eq 0 ]] || TEST_FAILED=1
    assert_contain `cat ${TMP_DIR}/config/nats.json | jq -r '.leafnodes.remotes[0].credentials'` "creds/myAccount.creds" "remote credentials not equal"
    [[ $? -eq 0 ]] || TEST_FAILED=1
    assert_eq `cat ${TMP_DIR}/config/nats.json | jq -r '.leafnodes.remotes[0].account'` "myAccount" "remote account not equal"
    [[ $? -eq 0 ]] || TEST_FAILED=1
    check_accounts_file_test1 ${TMP_DIR}
    sleep 0.5

    check_test1_after_runtime $client_pid $registry_pid
    test_status $TEST_FAILED "Test1" ${TMP_DIR}
    cleanup_test1 ${TMP_DIR} $client_pid $registry_pid
}


function prepare_test2 {
    TEST_FAILED=0
    mkdir -p ${TMP_DIR}/nats-credentials

    cat <<EOF > ${1}/nats-credentials/secret
{"userAccountName":"myAccount","username":"myUser","password":"myPassword","creds":"-----BEGIN NATS USER JWT-----\nmyToken\n------END NATS USER JWT------\n\n************************* IMPORTANT *************************\nNKEY Seed printed below can be used to sign and prove identity.\nNKEYs are sensitive and should be treated as secrets.\n\n-----BEGIN USER NKEY SEED-----\nmyNkeySeed\n------END USER NKEY SEED------\n\n*************************************************************\n"}
EOF
    ln -s ${1}/nats-credentials/secret ${1}/nats-credentials/edgefarm.network-natsUserData
    mkdir ${TMP_DIR}/config && chmod 777 -R ${TMP_DIR}/config
}

function check_test2_after_runtime() {
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
        error_message "check_test2_after_runtime failed"
        TEST_FAILED=1
    fi
}

function cleanup_test2() {
    kill ${2} &> ${1}/${TEST_LOG_FILE}
    sleep 0.5
    kill ${3} &> ${1}/${TEST_LOG_FILE}
    rm -rf ${1}
    TEST_FAILED=0
}

function check_accounts_file_test2 {
    if [ -f "${1}/creds/myAccount.creds" ]; then
        if [ -z $(cat ${1}/creds/myAccount.creds | grep "myToken") ]; then
            error_message "token not found in creds file"
            TEST_FAILED=1
        fi
        if [ -z $(cat ${1}/creds/myAccount.creds | grep "myNkeySeed") ]; then
            error_message "nkey seed not found in creds file"
            TEST_FAILED=1
        fi

    else
        error_message "creds file not found"
        TEST_FAILED=1
    fi
    ok_message "check_accounts_file_test2"
}

function test2 {
    echo -e "\e[7mRunning test2...\e[27m"
    TMP_DIR=$(mktemp -d)
    # TMP_DIR=/tmp/sidecar
    prepare_test2 ${TMP_DIR}

    rm nats-leafnode-sidecar-registry-test2 nats-leafnode-sidecar-client-test2 >& /dev/null
    go build -o ${TMP_DIR}/nats-leafnode-sidecar-registry-test2 ../../cmd/registry/main.go
    go build -o ${TMP_DIR}/nats-leafnode-sidecar-client-test2 ../../cmd/client/main.go

    ${TMP_DIR}/nats-leafnode-sidecar-registry-test2 --natsconfig ${TMP_DIR}/config/nats.json --creds ${TMP_DIR}/creds --natsuri nats://127.0.0.1:$NATS_PORT &> ${TMP_DIR}/${TEST_LOG_FILE} &
    registry_pid=$!
    sleep 1
    ${TMP_DIR}/nats-leafnode-sidecar-client-test2 --creds ${TMP_DIR}/nats-credentials --natsuri nats://127.0.0.1:$NATS_PORT &> ${TMP_DIR}/${TEST_LOG_FILE} &
    client_pid=$!
    sleep 1
    assert_eq `cat ${TMP_DIR}/config/nats.json | jq -r '.http'` "8222" "http port not equal"
    [[ $? -eq 0 ]] || TEST_FAILED=1
    assert_eq `cat ${TMP_DIR}/config/nats.json | jq -r '.accounts.myAccount.users[0].password'` myPassword "password not equal"
    [[ $? -eq 0 ]] || TEST_FAILED=1
    assert_eq `cat ${TMP_DIR}/config/nats.json | jq -r '.accounts.myAccount.users[0].user'` myUser "user not equal"
    [[ $? -eq 0 ]] || TEST_FAILED=1
    [ -z $(cat ${TMP_DIR}/config/nats.json | jq -r '.accounts.default.users[0].password') ] || TEST_FAILED=1
    assert_eq `cat ${TMP_DIR}/config/nats.json | jq -r '.accounts.default.users[0].user'` default "user not equal"
    [[ $? -eq 0 ]] || TEST_FAILED=1
    assert_eq `cat ${TMP_DIR}/config/nats.json | jq -r '.leafnodes.remotes[0].url'` "tls://connect.ngs.global:7422" "remote url not equal"
    [[ $? -eq 0 ]] || TEST_FAILED=1
    assert_contain `cat ${TMP_DIR}/config/nats.json | jq -r '.leafnodes.remotes[0].credentials'` "creds/myAccount.creds" "remote credentials not equal"
    [[ $? -eq 0 ]] || TEST_FAILED=1
    assert_eq `cat ${TMP_DIR}/config/nats.json | jq -r '.leafnodes.remotes[0].account'` "myAccount" "remote account not equal"
    [[ $? -eq 0 ]] || TEST_FAILED=1
    check_accounts_file_test2 ${TMP_DIR}
    sleep 0.5

    check_test2_after_runtime $client_pid $registry_pid
    test_status $TEST_FAILED "Test2" ${TMP_DIR}
    cleanup_test2 ${TMP_DIR} $client_pid $registry_pid
}

function main {
    prepare
    test1
    test2
    sleep 1
    cleanup
}

main
