<#
.SYNOPSIS
    Runs acceptance tests against a local Netbox Docker instance.

.DESCRIPTION
    This script starts the Docker Compose environment, waits for Netbox to be healthy,
    and then runs the Terraform provider acceptance tests.

.PARAMETER StartOnly
    Only start the Docker environment, don't run tests.

.PARAMETER StopOnly
    Only stop and clean up the Docker environment.

.PARAMETER SkipDocker
    Skip Docker operations and just run tests (assumes Netbox is already running).

.PARAMETER Timeout
    Timeout in seconds to wait for Netbox to be ready. Default is 300 (5 minutes).

.PARAMETER TestPattern
    Optional test pattern to run specific tests. E.g., "TestAccSite" to run only site tests.

.EXAMPLE
    .\scripts\run-acceptance-tests.ps1
    # Starts Docker, waits for Netbox, runs all tests, keeps Docker running

.EXAMPLE
    .\scripts\run-acceptance-tests.ps1 -TestPattern "TestAccSite"
    # Runs only tests matching "TestAccSite"

.EXAMPLE
    .\scripts\run-acceptance-tests.ps1 -StopOnly
    # Stops and removes all Docker containers and volumes
#>

param(
    [switch]$StartOnly,
    [switch]$StopOnly,
    [switch]$SkipDocker,
    [int]$Timeout = 300,
    [string]$TestPattern = ""
)

$ErrorActionPreference = "Stop"

# Configuration
$NetboxUrl = "http://localhost:8000"
$ApiToken = "0123456789abcdef0123456789abcdef01234567"
$ProjectRoot = Split-Path -Parent (Split-Path -Parent $MyInvocation.MyCommand.Path)

function Write-Status {
    param([string]$Message, [string]$Color = "Cyan")
    Write-Host "[$((Get-Date).ToString('HH:mm:ss'))] " -NoNewline -ForegroundColor Gray
    Write-Host $Message -ForegroundColor $Color
}

function Start-NetboxEnvironment {
    Write-Status "Starting Netbox Docker environment..." "Yellow"
    
    Push-Location $ProjectRoot
    try {
        docker-compose up -d
        if ($LASTEXITCODE -ne 0) {
            throw "Failed to start Docker Compose environment"
        }
        Write-Status "Docker containers started" "Green"
    }
    finally {
        Pop-Location
    }
}

function Stop-NetboxEnvironment {
    Write-Status "Stopping Netbox Docker environment..." "Yellow"
    
    Push-Location $ProjectRoot
    try {
        docker-compose down -v
        if ($LASTEXITCODE -ne 0) {
            Write-Status "Warning: Failed to stop Docker Compose environment" "Yellow"
        }
        else {
            Write-Status "Docker environment stopped and cleaned up" "Green"
        }
    }
    finally {
        Pop-Location
    }
}

function Wait-ForNetbox {
    param([int]$TimeoutSeconds)
    
    Write-Status "Waiting for Netbox to be ready (timeout: ${TimeoutSeconds}s)..." "Yellow"
    
    $stopwatch = [System.Diagnostics.Stopwatch]::StartNew()
    $ready = $false
    $lastError = ""
    
    while ($stopwatch.Elapsed.TotalSeconds -lt $TimeoutSeconds) {
        try {
            $response = Invoke-RestMethod -Uri "$NetboxUrl/api/" -Headers @{
                "Authorization" = "Token $ApiToken"
                "Accept" = "application/json"
            } -TimeoutSec 5 -ErrorAction Stop
            
            if ($response) {
                $ready = $true
                break
            }
        }
        catch {
            $lastError = $_.Exception.Message
        }
        
        $elapsed = [math]::Round($stopwatch.Elapsed.TotalSeconds)
        Write-Host "`r  Waiting... ($elapsed/${TimeoutSeconds}s) " -NoNewline
        Start-Sleep -Seconds 5
    }
    
    Write-Host "" # New line after progress
    
    if (-not $ready) {
        throw "Netbox did not become ready within $TimeoutSeconds seconds. Last error: $lastError"
    }
    
    Write-Status "Netbox is ready!" "Green"
}

function Test-DockerAvailable {
    try {
        docker --version | Out-Null
        docker-compose --version | Out-Null
        return $true
    }
    catch {
        return $false
    }
}

function Run-AcceptanceTests {
    param([string]$Pattern)
    
    Write-Status "Running acceptance tests..." "Yellow"
    
    # Set environment variables
    $env:NETBOX_SERVER_URL = $NetboxUrl
    $env:NETBOX_API_TOKEN = $ApiToken
    $env:TF_ACC = "1"
    
    Push-Location $ProjectRoot
    try {
        if ($Pattern) {
            Write-Status "Running tests matching: $Pattern" "Cyan"
            go test ./... -v -timeout 120m -run $Pattern
        }
        else {
            go test ./... -v -timeout 120m
        }
        
        $testResult = $LASTEXITCODE
        
        if ($testResult -eq 0) {
            Write-Status "All tests passed!" "Green"
        }
        else {
            Write-Status "Some tests failed (exit code: $testResult)" "Red"
        }
        
        return $testResult
    }
    finally {
        Pop-Location
    }
}

# Main execution
Write-Host ""
Write-Host "========================================" -ForegroundColor Cyan
Write-Host " Terraform Provider Netbox - Test Runner" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

# Handle stop-only mode
if ($StopOnly) {
    Stop-NetboxEnvironment
    exit 0
}

# Check Docker is available (unless skipping)
if (-not $SkipDocker) {
    if (-not (Test-DockerAvailable)) {
        Write-Status "Docker or docker-compose not found. Please install Docker Desktop." "Red"
        exit 1
    }
}

# Start Docker environment
if (-not $SkipDocker) {
    Start-NetboxEnvironment
    Wait-ForNetbox -TimeoutSeconds $Timeout
}

# Handle start-only mode
if ($StartOnly) {
    Write-Status "Docker environment is running. Netbox is available at: $NetboxUrl" "Green"
    Write-Status "API Token: $ApiToken" "Cyan"
    Write-Status "Run tests manually with: `$env:TF_ACC='1'; go test ./... -v" "Cyan"
    exit 0
}

# Run tests
$exitCode = Run-AcceptanceTests -Pattern $TestPattern

Write-Host ""
Write-Status "Docker environment is still running. To stop: .\scripts\run-acceptance-tests.ps1 -StopOnly" "Cyan"
Write-Host ""

exit $exitCode
