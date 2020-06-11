run_deploy_charm() {
    echo

    file="${TEST_DIR}/test-deploy-charm.txt"

    ensure "test-deploy-charm" "${file}"

    juju deploy cs:~jameinel/ubuntu-lite-7
    wait_for "ubuntu-lite" "$(idle_condition "ubuntu-lite")"

    destroy_model "test-deploy-charm"
}

run_deploy_lxd_profile_charm() {
    echo

    file="${TEST_DIR}/test-deploy-lxd-profile.txt"

    ensure "test-deploy-lxd-profile" "${file}"

    juju deploy cs:~juju-qa/bionic/lxd-profile-without-devices-5
    wait_for "lxd-profile" "$(idle_condition "lxd-profile")"

    juju status --format=json | jq ".machines | .[\"0\"] | .[\"lxd-profiles\"] | keys[0]" | check "juju-test-deploy-lxd-profile-lxd-profile"

    destroy_model "test-deploy-lxd-profile"
}

run_deploy_lxd_profile_charm_container() {
    echo

    file="${TEST_DIR}/test-deploy-lxd-profile.txt"

    ensure "test-deploy-lxd-profile-container" "${file}"

    juju deploy cs:~juju-qa/bionic/lxd-profile-without-devices-5 --to lxd
    wait_for "lxd-profile" "$(idle_condition "lxd-profile")"

    juju status --format=json | jq ".machines | .[\"0\"] | .containers | .[\"0/lxd/0\"] | .[\"lxd-profiles\"] | keys[0]" |
    check "juju-test-deploy-lxd-profile-container-lxd-profile"

    destroy_model "test-deploy-lxd-profile-container"
}

run_deploy_local_lxd_profile_charm() {
    echo

    file="${TEST_DIR}/test-deploy-local-lxd-profile.txt"

    ensure "test-deploy-local-lxd-profile" "${file}"

    juju deploy ./tests/suites/deploy/charms/lxd-profile
    juju deploy ./tests/suites/deploy/charms/lxd-profile-subordinate
    juju add-relation lxd-profile-subordinate lxd-profile

    wait_for "lxd-profile" "$(idle_condition "lxd-profile")"
    wait_for "lxd-profile-subordinate" ".applications | keys[1]"

    lxd_profile_name="juju-test-deploy-local-lxd-profile-lxd-profile"
    lxd_profile_sub_name="juju-test-deploy-local-lxd-profile-lxd-profile-subordinate"

    # subordinates take longer to show, so use wait_for
    machine_0="$(machine_path 0)"
    wait_for "${lxd_profile_sub_name}" "${machine_0}"

    juju status --format=json | jq "${machine_0}" | check "${lxd_profile_name}"
    juju status --format=json | jq "${machine_0}" | check "${lxd_profile_sub_name}"

    juju add-unit "lxd-profile"

    machine_1="$(machine_path 1)"
    wait_for "${lxd_profile_sub_name}" "${machine_1}"

    juju status --format=json | jq "${machine_1}" | check "${lxd_profile_name}"
    juju status --format=json | jq "${machine_1}" | check "${lxd_profile_sub_name}"

    juju add-unit "lxd-profile" --to lxd

    machine_2="$(machine_container_path 2 2/lxd/0)"
    wait_for "${lxd_profile_sub_name}" "${machine_2}"

    juju status --format=json | jq "${machine_2}" | check "${lxd_profile_name}"
    juju status --format=json | jq "${machine_2}" | check "${lxd_profile_sub_name}"

    destroy_model "test-deploy-local-lxd-profile"
}

run_deploy_lxd_to_machine() {
    echo

    model_name="test-deploy-lxd-machine"
    file="${TEST_DIR}/${model_name}.txt"

    ensure "${model_name}" "${file}"

    juju add-machine -n 1

    charm=./tests/suites/deploy/charms/lxd-profile-alt
    juju deploy "${charm}" --to 0

    wait_for "lxd-profile-alt" "$(idle_condition "lxd-profile-alt")"

    lxc profile show "juju-test-deploy-lxd-machine-lxd-profile-alt-0" | \
        grep "linux.kernel_modules: nbd,ip_tables,ip6_tables"

    juju upgrade-charm "lxd-profile-alt" --path "${charm}"

    wait_for "lxd-profile-alt" "$(charm_rev "lxd-profile-alt" 1)"
    wait_for "lxd-profile-alt" "$(idle_condition "lxd-profile-alt")"

    attempt=0
    while true; do
        OUT=$(lxc profile show "juju-test-deploy-lxd-machine-lxd-profile-alt-1" | grep "linux.kernel_modules: nbd,ip_tables,ip6_tables" || echo 'NOT FOUND')
        if [ "${OUT}" != "NOT FOUND" ]; then
            break
        fi
        attempt=$((attempt+1))
        if [ $attempt -eq 10 ]; then
             # shellcheck disable=SC2046
             echo $(red "timeout: waiting for lxc profile to show 50sec")
             exit 5
        fi
        sleep 1
    done

    # Ensure that the old one is removed
    attempt=0
    while true; do
        OUT=$(lxc profile show "juju-test-deploy-lxd-machine-lxd-profile-alt-0" || echo 'NOT FOUND')
        if [ "${OUT}" = "NOT FOUND" ]; then
            break
        fi
        attempt=$((attempt+1))
        if [ $attempt -eq 10 ]; then
             # shellcheck disable=SC2046
             echo $(red "timeout: waiting for lxc profile to show 50sec")
             exit 5
        fi
        sleep 1
    done
}

run_deploy_lxd_to_container() {
    echo

    model_name="test-deploy-lxd-container"
    file="${TEST_DIR}/${model_name}.txt"

    ensure "${model_name}" "${file}"

    charm=./tests/suites/deploy/charms/lxd-profile-alt
    juju deploy "${charm}" --to lxd

    wait_for "lxd-profile-alt" "$(idle_condition "lxd-profile-alt")"

    juju run --application lxd-profile-alt -- su root && lxc profile show "juju-test-deploy-lxd-container-lxd-profile-alt-0" | \
        grep "linux.kernel_modules: nbd,ip_tables,ip6_tables"

    juju upgrade-charm "lxd-profile-alt" --path "${charm}"

    wait_for "lxd-profile-alt" "$(charm_rev "lxd-profile-alt" 1)"
    wait_for "lxd-profile-alt" "$(idle_condition "lxd-profile-alt")"

    attempt=0
    while true; do
        OUT=$(juju run --application lxd-profile-alt -- su root && lxc profile show "juju-test-deploy-lxd-machine-lxd-profile-alt-1" | grep "linux.kernel_modules: nbd,ip_tables,ip6_tables" || echo 'NOT FOUND')
        if [ "${OUT}" != "NOT FOUND" ]; then
            break
        fi
        attempt=$((attempt+1))
        if [ $attempt -eq 10 ]; then
             # shellcheck disable=SC2046
             echo $(red "timeout: waiting for lxc profile to show 50sec")
             exit 5
        fi
        sleep 1
    done

    # Ensure that the old one is removed
    attempt=0
    while true; do
        OUT=$(juju run --application lxd-profile-alt -- su root && lxc profile show "juju-test-deploy-lxd-machine-lxd-profile-alt-0" || echo 'NOT FOUND')
        if [ "${OUT}" = "NOT FOUND" ]; then
            break
        fi
        attempt=$((attempt+1))
        if [ $attempt -eq 10 ]; then
             # shellcheck disable=SC2046
             echo $(red "timeout: waiting for lxc profile to show 50sec")
             exit 5
        fi
        sleep 1
    done
}

test_deploy_charms() {
    if [ "$(skip 'test_deploy_charms')" ]; then
        echo "==> TEST SKIPPED: deploy charms"
        return
    fi

    (
        set_verbosity

        cd .. || exit

        run "run_deploy_charm"
        run "run_deploy_lxd_to_container"
        run "run_deploy_lxd_profile_charm_container"

        case "${BOOTSTRAP_PROVIDER:-}" in
            "lxd")
                run "run_deploy_lxd_to_machine"
                run "run_deploy_lxd_profile_charm"
                run "run_deploy_local_lxd_profile_charm"
                ;;
            "localhost")
                run "run_deploy_lxd_to_machine"
                run "run_deploy_lxd_profile_charm"
                run "run_deploy_local_lxd_profile_charm"
                ;;
            *)
                echo "==> TEST SKIPPED: deploy_local_lxd_profile_charm - tests for LXD only"
                ;;
        esac
    )
}

machine_path() {
    local machine

    machine=${1}

    echo ".machines | .[\"${machine}\"] | .[\"lxd-profiles\"] | keys"
}

machine_container_path() {
    local machine container

    machine=${1}
    container=${2}

    echo ".machines | .[\"${machine}\"] | .containers | .[\"${container}\"] | .[\"lxd-profiles\"] | keys"
}