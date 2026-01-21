# Import Behavior Test Script
# Creates resources directly in NetBox via API, then tests Terraform import behavior
# This simulates real-world import scenarios where resources exist before Terraform management

param(
    [switch]$SkipCleanup,
    [switch]$Verbose
)

$ErrorActionPreference = "Stop"
$ScriptRoot = Split-Path -Parent $MyInvocation.MyCommand.Path
$ProjectRoot = Split-Path -Parent $ScriptRoot

# Configuration
$NetboxUrl = $env:NETBOX_SERVER_URL
if (-not $NetboxUrl) { $NetboxUrl = "http://localhost:8000" }

$ApiToken = $env:NETBOX_API_TOKEN
if (-not $ApiToken) { $ApiToken = "0123456789abcdef0123456789abcdef01234567" }

$TestDir = Join-Path $ProjectRoot "test\import-behavior"

function Write-Success { param($msg) Write-Host $msg -ForegroundColor Green }
function Write-Failure { param($msg) Write-Host $msg -ForegroundColor Red }
function Write-Info { param($msg) Write-Host $msg -ForegroundColor Cyan }
function Write-Warn { param($msg) Write-Host $msg -ForegroundColor Yellow }
function Write-Detail { param($msg) if ($Verbose) { Write-Host $msg -ForegroundColor Gray } }

# Generate unique test suffix
$TestSuffix = -join ((97..122) | Get-Random -Count 8 | ForEach-Object { [char]$_ })
Write-Info "Test suffix: $TestSuffix"

# API Helper Functions
function Invoke-NetboxApi {
    param(
        [string]$Method,
        [string]$Endpoint,
        [hashtable]$Body = $null
    )

    $headers = @{
        "Authorization" = "Token $ApiToken"
        "Content-Type" = "application/json"
        "Accept" = "application/json"
    }

    $uri = "$NetboxUrl$Endpoint"
    Write-Detail "API $Method $uri"

    $params = @{
        Method = $Method
        Uri = $uri
        Headers = $headers
    }

    if ($Body) {
        $params.Body = ($Body | ConvertTo-Json -Depth 10)
        Write-Detail "Body: $($params.Body)"
    }

    try {
        $response = Invoke-RestMethod @params
        return $response
    }
    catch {
        Write-Failure "API Error: $_"
        if ($_.Exception.Response) {
            $reader = New-Object System.IO.StreamReader($_.Exception.Response.GetResponseStream())
            $reader.BaseStream.Position = 0
            $responseBody = $reader.ReadToEnd()
            Write-Failure "Response: $responseBody"
        }
        throw
    }
}

function New-NetboxResource {
    param(
        [string]$Endpoint,
        [hashtable]$Body
    )
    return Invoke-NetboxApi -Method "POST" -Endpoint $Endpoint -Body $Body
}

function Remove-NetboxResource {
    param(
        [string]$Endpoint,
        [int]$Id
    )
    try {
        Invoke-NetboxApi -Method "DELETE" -Endpoint "$Endpoint$Id/"
        return $true
    }
    catch {
        Write-Warn "Failed to delete resource at $Endpoint$Id/: $_"
        return $false
    }
}

# Store created resource IDs for cleanup
$CreatedResources = @()

function Register-Cleanup {
    param(
        [string]$Endpoint,
        [int]$Id
    )
    $script:CreatedResources += @{ Endpoint = $Endpoint; Id = $Id }
}

function Invoke-Cleanup {
    Write-Info "`nCleaning up NetBox resources..."
    # Delete in reverse order (dependencies first)
    for ($i = $CreatedResources.Count - 1; $i -ge 0; $i--) {
        $resource = $CreatedResources[$i]
        Write-Detail "Deleting $($resource.Endpoint)$($resource.Id)/"
        Remove-NetboxResource -Endpoint $resource.Endpoint -Id $resource.Id | Out-Null
    }

    # Clean up test directory
    if (Test-Path $TestDir) {
        Write-Detail "Removing test directory: $TestDir"
        Remove-Item -Path $TestDir -Recurse -Force
    }
}

# ============================================================================
# STEP 1: Create resources directly in NetBox via API
# ============================================================================
Write-Info "`n=== Step 1: Creating resources in NetBox via API ==="

# Create Cluster Type
Write-Info "Creating cluster type..."
$clusterType = New-NetboxResource -Endpoint "/api/virtualization/cluster-types/" -Body @{
    name = "import-test-cluster-type-$TestSuffix"
    slug = "import-test-cluster-type-$TestSuffix"
}
Register-Cleanup -Endpoint "/api/virtualization/cluster-types/" -Id $clusterType.id
Write-Success "  Created cluster type: ID=$($clusterType.id), Name=$($clusterType.name)"

# Create Cluster
Write-Info "Creating cluster..."
$cluster = New-NetboxResource -Endpoint "/api/virtualization/clusters/" -Body @{
    name = "import-test-cluster-$TestSuffix"
    type = $clusterType.id
}
Register-Cleanup -Endpoint "/api/virtualization/clusters/" -Id $cluster.id
Write-Success "  Created cluster: ID=$($cluster.id), Name=$($cluster.name)"

# Create Device Role (for VM role)
Write-Info "Creating device role..."
$deviceRole = New-NetboxResource -Endpoint "/api/dcim/device-roles/" -Body @{
    name = "import-test-role-$TestSuffix"
    slug = "import-test-role-$TestSuffix"
    color = "ff0000"
    vm_role = $true
}
Register-Cleanup -Endpoint "/api/dcim/device-roles/" -Id $deviceRole.id
Write-Success "  Created device role: ID=$($deviceRole.id), Name=$($deviceRole.name)"

# Create Tenant
Write-Info "Creating tenant..."
$tenant = New-NetboxResource -Endpoint "/api/tenancy/tenants/" -Body @{
    name = "import-test-tenant-$TestSuffix"
    slug = "import-test-tenant-$TestSuffix"
}
Register-Cleanup -Endpoint "/api/tenancy/tenants/" -Id $tenant.id
Write-Success "  Created tenant: ID=$($tenant.id), Name=$($tenant.name)"

# Create Platform
Write-Info "Creating platform..."
$platform = New-NetboxResource -Endpoint "/api/dcim/platforms/" -Body @{
    name = "import-test-platform-$TestSuffix"
    slug = "import-test-platform-$TestSuffix"
}
Register-Cleanup -Endpoint "/api/dcim/platforms/" -Id $platform.id
Write-Success "  Created platform: ID=$($platform.id), Name=$($platform.name)"

# Create Virtual Machine with all reference fields
Write-Info "Creating virtual machine..."
$vm = New-NetboxResource -Endpoint "/api/virtualization/virtual-machines/" -Body @{
    name = "import-test-vm-$TestSuffix"
    cluster = $cluster.id
    role = $deviceRole.id
    tenant = $tenant.id
    platform = $platform.id
    status = "active"
    vcpus = 2
    memory = 2048
    disk = 50
}
Register-Cleanup -Endpoint "/api/virtualization/virtual-machines/" -Id $vm.id
Write-Success "  Created VM: ID=$($vm.id), Name=$($vm.name)"

Write-Info "`nCreated resource IDs:"
Write-Info "  Cluster Type: $($clusterType.id)"
Write-Info "  Cluster:      $($cluster.id)"
Write-Info "  Device Role:  $($deviceRole.id)"
Write-Info "  Tenant:       $($tenant.id)"
Write-Info "  Platform:     $($platform.id)"
Write-Info "  VM:           $($vm.id)"

# ============================================================================
# STEP 2: Create Terraform configuration using ID references
# ============================================================================
Write-Info "`n=== Step 2: Creating Terraform configuration ==="

# Create test directory
if (-not (Test-Path $TestDir)) {
    New-Item -Path $TestDir -ItemType Directory -Force | Out-Null
}

# Create Terraform config that uses IDs for reference fields
$tfConfig = @"
terraform {
  required_providers {
    netbox = {
      source = "bab3l/netbox"
    }
  }
}

provider "netbox" {
  server_url = "$NetboxUrl"
  api_token  = "$ApiToken"
}

# Reference resources (already exist in NetBox, we'll just reference them)
# These are created so Terraform knows the IDs to use in the VM resource

resource "netbox_cluster_type" "test" {
  name = "$($clusterType.name)"
  slug = "$($clusterType.slug)"
}

resource "netbox_cluster" "test" {
  name = "$($cluster.name)"
  type = netbox_cluster_type.test.id
}

resource "netbox_device_role" "test" {
  name    = "$($deviceRole.name)"
  slug    = "$($deviceRole.slug)"
  color   = "ff0000"
  vm_role = true
}

resource "netbox_tenant" "test" {
  name = "$($tenant.name)"
  slug = "$($tenant.slug)"
}

resource "netbox_platform" "test" {
  name = "$($platform.name)"
  slug = "$($platform.slug)"
}

# The VM resource - this is what we'll import
# Notice: We use ID references (resource.id) for all reference fields
resource "netbox_virtual_machine" "test" {
  name     = "$($vm.name)"
  cluster  = netbox_cluster.test.id      # ID reference
  role     = netbox_device_role.test.id  # ID reference
  tenant   = netbox_tenant.test.id       # ID reference
  platform = netbox_platform.test.id     # ID reference
  status   = "active"
  vcpus    = 2
  memory   = 2048
  disk     = 50
}

output "vm_id" {
  value = netbox_virtual_machine.test.id
}

output "vm_cluster" {
  value = netbox_virtual_machine.test.cluster
}

output "vm_role" {
  value = netbox_virtual_machine.test.role
}

output "vm_tenant" {
  value = netbox_virtual_machine.test.tenant
}

output "vm_platform" {
  value = netbox_virtual_machine.test.platform
}
"@

$tfConfigPath = Join-Path $TestDir "main.tf"
Set-Content -Path $tfConfigPath -Value $tfConfig
Write-Success "Created Terraform config at: $tfConfigPath"

# ============================================================================
# STEP 3: Initialize Terraform and import resources
# ============================================================================
Write-Info "`n=== Step 3: Initializing Terraform ==="

Push-Location $TestDir
try {
    # Initialize
    Write-Info "Running terraform init..."
    $env:TF_LOG = ""  # Disable terraform debug logging
    terraform init -upgrade -no-color *>&1 | Out-Null
    if ($LASTEXITCODE -ne 0) {
        Write-Failure "Terraform init failed"
        throw "Terraform init failed"
    }
    Write-Success "Terraform initialized successfully"

    # ============================================================================
    # STEP 4: Import the supporting resources first (so Terraform knows the IDs)
    # ============================================================================
    Write-Info "`n=== Step 4: Importing supporting resources ==="

    Write-Info "Importing cluster type..."
    terraform import netbox_cluster_type.test $($clusterType.id) 2>&1 | Out-Null

    Write-Info "Importing cluster..."
    terraform import netbox_cluster.test $($cluster.id) 2>&1 | Out-Null

    Write-Info "Importing device role..."
    terraform import netbox_device_role.test $($deviceRole.id) 2>&1 | Out-Null

    Write-Info "Importing tenant..."
    terraform import netbox_tenant.test $($tenant.id) 2>&1 | Out-Null

    Write-Info "Importing platform..."
    terraform import netbox_platform.test $($platform.id) 2>&1 | Out-Null

    Write-Success "Supporting resources imported"

    # ============================================================================
    # STEP 5: Import the VM (the resource with reference fields)
    # ============================================================================
    Write-Info "`n=== Step 5: Importing virtual machine ==="

    Write-Info "Importing VM with ID: $($vm.id)"
    $importOutput = terraform import netbox_virtual_machine.test $($vm.id) 2>&1
    Write-Detail $importOutput

    if ($LASTEXITCODE -ne 0) {
        Write-Failure "VM import failed:"
        Write-Host $importOutput
        throw "VM import failed"
    }
    Write-Success "VM imported successfully"

    # ============================================================================
    # STEP 6: Run terraform plan to check for diffs
    # ============================================================================
    Write-Info "`n=== Step 6: Checking for plan differences ==="

    Write-Info "Running terraform plan..."
    $planOutput = terraform plan -detailed-exitcode 2>&1
    $planExitCode = $LASTEXITCODE

    Write-Info "`nPlan output:"
    Write-Host "----------------------------------------"
    Write-Host $planOutput
    Write-Host "----------------------------------------"

    # Check plan result
    switch ($planExitCode) {
        0 {
            Write-Success "`n✅ SUCCESS: No changes detected after import!"
            Write-Success "   The provider correctly handles reference field import."
        }
        1 {
            Write-Failure "`n❌ ERROR: Terraform plan failed"
        }
        2 {
            Write-Warn "`n⚠️  CHANGES DETECTED: Plan shows differences after import!"
            Write-Warn "   This indicates the import behavior issue exists."
            Write-Info "`nLet's see the state values vs config values:"

            # Show state
            Write-Info "`nCurrent state for VM:"
            terraform state show netbox_virtual_machine.test 2>&1 | ForEach-Object { Write-Host "  $_" }
        }
    }

    # ============================================================================
    # STEP 7: Show detailed state information
    # ============================================================================
    Write-Info "`n=== Step 7: State Analysis ==="

    Write-Info "VM state after import:"
    $stateOutput = terraform state show netbox_virtual_machine.test 2>&1
    Write-Host $stateOutput

    Write-Info "`nExpected ID values:"
    Write-Info "  cluster  = $($cluster.id)"
    Write-Info "  role     = $($deviceRole.id)"
    Write-Info "  tenant   = $($tenant.id)"
    Write-Info "  platform = $($platform.id)"

    # Parse actual values from state
    $clusterMatch = $stateOutput | Select-String -Pattern 'cluster\s*=\s*"([^"]+)"'
    $roleMatch = $stateOutput | Select-String -Pattern 'role\s*=\s*"([^"]+)"'
    $tenantMatch = $stateOutput | Select-String -Pattern 'tenant\s*=\s*"([^"]+)"'
    $platformMatch = $stateOutput | Select-String -Pattern 'platform\s*=\s*"([^"]+)"'

    Write-Info "`nActual state values:"
    if ($clusterMatch) { Write-Info "  cluster  = $($clusterMatch.Matches.Groups[1].Value)" }
    if ($roleMatch) { Write-Info "  role     = $($roleMatch.Matches.Groups[1].Value)" }
    if ($tenantMatch) { Write-Info "  tenant   = $($tenantMatch.Matches.Groups[1].Value)" }
    if ($platformMatch) { Write-Info "  platform = $($platformMatch.Matches.Groups[1].Value)" }

    # Check if values are IDs or names
    Write-Info "`nValue type analysis:"
    $hasNameValues = $false

    if ($clusterMatch) {
        $val = $clusterMatch.Matches.Groups[1].Value
        if ($val -match '^\d+$') {
            Write-Success "  cluster:  ID ($val)"
        } else {
            Write-Warn "  cluster:  NAME ($val) - Expected ID: $($cluster.id)"
            $hasNameValues = $true
        }
    }

    if ($roleMatch) {
        $val = $roleMatch.Matches.Groups[1].Value
        if ($val -match '^\d+$') {
            Write-Success "  role:     ID ($val)"
        } else {
            Write-Warn "  role:     NAME ($val) - Expected ID: $($deviceRole.id)"
            $hasNameValues = $true
        }
    }

    if ($tenantMatch) {
        $val = $tenantMatch.Matches.Groups[1].Value
        if ($val -match '^\d+$') {
            Write-Success "  tenant:   ID ($val)"
        } else {
            Write-Warn "  tenant:   NAME ($val) - Expected ID: $($tenant.id)"
            $hasNameValues = $true
        }
    }

    if ($platformMatch) {
        $val = $platformMatch.Matches.Groups[1].Value
        if ($val -match '^\d+$') {
            Write-Success "  platform: ID ($val)"
        } else {
            Write-Warn "  platform: NAME ($val) - Expected ID: $($platform.id)"
            $hasNameValues = $true
        }
    }

    if ($hasNameValues) {
        Write-Failure "`n❌ ISSUE CONFIRMED: Some reference fields have NAME values instead of IDs!"
        Write-Failure "   This would cause unnecessary plan diffs when config uses ID references."
    } else {
        Write-Success "`n✅ All reference fields have ID values - import behavior is correct!"
    }

}
finally {
    Pop-Location

    # Cleanup
    if (-not $SkipCleanup) {
        Invoke-Cleanup
    } else {
        Write-Warn "`nSkipping cleanup. Test resources remain in NetBox."
        Write-Warn "Test directory: $TestDir"
        Write-Warn "To clean up manually, delete resources with suffix: $TestSuffix"
    }
}

Write-Info "`n=== Test Complete ==="
