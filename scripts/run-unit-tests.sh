#!/usr/bin/env bash
set -euo pipefail

pkgs=$(go list ./... | grep -v acceptance_tests)

go test ${pkgs} -v
