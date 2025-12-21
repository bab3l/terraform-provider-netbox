$filePath = "c:\GitRoot\terraform-provider-netbox\internal\datasources_acceptance_tests\acceptance_test.go"
$outputDir = "c:\GitRoot\terraform-provider-netbox\internal\datasources_acceptance_tests"

$content = Get-Content -Path $filePath -Raw

# Test function and datasource name mappings
$tests = @(
    @{test = "TestAccSiteDataSource_basic"; ds = "site"; configStart = 61; configEnd = 174},
    @{test = "TestAccTenantDataSource_basic"; ds = "tenant"; configStart = 220; configEnd = 329},
    @{test = "TestAccSiteGroupDataSource_basic"; ds = "site_group"; configStart = 375; configEnd = 484},
    @{test = "TestAccTenantGroupDataSource_basic"; ds = "tenant_group"; configStart = 530; configEnd = 639},
    @{test = "TestAccManufacturerDataSource_basic"; ds = "manufacturer"; configStart = 685; configEnd = 794},
    @{test = "TestAccPlatformDataSource_basic"; ds = "platform"; configStart = 853; configEnd = 986},
    @{test = "TestAccRegionDataSource_basic"; ds = "region"; configStart = 1032; configEnd = 1141},
    @{test = "TestAccLocationDataSource_basic"; ds = "location"; configStart = 1193; configEnd = 1330},
    @{test = "TestAccRackDataSource_basic"; ds = "rack"; configStart = 1378; configEnd = 1511},
    @{test = "TestAccRackRoleDataSource_basic"; ds = "rack_role"; configStart = 1557; configEnd = 1666},
    @{test = "TestAccDeviceRoleDataSource_basic"; ds = "device_role"; configStart = 1714; configEnd = 1823},
    @{test = "TestAccDeviceTypeDataSource_basic"; ds = "device_type"; configStart = 1884; configEnd = 2019},
    @{test = "TestAccRouteTargetDataSource_basic"; ds = "route_target"; configStart = 2061; configEnd = 2168},
    @{test = "TestAccVirtualDiskDataSource_basic"; ds = "virtual_disk"; configStart = 2235; configEnd = 2446},
    @{test = "TestAccASNRangeDataSource_basic"; ds = "asn_range"; configStart = 2507; configEnd = 2650},
    @{test = "TestAccDeviceBayTemplateDataSource_basic"; ds = "device_bay_template"; configStart = 2713; configEnd = 2878}
)

$header = @"
package datasources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/provider"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

"@

# Function to clean up excessive newlines
function Clean-Newlines {
    param([string]$text)
    # Replace 3+ consecutive newlines with 2
    $text = $text -replace "`n`n+", "`n`n"
    # Remove trailing whitespace on each line
    $text = ($text -split "`n" | ForEach-Object { $_.TrimEnd() }) -join "`n"
    return $text.Trim()
}

# Split content by lines
$lines = $content -split "`n"

foreach ($test in $tests) {
    $testName = $test.test
    $dsName = $test.ds
    $configEnd = $test.configEnd

    # Find test start and end lines
    $testStartIdx = [array]::FindIndex($lines, [System.Predicate[string]] { param($line) $line -match "func $testName" })

    # Find the next function (config function is right after test)
    $nextFuncIdx = $testStartIdx + 1
    while ($nextFuncIdx -lt $lines.Count -and $lines[$nextFuncIdx] -notmatch "^func testAcc") {
        $nextFuncIdx++
    }

    # Find where the config function ends (find next func or end of file)
    $configEndIdx = $nextFuncIdx + 1
    while ($configEndIdx -lt $lines.Count -and $lines[$configEndIdx] -notmatch "^func Test|^func testAcc|^$") {
        $configEndIdx++
    }
    $configEndIdx--  # Go back one line

    # Extract test and config
    $testLines = $lines[$testStartIdx..$configEndIdx]
    $testContent = ($testLines -join "`n").TrimEnd()

    # Clean up newlines
    $testContent = Clean-Newlines -text $testContent

    # Create file
    $fileContent = "$header`n$testContent`n"
    $outputFile = Join-Path -Path $outputDir -ChildPath "${dsName}_data_source_test.go"

    Set-Content -Path $outputFile -Value $fileContent -Encoding UTF8
    Write-Host "Created: $outputFile"
}

Write-Host "`nAll test files have been created!"
