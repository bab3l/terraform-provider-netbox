$ErrorActionPreference = "Stop"

$packages = go list ./... | Where-Object { $_ -notmatch 'acceptance_tests' }
if (-not $packages -or $packages.Count -eq 0) {
	throw "No packages found to test."
}

& go test @packages -v
