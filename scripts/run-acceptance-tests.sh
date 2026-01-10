#!/usr/bin/env bash
set -euo pipefail

: "${NETBOX_SERVER_URL:?NETBOX_SERVER_URL is not set}"
: "${NETBOX_API_TOKEN:?NETBOX_API_TOKEN is not set}"

export TF_ACC=${TF_ACC:-1}
TIMEOUT=${NETBOX_TEST_TIMEOUT:-120m}

ran_any=false

# Run all acceptance tests EXCEPT customfields packages (safe to run with default parallelism)
acceptance_pkgs=$(go list ./... | grep acceptance_tests | grep -v acceptance_tests_customfields || true)
if [[ -n "${acceptance_pkgs}" ]]; then
	ran_any=true
	echo "Running acceptance tests (non-customfields)..."
	go test ${acceptance_pkgs} -v -timeout "${TIMEOUT}"
fi

# Run customfields acceptance tests SERIAL only (shared resources)
# Note: these packages often contain ONLY files with `//go:build customfields`, so we must
# discover them with the build tag enabled.
customfields_pkgs=$(go list -tags customfields ./... | grep acceptance_tests_customfields || true)
if [[ -n "${customfields_pkgs}" ]]; then
	ran_any=true
	echo "Running acceptance tests (customfields, serial)..."
	go test -tags customfields ${customfields_pkgs} -v -timeout "${TIMEOUT}" -p 1 -parallel 1
	exit $?
fi

if [[ "${ran_any}" == "false" ]]; then
	echo "No acceptance test packages found."
else
	echo "No customfields acceptance test packages found; skipping customfields."
fi
