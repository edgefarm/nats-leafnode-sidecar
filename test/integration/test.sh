#!/bin/bash
source './assert.sh'
NATS_SERVER_NAME='nats-leafnode-sidecar-test-nats'
NATS_PORT=0
RED='\e[31m'
GREEN='\e[32m'
NC='\e[39m' # No Color
TEST_LOG_FILE='test.log'
TEST_FAILED=0
# function available_port {
#     read LOWERPORT UPPERPORT < /proc/sys/net/ipv4/ip_local_port_range
#     cat /proc/sys/net/ipv4/ip_local_port_range
#     printf "%s" $LOWERPORT
#     echo $UPPERPORT
#     sleep 1000
#     while :
#     do
#         PORT="`shuf -i $LOWERPORT-$UPPERPORT -n 1`"
#         ss -lpn | grep -q ":$PORT " || break
#     done
#     echo $PORT
#     echo $PORT
# }

# -------------------------------------------------------------------
# netstat_used_local_ports: get used tcp-ports
netstat_used_local_ports()
{
    netstat -atn | awk '{printf "%s\n%s\n", $4, $4}' | grep -oE '[0-9]*$' | sort -n | uniq
}

# -------------------------------------------------------------------
# netstat_free_local_port: get one free tcp-port
available_port()
{
    read lowerPort upperPort < /proc/sys/net/ipv4/ip_local_port_range
    # create a local array of used ports
    local all_used_ports=($(netstat_used_local_ports))

    for port in $(seq $lowerPort $upperPort); do
        for used_port in "${all_used_ports[@]}"; do
            if [ $used_port -eq $port ]; then
                continue
            else
                echo $port
                return 0
            fi
        done
    done
}

function prepare {
    echo "Preparing test environment..."
    NATS_PORT=$(available_port)

    [[ $(docker ps -f "name=$NATS_SERVER_NAME" --format '{{.Names}}') == $NATS_SERVER_NAME ]] && \
    docker rm -f $NATS_SERVER_NAME >& /dev/null
    docker run -p $NATS_PORT:4222 --rm -d --name ${NATS_SERVER_NAME} nats:latest >& /dev/null
    sleep 1
}

function cleanup {
    echo "Cleaning up"
    [[ $(docker ps -f "name=$NATS_SERVER_NAME" --format '{{.Names}}') == $NATS_SERVER_NAME ]] && \
    docker kill ${NATS_SERVER_NAME} >& /dev/null
}

function error_message {
    echo -e "${RED}${1}: failed${NC}"
}

function ok_message {
    echo -e "${GREEN}${1}: passed${NC}"
}

function test_status {
    if [ $1 -eq 0 ]; then
        ok_message ${2}
    else
        error_message ${2}
        echo Logs:
        echo ------------------------------------------------
        cat ${3}/${TEST_LOG_FILE}
        echo ------------------------------------------------
    fi
}

function prepare_test1 {
    mkdir -p ${TMP_DIR}/nats-credentials

    cat <<EOF > ${1}/nats-credentials/secret
{"userAccountName":"myAccount","username":"myUser","password":"myPassword","creds":"-----BEGIN NATS USER JWT-----\nmyToken\n------END NATS USER JWT------\n\n************************* IMPORTANT *************************\nNKEY Seed printed below can be used to sign and prove identity.\nNKEYs are sensitive and should be treated as secrets.\n\n-----BEGIN USER NKEY SEED-----\nmyNkeySeed\n------END USER NKEY SEED------\n\n*************************************************************\n"}
EOF
    ln -s ${1}/nats-credentials/secret ${1}/nats-credentials/edgefarm.network-natsUserData

    mkdir ${TMP_DIR}/config && chmod 777 -R ${TMP_DIR}/config
    cat <<EOF > ${TMP_DIR}/config/nats.json
{
  "accounts": {},
  "http": 8222,
  "leafnodes": {
    "remotes": []
  },
  "pid_file": "/var/run/nats/nats.pid"
}
EOF
}

function check_test1_after_runtime() {
    local err=0
    [ -z "$(ps -A | grep ${1})" ] && error_message "pid for registry not found" && err=1
    [ -z "$(ps -A | grep ${2})" ] && error_message "pid for client not found" && err=1
    [ $err -eq 1 ] && TEST_FAILED=1
}

function cleanup_test1() {
    kill ${2} &> ${1}/${TEST_LOG_FILE}
    sleep 0.5
    kill ${3} &> ${1}/${TEST_LOG_FILE}
    rm -rf ${1}
    TEST_FAILED=0
}

function test1 {
    echo -e "\e[7mRunning test1...\e[27m"
    TMP_DIR=$(mktemp -d)
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
    assert_eq `cat ${TMP_DIR}/config/nats.json | jq -r '.leafnodes.remotes[0].url'` "tls://connect.ngs.global:7422" "remote url not equal"
    [[ $? -eq 0 ]] || TEST_FAILED=1
    assert_eq `cat ${TMP_DIR}/config/nats.json | jq -r '.leafnodes.remotes[0].credentials'` "myUser.creds" "remote credentials not equal"
    [[ $? -eq 0 ]] || TEST_FAILED=1
    assert_eq `cat ${TMP_DIR}/config/nats.json | jq -r '.leafnodes.remotes[0].account'` "myAccount" "remote account not equal"
    [[ $? -eq 0 ]] || TEST_FAILED=1

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
