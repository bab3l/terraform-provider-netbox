$ErrorActionPreference = "Stop"

if (-not $env:NETBOX_SERVER_URL) { throw "NETBOX_SERVER_URL is not set" }
if (-not $env:NETBOX_API_TOKEN) { throw "NETBOX_API_TOKEN is not set" }

if (-not $env:TF_ACC) {
	$env:TF_ACC = "1"
}

$timeout = if ($env:NETBOX_TEST_TIMEOUT) { $env:NETBOX_TEST_TIMEOUT } else { "120m" }

$ranAny = $false

# Run all acceptance tests EXCEPT customfields packages (safe to run with default parallelism)
$acceptancePkgs = go list ./... | Where-Object { $_ -match 'acceptance_tests' -and $_ -notmatch 'acceptance_tests_customfields' }
if ($acceptancePkgs -and $acceptancePkgs.Count -gt 0) {
	$ranAny = $true
	Write-Host "Running acceptance tests (non-customfields)..." -ForegroundColor Cyan
	& go test @acceptancePkgs -v -timeout $timeout
	if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }
}

# Run customfields acceptance tests SERIAL only (shared resources)
# Note: these packages often contain ONLY files with `//go:build customfields`, so we must
# discover them with the build tag enabled.
$customfieldsPkgs = go list -tags customfields ./... | Where-Object { $_ -match 'acceptance_tests_customfields' }
if ($customfieldsPkgs -and $customfieldsPkgs.Count -gt 0) {
	$ranAny = $true
	Write-Host "Running acceptance tests (customfields, serial)..." -ForegroundColor Cyan
	& go test -tags customfields @customfieldsPkgs -v -timeout $timeout -p 1 -parallel 1
	exit $LASTEXITCODE
}

if (-not $ranAny) {
	Write-Host "No acceptance test packages found." -ForegroundColor Yellow
} else {
	Write-Host "No customfields acceptance test packages found; skipping customfields." -ForegroundColor Yellow
}
