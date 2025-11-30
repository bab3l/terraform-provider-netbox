#!/bin/bash
# Terraform Integration Test Runner
# Runs all Terraform test configurations against a running Netbox instance

set -e

# Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
TEST_ROOT="$PROJECT_ROOT/test/terraform"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Options
SKIP_DESTROY=false
VERBOSE=false
SPECIFIC_TEST=""

# Parse arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --skip-destroy)
            SKIP_DESTROY=true
            shift
            ;;
        --verbose)
            VERBOSE=true
            shift
            ;;
        --test)
            SPECIFIC_TEST="$2"
            shift 2
            ;;
        -h|--help)
            echo "Usage: $0 [--skip-destroy] [--verbose] [--test <path>]"
            echo ""
            echo "Options:"
            echo "  --skip-destroy  Don't destroy resources after test"
            echo "  --verbose       Show detailed output"
            echo "  --test <path>   Run specific test directory"
            exit 0
            ;;
        *)
            echo "Unknown option: $1"
            exit 1
            ;;
    esac
done

log_info() { echo -e "${CYAN}$1${NC}"; }
log_success() { echo -e "${GREEN}$1${NC}"; }
log_warning() { echo -e "${YELLOW}$1${NC}"; }
log_error() { echo -e "${RED}$1${NC}"; }

# Check environment
check_environment() {
    log_info "Checking environment..."
    
    # Check for required environment variables
    if [ -z "$NETBOX_SERVER_URL" ]; then
        export NETBOX_SERVER_URL="http://localhost:8000"
        log_warning "NETBOX_SERVER_URL not set, using default: $NETBOX_SERVER_URL"
    fi
    
    if [ -z "$NETBOX_API_TOKEN" ]; then
        export NETBOX_API_TOKEN="0123456789abcdef0123456789abcdef01234567"
        log_warning "NETBOX_API_TOKEN not set, using default"
    fi
    
    # Check Netbox connectivity
    if ! curl -sf -H "Authorization: Token $NETBOX_API_TOKEN" "$NETBOX_SERVER_URL/api/" > /dev/null 2>&1; then
        log_error "Cannot connect to Netbox API at $NETBOX_SERVER_URL"
        log_info "Make sure Netbox is running: docker-compose up -d"
        exit 1
    fi
    log_success "Connected to Netbox API"
    
    # Check for terraform
    if ! command -v terraform &> /dev/null; then
        log_error "Terraform not found in PATH"
        exit 1
    fi
    TF_VERSION=$(terraform version -json | jq -r '.terraform_version')
    log_success "Terraform version: $TF_VERSION"
    
    # Build and install provider
    log_info "Building and installing provider..."
    cd "$PROJECT_ROOT"
    
    if ! go build -o terraform-provider-netbox .; then
        log_error "Failed to build provider"
        exit 1
    fi
    
    # Determine OS/arch
    OS=$(go env GOOS)
    ARCH=$(go env GOARCH)
    
    # Install to local plugin directory
    PLUGIN_DIR="$HOME/.terraform.d/plugins/registry.terraform.io/bab3l/netbox/0.0.1/${OS}_${ARCH}"
    mkdir -p "$PLUGIN_DIR"
    cp terraform-provider-netbox "$PLUGIN_DIR/terraform-provider-netbox_v0.0.1"
    log_success "Provider installed to local plugin directory"
}

# Run a single test
run_test() {
    local test_path="$1"
    local test_name="$2"
    
    log_info ""
    log_info "============================================================"
    log_info "Running test: $test_name"
    log_info "Path: $test_path"
    log_info "============================================================"
    
    cd "$test_path"
    
    # Clean up any existing state
    rm -rf .terraform .terraform.lock.hcl terraform.tfstate terraform.tfstate.backup 2>/dev/null || true
    
    local status="Passed"
    local error=""
    
    # Init
    log_info "  terraform init..."
    if ! terraform init -no-color > /tmp/tf_init.log 2>&1; then
        status="Failed"
        error="terraform init failed"
        log_error "  Init failed"
        $VERBOSE && cat /tmp/tf_init.log
        return 1
    fi
    log_success "  Init completed"
    
    # Plan
    log_info "  terraform plan..."
    if ! terraform plan -no-color > /tmp/tf_plan.log 2>&1; then
        status="Failed"
        error="terraform plan failed"
        log_error "  Plan failed"
        $VERBOSE && cat /tmp/tf_plan.log
        return 1
    fi
    log_success "  Plan completed"
    
    # Apply
    log_info "  terraform apply..."
    if ! terraform apply -auto-approve -no-color > /tmp/tf_apply.log 2>&1; then
        status="Failed"
        error="terraform apply failed"
        log_error "  Apply failed"
        $VERBOSE && cat /tmp/tf_apply.log
        
        # Try cleanup
        if [ "$SKIP_DESTROY" = false ]; then
            log_warning "  Attempting cleanup..."
            terraform destroy -auto-approve -no-color > /dev/null 2>&1 || true
        fi
        return 1
    fi
    log_success "  Apply completed"
    
    # Get outputs and verify
    log_info "  terraform output..."
    local outputs
    outputs=$(terraform output -json 2>/dev/null)
    
    local verification_failed=false
    for key in $(echo "$outputs" | jq -r 'keys[]'); do
        local value
        value=$(echo "$outputs" | jq -r ".[\"$key\"].value")
        if [[ "$key" =~ (valid|match)$ ]]; then
            if [ "$value" = "true" ]; then
                log_success "    ✓ $key = $value"
            else
                log_error "    ✗ $key = $value (expected true)"
                verification_failed=true
            fi
        fi
    done
    
    if [ "$verification_failed" = true ]; then
        status="Failed"
        error="Output verification failed"
    fi
    
    # Destroy
    if [ "$SKIP_DESTROY" = false ]; then
        log_info "  terraform destroy..."
        if ! terraform destroy -auto-approve -no-color > /tmp/tf_destroy.log 2>&1; then
            log_warning "  Destroy had issues (continuing)"
        fi
        log_success "  Destroy completed"
    else
        log_warning "  Skipping destroy (--skip-destroy)"
    fi
    
    if [ "$status" = "Passed" ]; then
        log_success "  TEST PASSED"
        return 0
    else
        log_error "  TEST FAILED: $error"
        return 1
    fi
}

# Main
main() {
    log_info "Terraform Provider for Netbox - Integration Tests"
    log_info "=================================================="
    
    check_environment
    
    # Collect test directories
    declare -a test_dirs
    declare -a test_names
    
    if [ -n "$SPECIFIC_TEST" ]; then
        test_dirs+=("$SPECIFIC_TEST")
        test_names+=("$(basename "$SPECIFIC_TEST")")
    else
        # Ordered test sequence - dependencies first
        test_order=("tenant_group" "tenant" "site_group" "site")
        
        # Add resource tests in order
        for name in "${test_order[@]}"; do
            if [ -d "$TEST_ROOT/resources/$name" ]; then
                test_dirs+=("$TEST_ROOT/resources/$name")
                test_names+=("resource/$name")
            fi
        done
        
        # Add data source tests in order
        for name in "${test_order[@]}"; do
            if [ -d "$TEST_ROOT/data-sources/$name" ]; then
                test_dirs+=("$TEST_ROOT/data-sources/$name")
                test_names+=("data-source/$name")
            fi
        done
    fi
    
    if [ ${#test_dirs[@]} -eq 0 ]; then
        log_warning "No test directories found"
        exit 1
    fi
    
    log_info ""
    log_info "Found ${#test_dirs[@]} test(s) to run"
    
    # Run tests
    declare -a results
    passed=0
    failed=0
    
    for i in "${!test_dirs[@]}"; do
        if run_test "${test_dirs[$i]}" "${test_names[$i]}"; then
            results+=("PASSED: ${test_names[$i]}")
            ((passed++))
        else
            results+=("FAILED: ${test_names[$i]}")
            ((failed++))
        fi
    done
    
    # Summary
    log_info ""
    log_info "============================================================"
    log_info "TEST SUMMARY"
    log_info "============================================================"
    
    for result in "${results[@]}"; do
        if [[ "$result" == PASSED* ]]; then
            log_success "  ✓ ${result#PASSED: }"
        else
            log_error "  ✗ ${result#FAILED: }"
        fi
    done
    
    total=$((passed + failed))
    log_info ""
    
    if [ "$failed" -eq 0 ]; then
        log_success "All $total tests passed!"
        exit 0
    else
        log_error "$passed/$total tests passed, $failed failed"
        exit 1
    fi
}

main
