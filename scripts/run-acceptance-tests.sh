#!/bin/bash
# Run acceptance tests against a local Netbox Docker instance.
# Usage: ./scripts/run-acceptance-tests.sh [options]
#
# Options:
#   --start-only    Only start Docker, don't run tests
#   --stop-only     Stop and clean up Docker environment
#   --skip-docker   Skip Docker operations, just run tests
#   --timeout SEC   Timeout for Netbox readiness (default: 300)
#   --pattern PAT   Test pattern to run specific tests

set -e

# Configuration
NETBOX_URL="http://localhost:8000"
API_TOKEN="0123456789abcdef0123456789abcdef01234567"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

# Defaults
START_ONLY=false
STOP_ONLY=false
SKIP_DOCKER=false
TIMEOUT=300
TEST_PATTERN=""

# Parse arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --start-only) START_ONLY=true; shift ;;
        --stop-only) STOP_ONLY=true; shift ;;
        --skip-docker) SKIP_DOCKER=true; shift ;;
        --timeout) TIMEOUT="$2"; shift 2 ;;
        --pattern) TEST_PATTERN="$2"; shift 2 ;;
        *) echo "Unknown option: $1"; exit 1 ;;
    esac
done

log() {
    echo "[$(date '+%H:%M:%S')] $1"
}

start_netbox() {
    log "Starting Netbox Docker environment..."
    cd "$PROJECT_ROOT"
    docker-compose up -d
    log "Docker containers started"
}

stop_netbox() {
    log "Stopping Netbox Docker environment..."
    cd "$PROJECT_ROOT"
    docker-compose down -v
    log "Docker environment stopped and cleaned up"
}

wait_for_netbox() {
    log "Waiting for Netbox to be ready (timeout: ${TIMEOUT}s)..."
    
    local elapsed=0
    while [ $elapsed -lt $TIMEOUT ]; do
        if curl -sf -H "Authorization: Token $API_TOKEN" "$NETBOX_URL/api/" > /dev/null 2>&1; then
            log "Netbox is ready!"
            return 0
        fi
        echo -ne "\r  Waiting... ($elapsed/${TIMEOUT}s) "
        sleep 5
        elapsed=$((elapsed + 5))
    done
    
    echo ""
    log "ERROR: Netbox did not become ready within ${TIMEOUT} seconds"
    return 1
}

run_tests() {
    log "Running acceptance tests..."
    
    export NETBOX_SERVER_URL="$NETBOX_URL"
    export NETBOX_API_TOKEN="$API_TOKEN"
    export TF_ACC=1
    
    cd "$PROJECT_ROOT"
    
    if [ -n "$TEST_PATTERN" ]; then
        log "Running tests matching: $TEST_PATTERN"
        go test ./... -v -timeout 120m -run "$TEST_PATTERN"
    else
        go test ./... -v -timeout 120m
    fi
}

# Main
echo ""
echo "========================================"
echo " Terraform Provider Netbox - Test Runner"
echo "========================================"
echo ""

if [ "$STOP_ONLY" = true ]; then
    stop_netbox
    exit 0
fi

if [ "$SKIP_DOCKER" = false ]; then
    if ! command -v docker &> /dev/null || ! command -v docker-compose &> /dev/null; then
        log "ERROR: Docker or docker-compose not found"
        exit 1
    fi
    
    start_netbox
    wait_for_netbox
fi

if [ "$START_ONLY" = true ]; then
    log "Docker environment is running. Netbox is available at: $NETBOX_URL"
    log "API Token: $API_TOKEN"
    log "Run tests manually with: TF_ACC=1 go test ./... -v"
    exit 0
fi

run_tests
TEST_EXIT=$?

echo ""
log "Docker environment is still running. To stop: ./scripts/run-acceptance-tests.sh --stop-only"
echo ""

exit $TEST_EXIT
