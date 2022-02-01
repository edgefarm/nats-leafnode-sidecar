NATS_SERVER_NAME='nats-leafnode-sidecar-test-nats'
NATS_PORT=0
RED='\e[31m'
GREEN='\e[32m'
NC='\e[39m' # No Color
TEST_LOG_FILE='test.log'
TEST_FAILED=0

netstat_used_local_ports()
{
    netstat -atn | awk '{printf "%s\n%s\n", $4, $4}' | grep -oE '[0-9]*$' | sort -n | uniq
}

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
