# Terraform Integration Test Runner
# Runs all Terraform test configurations against a running Netbox instance

param(
    [string]$TestDir = "",
    [string]$StartFrom = "",
    [switch]$SkipDestroy,
    [switch]$ShowDetails
)

$ErrorActionPreference = "Stop"
$ScriptRoot = Split-Path -Parent $MyInvocation.MyCommand.Path
$ProjectRoot = Split-Path -Parent $ScriptRoot
$TestRoot = Join-Path $ProjectRoot "test\terraform"

function Write-Success { param($msg) Write-Host $msg -ForegroundColor Green }
function Write-Failure { param($msg) Write-Host $msg -ForegroundColor Red }
function Write-Info { param($msg) Write-Host $msg -ForegroundColor Cyan }
function Write-Warn { param($msg) Write-Host $msg -ForegroundColor Yellow }

# Mapping of Terraform resource types to NetBox API endpoints
$global:ResourceApiMap = @{
    # DCIM
    "netbox_manufacturer" = @{ Endpoint = "/api/dcim/manufacturers/"; NameField = "name" }
    "netbox_site" = @{ Endpoint = "/api/dcim/sites/"; NameField = "name" }
    "netbox_site_group" = @{ Endpoint = "/api/dcim/site-groups/"; NameField = "name" }
    "netbox_region" = @{ Endpoint = "/api/dcim/regions/"; NameField = "name" }
    "netbox_location" = @{ Endpoint = "/api/dcim/locations/"; NameField = "name" }
    "netbox_rack" = @{ Endpoint = "/api/dcim/racks/"; NameField = "name" }
    "netbox_rack_role" = @{ Endpoint = "/api/dcim/rack-roles/"; NameField = "name" }
    "netbox_rack_type" = @{ Endpoint = "/api/dcim/rack-types/"; NameField = "model" }
    "netbox_device" = @{ Endpoint = "/api/dcim/devices/"; NameField = "name" }
    "netbox_device_role" = @{ Endpoint = "/api/dcim/device-roles/"; NameField = "name" }
    "netbox_device_type" = @{ Endpoint = "/api/dcim/device-types/"; NameField = "model" }
    "netbox_module_type" = @{ Endpoint = "/api/dcim/module-types/"; NameField = "model" }
    "netbox_module" = @{ Endpoint = "/api/dcim/modules/"; NameField = "id" }
    "netbox_module_bay" = @{ Endpoint = "/api/dcim/module-bays/"; NameField = "name" }
    "netbox_device_bay" = @{ Endpoint = "/api/dcim/device-bays/"; NameField = "name" }
    "netbox_device_bay_template" = @{ Endpoint = "/api/dcim/device-bay-templates/"; NameField = "name" }
    "netbox_interface" = @{ Endpoint = "/api/dcim/interfaces/"; NameField = "name" }
    "netbox_interface_template" = @{ Endpoint = "/api/dcim/interface-templates/"; NameField = "name" }
    "netbox_console_port" = @{ Endpoint = "/api/dcim/console-ports/"; NameField = "name" }
    "netbox_console_port_template" = @{ Endpoint = "/api/dcim/console-port-templates/"; NameField = "name" }
    "netbox_console_server_port" = @{ Endpoint = "/api/dcim/console-server-ports/"; NameField = "name" }
    "netbox_console_server_port_template" = @{ Endpoint = "/api/dcim/console-server-port-templates/"; NameField = "name" }
    "netbox_power_port" = @{ Endpoint = "/api/dcim/power-ports/"; NameField = "name" }
    "netbox_power_port_template" = @{ Endpoint = "/api/dcim/power-port-templates/"; NameField = "name" }
    "netbox_power_outlet" = @{ Endpoint = "/api/dcim/power-outlets/"; NameField = "name" }
    "netbox_power_outlet_template" = @{ Endpoint = "/api/dcim/power-outlet-templates/"; NameField = "name" }
    "netbox_rear_port" = @{ Endpoint = "/api/dcim/rear-ports/"; NameField = "name" }
    "netbox_rear_port_template" = @{ Endpoint = "/api/dcim/rear-port-templates/"; NameField = "name" }
    "netbox_front_port" = @{ Endpoint = "/api/dcim/front-ports/"; NameField = "name" }
    "netbox_front_port_template" = @{ Endpoint = "/api/dcim/front-port-templates/"; NameField = "name" }
    "netbox_power_panel" = @{ Endpoint = "/api/dcim/power-panels/"; NameField = "name" }
    "netbox_power_feed" = @{ Endpoint = "/api/dcim/power-feeds/"; NameField = "name" }
    "netbox_platform" = @{ Endpoint = "/api/dcim/platforms/"; NameField = "name" }
    "netbox_inventory_item" = @{ Endpoint = "/api/dcim/inventory-items/"; NameField = "name" }
    "netbox_inventory_item_role" = @{ Endpoint = "/api/dcim/inventory-item-roles/"; NameField = "name" }
    "netbox_inventory_item_template" = @{ Endpoint = "/api/dcim/inventory-item-templates/"; NameField = "name" }
    "netbox_virtual_chassis" = @{ Endpoint = "/api/dcim/virtual-chassis/"; NameField = "name" }
    "netbox_rack_reservation" = @{ Endpoint = "/api/dcim/rack-reservations/"; NameField = "id" }
    "netbox_virtual_device_context" = @{ Endpoint = "/api/dcim/virtual-device-contexts/"; NameField = "name" }
    "netbox_module_bay_template" = @{ Endpoint = "/api/dcim/module-bay-templates/"; NameField = "name" }
    "netbox_cable" = @{ Endpoint = "/api/dcim/cables/"; NameField = "id" }

    # IPAM
    "netbox_vrf" = @{ Endpoint = "/api/ipam/vrfs/"; NameField = "name" }
    "netbox_vlan" = @{ Endpoint = "/api/ipam/vlans/"; NameField = "name" }
    "netbox_vlan_group" = @{ Endpoint = "/api/ipam/vlan-groups/"; NameField = "name" }
    "netbox_prefix" = @{ Endpoint = "/api/ipam/prefixes/"; NameField = "prefix" }
    "netbox_ip_address" = @{ Endpoint = "/api/ipam/ip-addresses/"; NameField = "address" }
    "netbox_ip_range" = @{ Endpoint = "/api/ipam/ip-ranges/"; NameField = "start_address" }
    "netbox_aggregate" = @{ Endpoint = "/api/ipam/aggregates/"; NameField = "prefix" }
    "netbox_rir" = @{ Endpoint = "/api/ipam/rirs/"; NameField = "name" }
    "netbox_role" = @{ Endpoint = "/api/ipam/roles/"; NameField = "name" }
    "netbox_asn" = @{ Endpoint = "/api/ipam/asns/"; NameField = "asn" }
    "netbox_asn_range" = @{ Endpoint = "/api/ipam/asn-ranges/"; NameField = "name" }
    "netbox_service" = @{ Endpoint = "/api/ipam/services/"; NameField = "name" }
    "netbox_route_target" = @{ Endpoint = "/api/ipam/route-targets/"; NameField = "name" }
    "netbox_fhrp_group" = @{ Endpoint = "/api/ipam/fhrp-groups/"; NameField = "id" }
    "netbox_fhrp_group_assignment" = @{ Endpoint = "/api/ipam/fhrp-group-assignments/"; NameField = "id" }
    "netbox_service_template" = @{ Endpoint = "/api/ipam/service-templates/"; NameField = "name" }
    "netbox_virtual_disk" = @{ Endpoint = "/api/virtualization/virtual-disks/"; NameField = "name" }

    # Tenancy
    "netbox_tenant" = @{ Endpoint = "/api/tenancy/tenants/"; NameField = "name" }
    "netbox_tenant_group" = @{ Endpoint = "/api/tenancy/tenant-groups/"; NameField = "name" }
    "netbox_contact" = @{ Endpoint = "/api/tenancy/contacts/"; NameField = "name" }
    "netbox_contact_group" = @{ Endpoint = "/api/tenancy/contact-groups/"; NameField = "name" }
    "netbox_contact_role" = @{ Endpoint = "/api/tenancy/contact-roles/"; NameField = "name" }

    # Virtualization
    "netbox_cluster" = @{ Endpoint = "/api/virtualization/clusters/"; NameField = "name" }
    "netbox_cluster_type" = @{ Endpoint = "/api/virtualization/cluster-types/"; NameField = "name" }
    "netbox_cluster_group" = @{ Endpoint = "/api/virtualization/cluster-groups/"; NameField = "name" }
    "netbox_virtual_machine" = @{ Endpoint = "/api/virtualization/virtual-machines/"; NameField = "name" }
    "netbox_vm_interface" = @{ Endpoint = "/api/virtualization/interfaces/"; NameField = "name" }

    # Circuits
    "netbox_provider" = @{ Endpoint = "/api/circuits/providers/"; NameField = "name" }
    "netbox_provider_account" = @{ Endpoint = "/api/circuits/provider-accounts/"; NameField = "name" }
    "netbox_provider_network" = @{ Endpoint = "/api/circuits/provider-networks/"; NameField = "name" }
    "netbox_circuit" = @{ Endpoint = "/api/circuits/circuits/"; NameField = "cid" }
    "netbox_circuit_type" = @{ Endpoint = "/api/circuits/circuit-types/"; NameField = "name" }
    "netbox_circuit_termination" = @{ Endpoint = "/api/circuits/circuit-terminations/"; NameField = "id" }
    "netbox_circuit_group" = @{ Endpoint = "/api/circuits/circuit-groups/"; NameField = "name" }
    "netbox_circuit_group_assignment" = @{ Endpoint = "/api/circuits/circuit-group-assignments/"; NameField = "id" }

    # Wireless
    "netbox_wireless_lan" = @{ Endpoint = "/api/wireless/wireless-lans/"; NameField = "ssid" }
    "netbox_wireless_lan_group" = @{ Endpoint = "/api/wireless/wireless-lan-groups/"; NameField = "name" }

    # VPN
    "netbox_ike_proposal" = @{ Endpoint = "/api/vpn/ike-proposals/"; NameField = "name" }
    "netbox_ike_policy" = @{ Endpoint = "/api/vpn/ike-policies/"; NameField = "name" }
    "netbox_ipsec_proposal" = @{ Endpoint = "/api/vpn/ipsec-proposals/"; NameField = "name" }
    "netbox_ipsec_policy" = @{ Endpoint = "/api/vpn/ipsec-policies/"; NameField = "name" }
    "netbox_ipsec_profile" = @{ Endpoint = "/api/vpn/ipsec-profiles/"; NameField = "name" }
    "netbox_tunnel_group" = @{ Endpoint = "/api/vpn/tunnel-groups/"; NameField = "name" }
    "netbox_tunnel" = @{ Endpoint = "/api/vpn/tunnels/"; NameField = "name" }
    "netbox_tunnel_termination" = @{ Endpoint = "/api/vpn/tunnel-terminations/"; NameField = "id" }
    "netbox_l2vpn" = @{ Endpoint = "/api/vpn/l2vpns/"; NameField = "name" }
    "netbox_l2vpn_termination" = @{ Endpoint = "/api/vpn/l2vpn-terminations/"; NameField = "id" }

    # Extras
    "netbox_tag" = @{ Endpoint = "/api/extras/tags/"; NameField = "name" }
    "netbox_custom_field" = @{ Endpoint = "/api/extras/custom-fields/"; NameField = "name" }
    "netbox_custom_field_choice_set" = @{ Endpoint = "/api/extras/custom-field-choice-sets/"; NameField = "name" }
    "netbox_custom_link" = @{ Endpoint = "/api/extras/custom-links/"; NameField = "name" }
    "netbox_webhook" = @{ Endpoint = "/api/extras/webhooks/"; NameField = "name" }
    "netbox_config_context" = @{ Endpoint = "/api/extras/config-contexts/"; NameField = "name" }
    "netbox_config_template" = @{ Endpoint = "/api/extras/config-templates/"; NameField = "name" }
    "netbox_journal_entry" = @{ Endpoint = "/api/extras/journal-entries/"; NameField = "id" }
    "netbox_export_template" = @{ Endpoint = "/api/extras/export-templates/"; NameField = "name" }
}

function Get-NetboxHeaders {
    return @{
        "Authorization" = "Token $($env:NETBOX_API_TOKEN)"
        "Content-Type" = "application/json"
    }
}

function Remove-NetboxResource {
    param(
        [string]$Endpoint,
        [int]$Id
    )

    $headers = Get-NetboxHeaders
    $uri = "$($env:NETBOX_SERVER_URL)$Endpoint$Id/"

    try {
        Invoke-RestMethod -Uri $uri -Method Delete -Headers $headers | Out-Null
        return $true
    }
    catch {
        # Resource couldn't be deleted (likely has dependencies)
        return $false
    }
}

function Find-NetboxResourceByName {
    param(
        [string]$Endpoint,
        [string]$NameField,
        [string]$Value
    )

    $headers = Get-NetboxHeaders
    $encodedValue = [uri]::EscapeDataString($Value)
    $uri = "$($env:NETBOX_SERVER_URL)$Endpoint`?$NameField=$encodedValue"

    try {
        $response = Invoke-RestMethod -Uri $uri -Method Get -Headers $headers
        if ($response.results -and $response.results.Count -gt 0) {
            return $response.results
        }
    }
    catch {
        # Ignore errors - resource might not exist
    }

    return @()
}

function Get-ResourceNamesFromTerraform {
    param(
        [string]$MainTfPath
    )

    $resources = @()

    if (-not (Test-Path $MainTfPath)) {
        return $resources
    }

    $content = Get-Content $MainTfPath -Raw

    # Pattern to match resource blocks: resource "netbox_xxx" "name" { ... name = "value" ... }
    $resourcePattern = 'resource\s+"(netbox_\w+)"\s+"(\w+)"\s*\{([^}]*(?:\{[^}]*\}[^}]*)*)\}'
    $regexMatches = [regex]::Matches($content, $resourcePattern, [System.Text.RegularExpressions.RegexOptions]::Singleline)

    foreach ($match in $regexMatches) {
        $resourceType = $match.Groups[1].Value
        $terraformName = $match.Groups[2].Value
        $resourceBody = $match.Groups[3].Value

        $apiInfo = $global:ResourceApiMap[$resourceType]
        if (-not $apiInfo) {
            continue
        }

        $nameField = $apiInfo.NameField
        $nameValue = $null

        # Extract the name/identifier value from the resource body
        # Handle different field names based on resource type
        switch ($nameField) {
            "name" {
                if ($resourceBody -match '\bname\s*=\s*"([^"]+)"') {
                    $nameValue = $Matches[1]
                }
            }
            "model" {
                if ($resourceBody -match '\bmodel\s*=\s*"([^"]+)"') {
                    $nameValue = $Matches[1]
                }
            }
            "ssid" {
                if ($resourceBody -match '\bssid\s*=\s*"([^"]+)"') {
                    $nameValue = $Matches[1]
                }
            }
            "cid" {
                if ($resourceBody -match '\bcid\s*=\s*"([^"]+)"') {
                    $nameValue = $Matches[1]
                }
            }
            "prefix" {
                if ($resourceBody -match '\bprefix\s*=\s*"([^"]+)"') {
                    $nameValue = $Matches[1]
                }
            }
            "address" {
                if ($resourceBody -match '\baddress\s*=\s*"([^"]+)"') {
                    $nameValue = $Matches[1]
                }
            }
            "start_address" {
                if ($resourceBody -match '\bstart_address\s*=\s*"([^"]+)"') {
                    $nameValue = $Matches[1]
                }
            }
            "asn" {
                if ($resourceBody -match '\basn\s*=\s*(\d+)') {
                    $nameValue = $Matches[1]
                }
            }
        }

        if ($nameValue) {
            $resources += @{
                ResourceType = $resourceType
                TerraformName = $terraformName
                NameField = $nameField
                NameValue = $nameValue
                Endpoint = $apiInfo.Endpoint
            }
        }
    }

    return $resources
}

function Clear-AllTestAggregates {
    <#
    .SYNOPSIS
    Removes ALL aggregates with private/test IP ranges to prevent overlapping aggregate errors.
    #>
    $headers = Get-NetboxHeaders
    $deleted = 0

    try {
        # Get all aggregates
        $uri = "$($env:NETBOX_SERVER_URL)/api/ipam/aggregates/?limit=1000"
        $response = Invoke-RestMethod -Uri $uri -Method Get -Headers $headers

        if (-not $response.results) {
            return $deleted
        }

        # Test prefixes to check for: private/reserved IP ranges commonly used in tests
        $testPrefixes = @(
            "10.",           # 10.0.0.0/8
            "172.16.",       # 172.16.0.0/12
            "172.17.",
            "172.18.",
            "172.19.",
            "172.20.",
            "172.21.",
            "172.22.",
            "172.23.",
            "172.24.",
            "172.25.",
            "172.26.",
            "172.27.",
            "172.28.",
            "172.29.",
            "172.30.",
            "172.31.",
            "192.168.",      # 192.168.0.0/16
            "192.0.2.",      # TEST-NET-1
            "192.0.",        # Reserved test ranges
            "198.51.",       # TEST-NET-2
            "203.0.113."     # TEST-NET-3
        )

        foreach ($aggregate in $response.results) {
            $prefix = $aggregate.prefix
            $isTestPrefix = $false

            foreach ($testPrefix in $testPrefixes) {
                if ($prefix.StartsWith($testPrefix)) {
                    $isTestPrefix = $true
                    break
                }
            }

            if ($isTestPrefix) {
                try {
                    $deleteUri = "$($env:NETBOX_SERVER_URL)/api/ipam/aggregates/$($aggregate.id)/"
                    Invoke-RestMethod -Uri $deleteUri -Method Delete -Headers $headers -ErrorAction Stop | Out-Null
                    $deleted++
                    Write-Info "      Deleted test aggregate: $prefix (ID: $($aggregate.id))"
                }
                catch {
                    # Might fail if it has dependencies, that's ok - we'll try again later
                    Write-Info "      Could not delete aggregate $prefix - may have dependencies"
                }
            }
        }
    }
    catch {
        Write-Info "      Error cleaning test aggregates: $_"
    }

    return $deleted
}

function Clear-AllTestIPRanges {
    <#
    .SYNOPSIS
    Removes ALL IP ranges with private/test IP ranges to prevent overlapping range errors.
    #>
    $headers = Get-NetboxHeaders
    $deleted = 0

    try {
        # Get all IP ranges
        $uri = "$($env:NETBOX_SERVER_URL)/api/ipam/ip-ranges/?limit=1000"
        $response = Invoke-RestMethod -Uri $uri -Method Get -Headers $headers

        if (-not $response.results) {
            return $deleted
        }

        # Test IP ranges to check for: private/reserved IP ranges commonly used in tests
        $testRanges = @(
            "10.",           # 10.0.0.0/8
            "172.16.",       # 172.16.0.0/12
            "172.17.",
            "172.18.",
            "172.19.",
            "172.20.",
            "172.21.",
            "172.22.",
            "172.23.",
            "172.24.",
            "172.25.",
            "172.26.",
            "172.27.",
            "172.28.",
            "172.29.",
            "172.30.",
            "172.31.",
            "192.168.",      # 192.168.0.0/16
            "192.0.2.",      # TEST-NET-1
            "192.0.",        # Reserved test ranges
            "198.51.",       # TEST-NET-2
            "203.0.113."     # TEST-NET-3
        )

        foreach ($range in $response.results) {
            $startAddr = $range.start_address
            $isTestRange = $false

            foreach ($testRange in $testRanges) {
                if ($startAddr.StartsWith($testRange)) {
                    $isTestRange = $true
                    break
                }
            }

            if ($isTestRange) {
                try {
                    $deleteUri = "$($env:NETBOX_SERVER_URL)/api/ipam/ip-ranges/$($range.id)/"
                    Invoke-RestMethod -Uri $deleteUri -Method Delete -Headers $headers -ErrorAction Stop | Out-Null
                    $deleted++
                    Write-Info "      Deleted test IP range: $startAddr (ID: $($range.id))"
                }
                catch {
                    # Might fail if it has dependencies, that's ok - we'll try again later
                    Write-Info "      Could not delete IP range $startAddr - may have dependencies"
                }
            }
        }
    }
    catch {
        Write-Info "      Error cleaning test IP ranges: $_"
    }

    return $deleted
}

function Clear-OrphanedNetboxResources {
    param(
        [string]$TestPath
    )

    # If we're already in the test directory (via Push-Location), just use main.tf
    # Otherwise join the path
    if (Test-Path "main.tf") {
        $mainTfPath = "main.tf"
    } else {
        $mainTfPath = Join-Path $TestPath "main.tf"
    }

    $resources = Get-ResourceNamesFromTerraform -MainTfPath $mainTfPath

    if ($resources.Count -eq 0) {
        Write-Info "  No resources found in main.tf to check for orphans"
        return
    }

    Write-Info "  Checking for orphaned resources in Netbox ($($resources.Count) resources defined)..."

    # Special handling for aggregates: clean up ALL test aggregates (those with private/test prefixes)
    # to avoid overlapping aggregate errors
    if ($resources | Where-Object { $_.ResourceType -eq 'netbox_aggregate' }) {
        Write-Info "    Pre-cleaning all test aggregates to avoid overlaps..."
        $testAggregates = Clear-AllTestAggregates
        if ($testAggregates -gt 0) {
            Write-Warn "    Pre-cleaned $testAggregates test aggregate(s)"
        }
    }

    # Special handling for IP ranges: clean up ALL test IP ranges (those with private/test ranges)
    # to avoid overlapping range errors
    if ($resources | Where-Object { $_.ResourceType -eq 'netbox_ip_range' }) {
        Write-Info "    Pre-cleaning all test IP ranges to avoid overlaps..."
        $testRanges = Clear-AllTestIPRanges
        if ($testRanges -gt 0) {
            Write-Warn "    Pre-cleaned $testRanges test IP range(s)"
        }
    }

    # Reverse order to handle dependencies (delete children first)
    $deletionOrder = @(
        # Components and connections first
        "netbox_cable",
        "netbox_circuit_termination",
        "netbox_vm_interface",
        "netbox_interface",
        "netbox_console_port",
        "netbox_console_server_port",
        "netbox_power_port",
        "netbox_power_outlet",
        "netbox_inventory_item",
        "netbox_module",
        "netbox_module_bay",
        "netbox_device_bay",
        "netbox_interface_template",
        "netbox_console_port_template",
        "netbox_console_server_port_template",
        "netbox_power_port_template",
        "netbox_power_outlet_template",
        "netbox_service",
        "netbox_ip_address",
        "netbox_ip_range",
        "netbox_prefix",
        "netbox_vlan",
        "netbox_aggregate",
        "netbox_asn",
        "netbox_asn_range",
        "netbox_route_target",
        "netbox_virtual_disk",
        "netbox_device_bay_template",
        "netbox_module_bay_template",
        "netbox_inventory_item_template",
        "netbox_rack_reservation",
        "netbox_virtual_device_context",
        "netbox_wireless_lan",
        "netbox_virtual_machine",
        "netbox_device",
        "netbox_virtual_chassis",
        "netbox_power_feed",
        "netbox_power_panel",
        "netbox_circuit",
        "netbox_rack",
        "netbox_cluster",
        "netbox_location",
        "netbox_site",
        "netbox_provider_network",
        "netbox_provider_account",
        "netbox_provider",
        "netbox_circuit_type",
        "netbox_tenant",
        "netbox_rack_type",
        "netbox_device_type",
        "netbox_module_type",
        "netbox_platform",
        "netbox_vrf",
        "netbox_vlan_group",
        "netbox_config_template",
        "netbox_config_context",
        "netbox_webhook",
        "netbox_custom_field",
        "netbox_contact",
        "netbox_manufacturer",
        "netbox_rir",
        "netbox_tag",
        "netbox_tenant_group",
        "netbox_site_group",
        "netbox_region",
        "netbox_cluster_group",
        "netbox_cluster_type",
        "netbox_contact_group",
        "netbox_contact_role",
        "netbox_rack_role",
        "netbox_role",
        "netbox_device_role",
        "netbox_wireless_lan_group",
        "netbox_inventory_item_role"
    )

    $resourcesByType = @{}
    foreach ($resource in $resources) {
        if (-not $resourcesByType.ContainsKey($resource.ResourceType)) {
            $resourcesByType[$resource.ResourceType] = @()
        }
        $resourcesByType[$resource.ResourceType] += $resource
    }

    # Loop until all resources are deleted or max iterations reached
    # This handles dependency issues where some resources can't be deleted until others are gone
    $maxIterations = 10
    $totalDeleted = 0

    for ($iteration = 1; $iteration -le $maxIterations; $iteration++) {
        $deletedThisIteration = 0
        $remainingCount = 0

        foreach ($resourceType in $deletionOrder) {
            if (-not $resourcesByType.ContainsKey($resourceType)) {
                continue
            }

            foreach ($resource in $resourcesByType[$resourceType]) {
                $existing = Find-NetboxResourceByName -Endpoint $resource.Endpoint -NameField $resource.NameField -Value $resource.NameValue

                if ($existing.Count -gt 0) {
                    Write-Info "    Found $($existing.Count) $($resource.ResourceType)(s) named '$($resource.NameValue)'"
                }

                foreach ($item in $existing) {
                    $deleted = Remove-NetboxResource -Endpoint $resource.Endpoint -Id $item.id
                    if ($deleted) {
                        Write-Warn "    Deleted orphaned $($resource.ResourceType): $($resource.NameValue) (ID: $($item.id))"
                        $deletedThisIteration++
                        $totalDeleted++
                    } else {
                        Write-Info "    Could not delete $($resource.ResourceType): $($resource.NameValue) (ID: $($item.id)) - has dependencies"
                        $remainingCount++
                    }
                }
            }
        }

        # If nothing was deleted this iteration, we're done (either all clean or stuck)
        if ($deletedThisIteration -eq 0) {
            break
        }

        # If we deleted something but there are still remaining, continue looping
        if ($remainingCount -eq 0) {
            break
        }

        Write-Info "    Iteration $iteration deleted $deletedThisIteration, $remainingCount remaining (retrying...)"
    }

    if ($totalDeleted -gt 0) {
        Write-Info "  Cleaned up $totalDeleted orphaned resource(s)"
    }
}

function Test-Environment {
    Write-Info "Checking environment..."

    # Disable Terraform debug logging to avoid output pollution
    $env:TF_LOG = ""

    if (-not $env:NETBOX_SERVER_URL) {
        $env:NETBOX_SERVER_URL = "http://localhost:8000"
        Write-Warn "NETBOX_SERVER_URL not set, using default: $($env:NETBOX_SERVER_URL)"
    }

    if (-not $env:NETBOX_API_TOKEN) {
        $env:NETBOX_API_TOKEN = "0123456789abcdef0123456789abcdef01234567"
        Write-Warn "NETBOX_API_TOKEN not set, using default"
    }

    try {
        $headers = @{ "Authorization" = "Token $($env:NETBOX_API_TOKEN)" }
        Invoke-RestMethod -Uri "$($env:NETBOX_SERVER_URL)/api/" -Headers $headers -TimeoutSec 5 | Out-Null
        Write-Success "Connected to Netbox API"
    }
    catch {
        Write-Failure "Cannot connect to Netbox API at $($env:NETBOX_SERVER_URL)"
        Write-Failure "Error: $($_.Exception.Message)"
        Write-Info "Make sure Netbox is running: docker-compose up -d"
        exit 1
    }

    try {
        $tfVersion = terraform version -json 2>$null | ConvertFrom-Json
        Write-Success "Terraform version: $($tfVersion.terraform_version)"
    }
    catch {
        Write-Failure "Terraform not found in PATH: $($_.Exception.Message)"
        exit 1
    }

    # Check for dev_overrides in terraform config
    $tfrcPath = "$env:APPDATA\terraform.rc"
    if (Test-Path $tfrcPath) {
        $tfrcContent = Get-Content $tfrcPath -Raw
        if ($tfrcContent -match "dev_overrides") {
            Write-Success "Using dev_overrides from terraform.rc"
            # Still need to build the provider
            Write-Info "Building provider..."
            Push-Location $ProjectRoot
            try {
                go build -o terraform-provider-netbox.exe .
                if ($LASTEXITCODE -ne 0) {
                    Write-Failure "Failed to build provider"
                    exit 1
                }
                Write-Success "Provider built successfully"
            }
            finally {
                Pop-Location
            }
            return
        }
    }

    Write-Info "Building and installing provider..."
    Push-Location $ProjectRoot
    try {
        go build -o terraform-provider-netbox.exe .
        if ($LASTEXITCODE -ne 0) {
            Write-Failure "Failed to build provider"
            exit 1
        }

        $pluginDir = "$env:APPDATA\terraform.d\plugins\registry.terraform.io\bab3l\netbox\0.0.1\windows_amd64"
        if (-not (Test-Path $pluginDir)) {
            New-Item -ItemType Directory -Path $pluginDir -Force | Out-Null
        }
        Copy-Item "terraform-provider-netbox.exe" "$pluginDir\terraform-provider-netbox_v0.0.1.exe" -Force
        Write-Success "Provider installed to local plugin directory"
    }
    finally {
        Pop-Location
    }
}

function Invoke-TerraformTest {
    param(
        [string]$TestPath,
        [string]$TestName
    )

    Write-Info ""
    Write-Info ("=" * 60)
    Write-Info "Running test: $TestName"
    Write-Info "Path: $TestPath"
    Write-Info ("=" * 60)

    Push-Location $TestPath
    $result = @{
        Name = $TestName
        Path = $TestPath
        Status = "Unknown"
        Error = ""
        TotalTime = 0
    }

    try {
        # Clean up any existing state
        Remove-Item -Recurse -Force ".terraform" -ErrorAction SilentlyContinue
        Remove-Item -Force ".terraform.lock.hcl" -ErrorAction SilentlyContinue
        Remove-Item -Force "terraform.tfstate" -ErrorAction SilentlyContinue
        Remove-Item -Force "terraform.tfstate.backup" -ErrorAction SilentlyContinue

        # Clean up orphaned resources in Netbox that might be left over from failed tests
        Clear-OrphanedNetboxResources -TestPath $TestPath

        $startTime = Get-Date

        # Check for dev_overrides - if present, skip init as per Terraform guidance
        $tfrcPath = "$env:APPDATA\terraform.rc"
        $useDevOverrides = $false
        if (Test-Path $tfrcPath) {
            $tfrcContent = Get-Content $tfrcPath -Raw
            if ($tfrcContent -match "dev_overrides") {
                $useDevOverrides = $true
            }
        }

        if ($useDevOverrides) {
            Write-Info "  Using dev_overrides - skipping terraform init"
        }
        else {
            # Init
            Write-Info "  terraform init..."
            $initOutput = terraform init -no-color 2>&1
            if ($LASTEXITCODE -ne 0) {
                throw "terraform init failed: $initOutput"
            }
            Write-Success "  Init completed"
        }

        # Plan
        Write-Info "  terraform plan..."
        $planOutput = terraform plan -no-color 2>&1 | Out-String
        if ($LASTEXITCODE -ne 0) {
            throw "terraform plan failed: $planOutput"
        }
        Write-Success "  Plan completed"

        # Apply
        Write-Info "  terraform apply..."
        $applyOutput = terraform apply -auto-approve -no-color 2>&1 | Out-String
        if ($LASTEXITCODE -ne 0) {
            Write-Failure "  Apply output: $applyOutput"
            throw "terraform apply failed"
        }
        Write-Success "  Apply completed"

        # Get outputs
        Write-Info "  terraform output..."
        $outputJson = terraform output -json 2>&1
        $outputs = $outputJson | ConvertFrom-Json

        if ($ShowDetails) {
            Write-Info "  Outputs:"
            $outputJson | Write-Host
        }

        # Verify outputs
        $verifyFailed = $false
        foreach ($prop in $outputs.PSObject.Properties) {
            $propName = $prop.Name
            $propValue = $prop.Value.value
            if ($propName -like "*valid" -or $propName -like "*match*") {
                if ($propValue -eq $true) {
                    Write-Success "    [OK] $propName = $propValue"
                }
                else {
                    Write-Failure "    [FAIL] $propName = $propValue (expected true)"
                    $verifyFailed = $true
                }
            }
        }

        if ($verifyFailed) {
            throw "Output verification failed"
        }

        # Destroy
        if (-not $SkipDestroy) {
            Write-Info "  terraform destroy..."
            $destroyOutput = terraform destroy -auto-approve -no-color 2>&1 | Out-String
            if ($LASTEXITCODE -ne 0) {
                throw "terraform destroy failed: $destroyOutput"
            }
            Write-Success "  Destroy completed"
        }
        else {
            Write-Warn "  Skipping destroy (SkipDestroy flag set)"
        }

        $result.TotalTime = ((Get-Date) - $startTime).TotalSeconds
        $result.Status = "Passed"
        Write-Success "  TEST PASSED"
    }
    catch {
        $result.Status = "Failed"
        $result.Error = $_.Exception.Message
        Write-Failure "  TEST FAILED: $($_.Exception.Message)"

        if (-not $SkipDestroy) {
            Write-Warn "  Attempting cleanup..."
            terraform destroy -auto-approve -no-color 2>&1 | Out-Null
        }
    }
    finally {
        Pop-Location
    }

    return $result
}

function Main {
    Write-Info "Terraform Provider for Netbox - Integration Tests"
    Write-Info "=================================================="

    Test-Environment

    $testDirs = @()

    if ($TestDir) {
        $testDirs += @{ Path = $TestDir; Name = (Split-Path -Leaf $TestDir) }
    }
    else {
        $resourceTests = Get-ChildItem -Path (Join-Path $TestRoot "resources") -Directory -ErrorAction SilentlyContinue
        $dataSourceTests = Get-ChildItem -Path (Join-Path $TestRoot "data-sources") -Directory -ErrorAction SilentlyContinue

        # Order matters: dependencies must come before dependents
        # This is a comprehensive ordering of all resources based on their dependencies
        $testOrder = @(
            # Phase 1: Core Infrastructure - No dependencies
            "manufacturer",
            "rir",
            "tag",
            "tenant_group",
            "site_group",
            "region",
            "cluster_group",
            "cluster_type",
            "contact_group",
            "contact_role",
            "rack_role",
            "role",                      # IPAM role
            "device_role",
            "wireless_lan_group",
            "vlan_group",

            # Phase 2: Core Infrastructure - Simple dependencies
            "platform",                  # depends on manufacturer
            "tenant",                    # depends on tenant_group
            "site",                      # depends on site_group, region, tenant
            "rack_type",                 # depends on manufacturer
            "device_type",               # depends on manufacturer
            "module_type",               # depends on manufacturer
            "inventory_item_role",
            "provider",                  # circuit provider
            "vrf",
            "custom_field",
            "config_template",

            # Phase 3: Location & Infrastructure
            "location",                  # depends on site
            "rack",                      # depends on site, location, rack_role, rack_type, tenant
            "power_panel",               # depends on site, location
            "power_feed",                # depends on power_panel, rack
            "cluster",                   # depends on cluster_type, cluster_group, site, tenant

            # Phase 4: Device Infrastructure
            "device",                    # depends on device_type, device_role, site, location, rack, tenant, platform
            "virtual_chassis",
            "device_bay",                # depends on device
            "module_bay",                # depends on device
            "module",                    # depends on device, module_bay, module_type

            # Phase 5: Device Components & Templates
            "console_port_template",     # depends on device_type
            "console_server_port_template", # depends on device_type
            "power_port_template",       # depends on device_type
            "power_outlet_template",     # depends on device_type
            "interface_template",        # depends on device_type
            "device_bay_template",       # depends on device_type
            "module_bay_template",       # depends on device_type
            "inventory_item_template",   # depends on device_type
            "console_port",              # depends on device
            "console_server_port",       # depends on device
            "power_port",                # depends on device
            "power_outlet",              # depends on device, power_port
            "interface",                 # depends on device
            "inventory_item",            # depends on device, inventory_item_role, manufacturer
            "rack_reservation",          # depends on rack, user
            "virtual_device_context",    # depends on device

            # Phase 6: Virtualization
            "virtual_machine",           # depends on cluster, site, tenant, device, platform
            "vm_interface",              # depends on virtual_machine
            "virtual_disk",              # depends on virtual_machine

            # Phase 7: IPAM - IP Address Management
            "asn",                       # depends on rir, tenant
            "asn_range",                 # depends on rir, tenant
            "aggregate",                 # depends on rir, tenant
            "vlan",                      # depends on vlan_group, site, tenant, role
            "prefix",                    # depends on vrf, vlan, site, tenant, role
            "ip_address",                # depends on vrf, tenant, interface, vm_interface
            "ip_range",                  # depends on vrf, tenant, role
            "route_target",              # depends on tenant
            "service",                   # depends on device, virtual_machine, ip_address
            "service_template",          # no dependencies
            "fhrp_group",                # depends on auth settings
            "fhrp_group_assignment",     # depends on fhrp_group, interface

            # Phase 8: Circuits
            "circuit_type",
            "provider_account",          # depends on provider
            "provider_network",          # depends on provider
            "circuit",                   # depends on provider, circuit_type, tenant
            "circuit_termination",       # depends on circuit, site, provider_network

            # Phase 9: Wireless
            "wireless_lan",              # depends on wireless_lan_group, vlan, tenant

            # Phase 10: Cabling & Connections
            "cable",                     # depends on interfaces, console ports, power ports, etc.

            # Phase 11: Extras & Customization
            "contact",                   # depends on contact_group
            "webhook",
            "config_context",            # depends on many resources for assignment
            "export_template"            # no dependencies
        )

        foreach ($name in $testOrder) {
            $test = $resourceTests | Where-Object { $_.Name -eq $name }
            if ($test) {
                $testDirs += @{ Path = $test.FullName; Name = "resource/$($test.Name)" }
            }
        }

        foreach ($name in $testOrder) {
            $test = $dataSourceTests | Where-Object { $_.Name -eq $name }
            if ($test) {
                $testDirs += @{ Path = $test.FullName; Name = "data-source/$($test.Name)" }
            }
        }
    }

    if ($testDirs.Count -eq 0) {
        Write-Warn "No test directories found"
        exit 1
    }

    # If StartFrom is specified, skip tests until we reach it
    if ($StartFrom) {
        $startIndex = -1
        for ($i = 0; $i -lt $testDirs.Count; $i++) {
            if ($testDirs[$i].Name -eq $StartFrom -or $testDirs[$i].Name -like "*/$StartFrom") {
                $startIndex = $i
                break
            }
        }

        if ($startIndex -eq -1) {
            Write-Warn "StartFrom test '$StartFrom' not found. Available tests:"
            foreach ($test in $testDirs) {
                Write-Info "  $($test.Name)"
            }
            exit 1
        }

        $testDirs = $testDirs[$startIndex..($testDirs.Count - 1)]
        Write-Info "Starting from test: $StartFrom (skipping $startIndex test(s))"
    }

    Write-Info ""
    Write-Info "Found $($testDirs.Count) test(s) to run"

    $results = @()
    foreach ($test in $testDirs) {
        $result = Invoke-TerraformTest -TestPath $test.Path -TestName $test.Name
        $results += $result
    }

    Write-Info ""
    Write-Info ("=" * 60)
    Write-Info "TEST SUMMARY"
    Write-Info ("=" * 60)

    $passed = @($results | Where-Object { $_.Status -eq "Passed" }).Count
    $failed = @($results | Where-Object { $_.Status -eq "Failed" }).Count
    $total = $results.Count

    foreach ($result in $results) {
        if ($result.Status -eq "Passed") {
            Write-Success "  [PASS] $($result.Name) ($([math]::Round($result.TotalTime, 1))s)"
        }
        else {
            Write-Failure "  [FAIL] $($result.Name): $($result.Error)"
        }
    }

    Write-Info ""
    if ($failed -eq 0) {
        Write-Success "All $total tests passed!"
        exit 0
    }
    else {
        Write-Failure "$passed/$total tests passed, $failed failed"
        exit 1
    }
}

Main
