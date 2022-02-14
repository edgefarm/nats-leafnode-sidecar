` "creds/myAccount.creds" "remote credentials not equal"
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
