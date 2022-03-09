#!/bin/bash
source './assert.sh'
source './test_lib.sh'

function add_credsfile1 {
    path=${1} # where the raw file is stored
    symlink=${2} # where the symlink is stored
    jwt=${3}
    nkey=${4}
    echo "-----BEGIN NATS USER JWT-----" > ${path}
    echo ${jwt} >> ${path}
    echo "-----END NATS USER JWT-----" >> ${path}
    echo "************************* IMPORTANT *************************" >> ${path}
    echo "" >> ${path}
    echo "NKEY Seed printed below can be used to sign and prove identity." >> ${path}
    echo "NKEYs are sensitive and should be treated as secrets." >> ${path}
    echo "" >> ${path}
    echo "-----BEGIN NATS RSA PRIVATE KEY-----" >> ${path}
    echo ${nkey} >> ${path}
    echo "-----END NATS RSA PRIVATE KEY-----" >> ${path}
    echo "" >> ${path}
    echo "*************************************************************" >> ${path}

    ln -s ${path} ${symlink}
}

function prepare_test1 {
    TEST_FAILED=0
    mkdir -p ${1}/nats-credentials/..mytimeanddate/
    ln -s ${1}/nats-credentials/..mytimeanddate ${1}/nats-credentials/..data

    mkdir ${TMP_DIR}/config && chmod 777 -R ${TMP_DIR}/config
    cat <<EOF > ${TMP_DIR}/config/nats.json
{
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

function check_network_files_test1  {
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
    ok_message "check_network_files_test1"
}

function test1 {
    echo -e "\e[7mRunning test1...\e[27m"
    TMP_DIR=$(mktemp -d)
    # TMP_DIR=/tmp/sidecar
    prepare_test1 ${TMP_DIR}

    jwt1="myjwt1"
    nkey1="mynkey1"
    add_credsfile1 ${TMP_DIR}/nats-credentials/..data/secret ${TMP_DIR}/nats-credentials/mynetwork ${jwt1} ${nkey1}

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
    assert_contain `cat ${TMP_DIR}/config/nats.json | jq -r '.leafnodes.remotes[0].credentials'` "creds/mynetwork.creds" "remote credentials not equal"
    [[ $? -eq 0 ]] || TEST_FAILED=1
    check_network_files_test1 ${TMP_DIR} mynetwork ${jwt1} ${nkey1}
    sleep 1
    #read  -n 1 -p "Enter for continue to add second credential" mainmenuinput # for debug
    jwt2="myjwt2"
    nkey2="mynkey2"

    ## Add another credentials file
    add_credsfile1 ${TMP_DIR}/nats-credentials/..data/secret2 ${TMP_DIR}/nats-credentials/mynetwork2 ${jwt2} ${nkey2}
    sleep 1
    #read  -n 1 -p "Enter for continue to registering second credential" mainmenuinput # for debug
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
    check_network_files_test1 ${TMP_DIR} mynetwork ${jwy1} ${nkey1}
    check_network_files_test1 ${TMP_DIR} mynetwork2 ${jwy2} ${nkey2}

    ## Remove first credentials file
    #read  -n 1 -p "Enter for continue to end the test" mainmenuinput # for debug
    rm ${TMP_DIR}/nats-credentials/mynetwork
    sleep 1
    cat ${TMP_DIR}/config/nats.json
    assert_eq `cat ${TMP_DIR}/config/nats.json | jq -r '.http'` "8222" "http port not equal"
    [[ $? -eq 0 ]] || TEST_FAILED=1
    assert_eq `cat ${TMP_DIR}/config/nats.json | jq -r '.leafnodes.remotes[0].url'` "tls://connect.ngs.global:7422" "remote url not equal"
    [[ $? -eq 0 ]] || TEST_FAILED=1
    assert_contain `cat ${TMP_DIR}/config/nats.json | jq -r '.leafnodes.remotes[0].credentials'` "creds/mynetwork2.creds" "remote credentials not equal"
    [[ $? -eq 0 ]] || TEST_FAILED=1
    check_network_files_test1 ${TMP_DIR} mynetwork2 ${jwy2} ${nkey2}
    #read  -n 1 -p "Enter for continue to end the test" mainmenuinput # for debug
    check_test1_after_runtime $client_pid $registry_pid
    test_status $TEST_FAILED "Test1" ${TMP_DIR}
    cleanup_test1 ${TMP_DIR} $client_pid $registry_pid
}

function main {
    prepare
    test1
    sleep 1
    cleanup
}

main
