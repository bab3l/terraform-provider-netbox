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

# Mapping of Terraform resource types to NetBox API endpoints
# Format: "resource_type:endpoint:name_field"
declare -A RESOURCE_API_MAP=(
    # DCIM
    ["netbox_manufacturer"]="/api/dcim/manufacturers/:name"
    ["netbox_site"]="/api/dcim/sites/:name"
    ["netbox_site_group"]="/api/dcim/site-groups/:name"
    ["netbox_region"]="/api/dcim/regions/:name"
    ["netbox_location"]="/api/dcim/locations/:name"
    ["netbox_rack"]="/api/dcim/racks/:name"
    ["netbox_rack_role"]="/api/dcim/rack-roles/:name"
    ["netbox_rack_type"]="/api/dcim/rack-types/:model"
    ["netbox_device"]="/api/dcim/devices/:name"
    ["netbox_device_role"]="/api/dcim/device-roles/:name"
    ["netbox_device_type"]="/api/dcim/device-types/:model"
    ["netbox_module_type"]="/api/dcim/module-types/:model"
    ["netbox_module"]="/api/dcim/modules/:id"
    ["netbox_module_bay"]="/api/dcim/module-bays/:name"
    ["netbox_device_bay"]="/api/dcim/device-bays/:name"
    ["netbox_interface"]="/api/dcim/interfaces/:name"
    ["netbox_interface_template"]="/api/dcim/interface-templates/:name"
    ["netbox_console_port"]="/api/dcim/console-ports/:name"
    ["netbox_console_port_template"]="/api/dcim/console-port-templates/:name"
    ["netbox_console_server_port"]="/api/dcim/console-server-ports/:name"
    ["netbox_console_server_port_template"]="/api/dcim/console-server-port-templates/:name"
    ["netbox_power_port"]="/api/dcim/power-ports/:name"
    ["netbox_power_port_template"]="/api/dcim/power-port-templates/:name"
    ["netbox_power_outlet"]="/api/dcim/power-outlets/:name"
    ["netbox_power_outlet_template"]="/api/dcim/power-outlet-templates/:name"
    ["netbox_power_panel"]="/api/dcim/power-panels/:name"
    ["netbox_power_feed"]="/api/dcim/power-feeds/:name"
    ["netbox_platform"]="/api/dcim/platforms/:name"
    ["netbox_inventory_item"]="/api/dcim/inventory-items/:name"
    ["netbox_inventory_item_role"]="/api/dcim/inventory-item-roles/:name"
    ["netbox_virtual_chassis"]="/api/dcim/virtual-chassis/:name"
    ["netbox_cable"]="/api/dcim/cables/:id"
    # IPAM
    ["netbox_vrf"]="/api/ipam/vrfs/:name"
    ["netbox_vlan"]="/api/ipam/vlans/:name"
    ["netbox_vlan_group"]="/api/ipam/vlan-groups/:name"
    ["netbox_prefix"]="/api/ipam/prefixes/:prefix"
    ["netbox_ip_address"]="/api/ipam/ip-addresses/:address"
    ["netbox_ip_range"]="/api/ipam/ip-ranges/:start_address"
    ["netbox_aggregate"]="/api/ipam/aggregates/:prefix"
    ["netbox_rir"]="/api/ipam/rirs/:name"
    ["netbox_role"]="/api/ipam/roles/:name"
    ["netbox_asn"]="/api/ipam/asns/:asn"
    ["netbox_service"]="/api/ipam/services/:name"
    # Tenancy
    ["netbox_tenant"]="/api/tenancy/tenants/:name"
    ["netbox_tenant_group"]="/api/tenancy/tenant-groups/:name"
    ["netbox_contact"]="/api/tenancy/contacts/:name"
    ["netbox_contact_group"]="/api/tenancy/contact-groups/:name"
    ["netbox_contact_role"]="/api/tenancy/contact-roles/:name"
    # Virtualization
    ["netbox_cluster"]="/api/virtualization/clusters/:name"
    ["netbox_cluster_type"]="/api/virtualization/cluster-types/:name"
    ["netbox_cluster_group"]="/api/virtualization/cluster-groups/:name"
    ["netbox_virtual_machine"]="/api/virtualization/virtual-machines/:name"
    ["netbox_vm_interface"]="/api/virtualization/interfaces/:name"
    # Circuits
    ["netbox_provider"]="/api/circuits/providers/:name"
    ["netbox_provider_account"]="/api/circuits/provider-accounts/:name"
    ["netbox_provider_network"]="/api/circuits/provider-networks/:name"
    ["netbox_circuit"]="/api/circuits/circuits/:cid"
    ["netbox_circuit_type"]="/api/circuits/circuit-types/:name"
    ["netbox_circuit_termination"]="/api/circuits/circuit-terminations/:id"
    # Wireless
    ["netbox_wireless_lan"]="/api/wireless/wireless-lans/:ssid"
    ["netbox_wireless_lan_group"]="/api/wireless/wireless-lan-groups/:name"
    # Extras
    ["netbox_tag"]="/api/extras/tags/:name"
    ["netbox_custom_field"]="/api/extras/custom-fields/:name"
    ["netbox_webhook"]="/api/extras/webhooks/:name"
    ["netbox_config_context"]="/api/extras/config-contexts/:name"
    ["netbox_config_template"]="/api/extras/config-templates/:name"
)

# Deletion order - most dependent resources first (reverse dependency order)
DELETION_ORDER=(
    "netbox_cable"
    "netbox_circuit_termination"
    "netbox_vm_interface"
    "netbox_interface"
    "netbox_console_port"
    "netbox_console_server_port"
    "netbox_power_port"
    "netbox_power_outlet"
    "netbox_inventory_item"
    "netbox_module"
    "netbox_module_bay"
    "netbox_device_bay"
    "netbox_interface_template"
    "netbox_console_port_template"
    "netbox_console_server_port_template"
    "netbox_power_port_template"
    "netbox_power_outlet_template"
    "netbox_service"
    "netbox_ip_address"
    "netbox_ip_range"
    "netbox_prefix"
    "netbox_vlan"
    "netbox_aggregate"
    "netbox_asn"
    "netbox_wireless_lan"
    "netbox_virtual_machine"
    "netbox_device"
    "netbox_virtual_chassis"
    "netbox_power_feed"
    "netbox_power_panel"
    "netbox_circuit"
    "netbox_rack"
    "netbox_cluster"
    "netbox_location"
    "netbox_site"
    "netbox_provider_network"
    "netbox_provider_account"
    "netbox_provider"
    "netbox_circuit_type"
    "netbox_tenant"
    "netbox_rack_type"
    "netbox_device_type"
    "netbox_module_type"
    "netbox_platform"
    "netbox_vrf"
    "netbox_vlan_group"
    "netbox_config_template"
    "netbox_config_context"
    "netbox_webhook"
    "netbox_custom_field"
    "netbox_contact"
    "netbox_manufacturer"
    "netbox_rir"
    "netbox_tag"
    "netbox_tenant_group"
    "netbox_site_group"
    "netbox_region"
    "netbox_cluster_group"
    "netbox_cluster_type"
    "netbox_contact_group"
    "netbox_contact_role"
    "netbox_rack_role"
    "netbox_role"
    "netbox_device_role"
    "netbox_wireless_lan_group"
    "netbox_inventory_item_role"
)

# Find resource in NetBox by name
find_netbox_resource() {
    local endpoint="$1"
    local name_field="$2"
    local value="$3"
    
    local encoded_value
    encoded_value=$(printf '%s' "$value" | jq -sRr @uri)
    
    local url="${NETBOX_SERVER_URL}${endpoint}?${name_field}=${encoded_value}"
    local response
    response=$(curl -sf -H "Authorization: Token $NETBOX_API_TOKEN" "$url" 2>/dev/null)
    
    if [ $? -eq 0 ] && [ -n "$response" ]; then
        echo "$response" | jq -r '.results[].id' 2>/dev/null
    fi
}

# Delete resource from NetBox
delete_netbox_resource() {
    local endpoint="$1"
    local id="$2"
    
    local url="${NETBOX_SERVER_URL}${endpoint}${id}/"
    if curl -sf -X DELETE -H "Authorization: Token $NETBOX_API_TOKEN" "$url" > /dev/null 2>&1; then
        return 0
    fi
    return 1
}

# Extract resource info from main.tf
parse_terraform_resources() {
    local main_tf="$1"
    
    if [ ! -f "$main_tf" ]; then
        return
    fi
    
    # Use grep and sed to extract resource definitions
    # Look for: resource "netbox_xxx" "name" { ... name = "value" ... }
    local content
    content=$(cat "$main_tf")
    
    # Use perl for multiline regex parsing (more portable than complex bash)
    perl -0777 -ne '
        while (/resource\s+"(netbox_\w+)"\s+"(\w+)"\s*\{([^}]*(?:\{[^}]*\}[^}]*)*)\}/gs) {
            my ($type, $tf_name, $body) = ($1, $2, $3);
            my $value = "";
            
            if ($body =~ /\bname\s*=\s*"([^"]+)"/) { $value = $1; print "$type:name:$value\n"; }
            elsif ($body =~ /\bmodel\s*=\s*"([^"]+)"/) { $value = $1; print "$type:model:$value\n"; }
            elsif ($body =~ /\bssid\s*=\s*"([^"]+)"/) { $value = $1; print "$type:ssid:$value\n"; }
            elsif ($body =~ /\bcid\s*=\s*"([^"]+)"/) { $value = $1; print "$type:cid:$value\n"; }
            elsif ($body =~ /\bprefix\s*=\s*"([^"]+)"/) { $value = $1; print "$type:prefix:$value\n"; }
            elsif ($body =~ /\baddress\s*=\s*"([^"]+)"/) { $value = $1; print "$type:address:$value\n"; }
            elsif ($body =~ /\bstart_address\s*=\s*"([^"]+)"/) { $value = $1; print "$type:start_address:$value\n"; }
            elsif ($body =~ /\basn\s*=\s*(\d+)/) { $value = $1; print "$type:asn:$value\n"; }
        }
    ' "$main_tf"
}

# Clean up orphaned NetBox resources
cleanup_orphaned_resources() {
    local test_path="$1"
    local main_tf
    
    # If we're already in the test directory, just use main.tf
    if [ -f "main.tf" ]; then
        main_tf="main.tf"
    else
        main_tf="$test_path/main.tf"
    fi
    
    if [ ! -f "$main_tf" ]; then
        return
    fi
    
    log_info "  Checking for orphaned resources in Netbox..."
    
    # Get resources from terraform file
    local resources
    resources=$(parse_terraform_resources "$main_tf")
    
    if [ -z "$resources" ]; then
        log_info "  No resources found in main.tf to check for orphans"
        return
    fi
    
    local resource_count
    resource_count=$(echo "$resources" | wc -l)
    log_info "  Found $resource_count resources defined in main.tf"
    
    # Loop until all resources are deleted or max iterations reached
    local max_iterations=10
    local total_deleted=0
    
    for iteration in $(seq 1 $max_iterations); do
        local deleted_this_iteration=0
        local remaining_count=0
        
        # Process in deletion order
        for resource_type in "${DELETION_ORDER[@]}"; do
            local api_info="${RESOURCE_API_MAP[$resource_type]}"
            if [ -z "$api_info" ]; then
                continue
            fi
            
            local endpoint="${api_info%:*}"
            local expected_field="${api_info##*:}"
            
            # Find matching resources in our list
            while IFS=: read -r type field value; do
                if [ "$type" = "$resource_type" ] && [ -n "$value" ]; then
                    # Find in NetBox
                    local ids
                    ids=$(find_netbox_resource "$endpoint" "$expected_field" "$value")
                    
                    for id in $ids; do
                        if [ -n "$id" ]; then
                            if delete_netbox_resource "$endpoint" "$id"; then
                                log_warning "    Deleted orphaned $resource_type: $value (ID: $id)"
                                deleted_this_iteration=$((deleted_this_iteration + 1))
                                total_deleted=$((total_deleted + 1))
                            else
                                remaining_count=$((remaining_count + 1))
                            fi
                        fi
                    done
                fi
            done <<< "$resources"
        done
        
        # If nothing was deleted this iteration, we're done
        if [ "$deleted_this_iteration" -eq 0 ]; then
            break
        fi
        
        # If we deleted something but there are still remaining, continue looping
        if [ "$remaining_count" -eq 0 ]; then
            break
        fi
        
        log_info "    Iteration $iteration deleted $deleted_this_iteration, $remaining_count remaining (retrying...)"
    done
    
    if [ "$total_deleted" -gt 0 ]; then
        log_info "  Cleaned up $total_deleted orphaned resource(s)"
    fi
}

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
    
    # Clean up orphaned resources in Netbox that might be left over from failed tests
    cleanup_orphaned_resources "$test_path"
    
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
        # Order matters: dependencies must come before dependents
        # This is a comprehensive ordering of all resources based on their dependencies
        test_order=(
            # Phase 1: Core Infrastructure - No dependencies
            "manufacturer"
            "rir"
            "tag"
            "tenant_group"
            "site_group"
            "region"
            "cluster_group"
            "cluster_type"
            "contact_group"
            "contact_role"
            "rack_role"
            "role"                      # IPAM role
            "device_role"
            "wireless_lan_group"
            "vlan_group"
            
            # Phase 2: Core Infrastructure - Simple dependencies
            "platform"                  # depends on manufacturer
            "tenant"                    # depends on tenant_group
            "site"                      # depends on site_group, region, tenant
            "rack_type"                 # depends on manufacturer
            "device_type"               # depends on manufacturer
            "module_type"               # depends on manufacturer
            "inventory_item_role"
            "provider"                  # circuit provider
            "vrf"
            "custom_field"
            "config_template"
            
            # Phase 3: Location & Infrastructure
            "location"                  # depends on site
            "rack"                      # depends on site, location, rack_role, rack_type, tenant
            "power_panel"               # depends on site, location
            "power_feed"                # depends on power_panel, rack
            "cluster"                   # depends on cluster_type, cluster_group, site, tenant
            
            # Phase 4: Device Infrastructure
            "device"                    # depends on device_type, device_role, site, location, rack, tenant, platform
            "virtual_chassis"
            "device_bay"                # depends on device
            "module_bay"                # depends on device
            "module"                    # depends on device, module_bay, module_type
            
            # Phase 5: Device Components & Templates
            "console_port_template"     # depends on device_type
            "console_server_port_template" # depends on device_type
            "power_port_template"       # depends on device_type
            "power_outlet_template"     # depends on device_type
            "interface_template"        # depends on device_type
            "console_port"              # depends on device
            "console_server_port"       # depends on device
            "power_port"                # depends on device
            "power_outlet"              # depends on device, power_port
            "interface"                 # depends on device
            "inventory_item"            # depends on device, inventory_item_role, manufacturer
            
            # Phase 6: Virtualization
            "virtual_machine"           # depends on cluster, site, tenant, device, platform
            "vm_interface"              # depends on virtual_machine
            
            # Phase 7: IPAM - IP Address Management
            "asn"                       # depends on rir, tenant
            "aggregate"                 # depends on rir, tenant
            "vlan"                      # depends on vlan_group, site, tenant, role
            "prefix"                    # depends on vrf, vlan, site, tenant, role
            "ip_address"                # depends on vrf, tenant, interface, vm_interface
            "ip_range"                  # depends on vrf, tenant, role
            "service"                   # depends on device, virtual_machine, ip_address
            
            # Phase 8: Circuits
            "circuit_type"
            "provider_account"          # depends on provider
            "provider_network"          # depends on provider
            "circuit"                   # depends on provider, circuit_type, tenant
            "circuit_termination"       # depends on circuit, site, provider_network
            
            # Phase 9: Wireless
            "wireless_lan"              # depends on wireless_lan_group, vlan, tenant
            
            # Phase 10: Cabling & Connections
            "cable"                     # depends on interfaces, console ports, power ports, etc.
            
            # Phase 11: Extras & Customization
            "contact"                   # depends on contact_group
            "webhook"
            "config_context"            # depends on many resources for assignment
        )
        
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
