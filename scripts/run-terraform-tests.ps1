# Terraform Integration Test Runner
# Runs all Terraform test configurations against a running Netbox instance

param(
    [string]$TestDir = "",
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

function Test-Environment {
    Write-Info "Checking environment..."
    
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
        $tfVersion = terraform version -json | ConvertFrom-Json
        Write-Success "Terraform version: $($tfVersion.terraform_version)"
    }
    catch {
        Write-Failure "Terraform not found in PATH"
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
        $planOutput = terraform plan -no-color 2>&1
        if ($LASTEXITCODE -ne 0) {
            throw "terraform plan failed: $planOutput"
        }
        Write-Success "  Plan completed"
        
        # Apply
        Write-Info "  terraform apply..."
        $applyOutput = terraform apply -auto-approve -no-color 2>&1
        if ($LASTEXITCODE -ne 0) {
            throw "terraform apply failed: $applyOutput"
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
            $destroyOutput = terraform destroy -auto-approve -no-color 2>&1
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
        
        $testOrder = @("tenant_group", "tenant", "site_group", "site")
        
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
