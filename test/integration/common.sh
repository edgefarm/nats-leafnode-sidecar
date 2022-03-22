#!/bin/bash

function prepare_test {
    prepare_volume ${1}
    prepare_nats_config ${1}
}

function prepare_nats_config {
    mkdir ${1}/config && chmod 777 -R ${1}/config
    cat <<EOF > ${TMP_DIR}/config/nats.json
{
  "pid_file": "/var/run/nats/nats.pid",
  "http": 8222,
  "leafnodes": {
    "remotes": [
      {
        "url": "tls://connect.ngs.global:7422",
        "credentials": "/creds/edgefarm-sys.creds",
        "account": "MY-SYS-ACCOUNT-PUB-KEY"
      }
    ]
  },
  "operator": "MY-OPERATOR-JWY,
  "system_account": "MY-SYS-ACCOUNT-PUB-KEY",
  "resolver": {
    "type": "cache",
    "dir": "/jwt",
    "ttl": "1h",
    "timeout": "1.9s"
  },
  "resolver_preload": {
    "MY-SYS-ACCOUNT-PUB-KEY": "MY-SYS-ACCOUNT-JWT"
  }
}

EOF
}

function prepare_volume {
    mkdir -p ${1}/nats-credentials/..mytimeanddate/
    ln -s ${1}/nats-credentials/..mytimeanddate ${1}/nats-credentials/..data

    mkdir ${1}/config && chmod 777 -R ${1}/config
}

function add_credsfile {
    path=${1} # where the raw file is stored
    symlink=${2} # where the symlink is stored
    jwt=${3}
    nkey=${4}
    accountPublicKey=${5}
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
    echo ${accountPublicKey} > ${path}.pub
    ln -s ${path}.pub ${symlink}.pub
}


function check_network_files  {
    network=${2}
    jwt=${3}
    nkey=${4}

    if [ -f "${1}/creds/${network}.creds" ]; then
        if [ -z $(cat ${1}/creds/${network}.creds | grep "${jwt}") ]; then
            error_message "token not found in creds file ${1}/creds/${network}.creds"
            TEST_FAILED=1
        fi
        if [ -z $(cat ${1}/creds/${network}.creds | grep "${nkey}") ]; then
            error_message "nkey seed not found in creds file ${1}/creds/${network}.creds"
            TEST_FAILED=1
        fi
    else
        error_message "creds ${network}.creds file not found"
        TEST_FAILED=1
    fi
    ok_message "check_network_files"
}

function user_input {
    if [ -n "${MANUAL_STEPPING}" ]; then
        read -p "Press enter to continue..."
    fi
}
