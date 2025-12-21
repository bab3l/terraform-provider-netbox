# Script to update all datasource unit tests with standard pattern

param(
    [string]$DatasourcesDir = "c:\GitRoot\terraform-provider-netbox\internal\datasources",
    [string]$TestsDir = "c:\GitRoot\terraform-provider-netbox\internal\datasources_unit_tests"
)

# Function to convert snake_case to proper CamelCase (handling acronyms)
function ConvertTo-CamelCase {
    param([string]$str)

    # Special cases for acronyms
    $acronyms = @{
        "asn" = "ASN"
        "ip" = "IP"
        "vm" = "VM"
        "id" = "ID"
        "rir" = "RIR"
        "iam" = "IAM"
        "api" = "API"
        "url" = "URL"
        "dcim" = "DCIM"
        "vlan" = "VLAN"
        "vrf" = "VRF"
        "lan" = "LAN"
    }

    $parts = $str.Split('_')
    $camel = @()

    foreach ($part in $parts) {
        if ($acronyms.ContainsKey($part)) {
            $camel += $acronyms[$part]
        } else {
            $camel += $part.Substring(0, 1).ToUpper() + $part.Substring(1).ToLower()
        }
    }

    return $camel -join ''
}

# Function to extract schema attributes from datasource file
function Get-SchemaAttributes {
    param([string]$filePath)

    $content = Get-Content -Path $filePath -Raw

    # Find the Schema method
    $schemaStart = $content.IndexOf('Attributes: map[string]schema.Attribute{')
    if ($schemaStart -eq -1) {
        return @{ lookup = @(); computed = @() }
    }

    # Find the closing brace of the Attributes map
    $level = 0
    $startIdx = $schemaStart + 'Attributes: map[string]schema.Attribute{'.Length
    $idx = $startIdx

    while ($idx -lt $content.Length) {
        if ($content[$idx] -eq '{') { $level++ }
        elseif ($content[$idx] -eq '}') {
            $level--
            if ($level -eq 0) { break }
        }
        $idx++
    }

    $schemaContent = $content.Substring($startIdx, $idx - $startIdx)

    # Split by attribute definitions
    $attrBlocks = $schemaContent -split '"[a-z_]+":\s*schema\.' | Select-Object -Skip 1

    $lookup = @()
    $computed = @()
    $currentIdx = $startIdx

    # Extract attribute names from the original content
    $attrNames = [regex]::Matches($schemaContent, '"([a-z_]+)":\s*schema\.', [System.Text.RegularExpressions.RegexOptions]::IgnoreCase)

    $blockIdx = 0
    foreach ($attrName in $attrNames) {
        $name = $attrName.Groups[1].Value

        if ($blockIdx -lt $attrBlocks.Count) {
            $block = $attrBlocks[$blockIdx]

            # Check for Optional: true
            $hasOptional = $block -match 'Optional:\s*true'
            # Check for Computed: true
            $hasComputed = $block -match 'Computed:\s*true'

            if ($hasOptional) {
                $lookup += $name
            }
            if ($hasComputed -and -not $hasOptional) {
                $computed += $name
            } elseif ($hasComputed -and $hasOptional) {
                # If it has both Optional and Computed, it's a lookup field
                # Keep it in lookup only
            }
        }
        $blockIdx++
    }

    return @{ lookup = $lookup; computed = $computed }
}

# Process each datasource
$testFiles = Get-ChildItem -Path $TestsDir -Filter "*_data_source_test.go" | Sort-Object Name

Write-Host "Updating $($testFiles.Count) datasource unit tests..."

foreach ($testFile in $testFiles) {
    $baseName = $testFile.Name -replace '_data_source_test\.go', ''
    $dataSourceName = ConvertTo-CamelCase -str $baseName
    $dataSourceTypeName = "netbox_$baseName"

    # Get the corresponding datasource implementation file
    $implFile = Join-Path -Path $DatasourcesDir -ChildPath "$baseName`_data_source.go"

    if (-not (Test-Path $implFile)) {
        Write-Warning "Implementation file not found: $implFile"
        continue
    }

    # Extract schema attributes
    $attrs = Get-SchemaAttributes -filePath $implFile

    # Create the test content
    $lookupStr = if ($attrs.lookup.Count -gt 0) {
        "`"" + ($attrs.lookup -join "`", `"") + "`""
    } else {
        ""
    }

    $computedAttrs = $attrs.computed | Where-Object { $_ -notin @($attrs.lookup) }
    $computedStr = if ($computedAttrs.Count -gt 0) {
        "`"" + ($computedAttrs -join "`", `"") + "`""
    } else {
        ""
    }

    $testContent = @"
package datasources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/datasources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func Test${dataSourceName}DataSourceSchema(t *testing.T) {
	d := datasources.New${dataSourceName}DataSource()

	req := datasource.SchemaRequest{}
	resp := &datasource.SchemaResponse{}

	d.Schema(context.Background(), req, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("Schema returned errors: %v", resp.Diagnostics)
	}

	testutil.ValidateDataSourceSchema(t, resp.Schema.Attributes, testutil.DataSourceValidation{
		LookupAttrs: []string{$lookupStr},
		ComputedAttrs: []string{$computedStr},
	})
}

func Test${dataSourceName}DataSourceMetadata(t *testing.T) {
	d := datasources.New${dataSourceName}DataSource()
	testutil.ValidateDataSourceMetadata(t, d, "netbox", "$dataSourceTypeName")
}

func Test${dataSourceName}DataSourceConfigure(t *testing.T) {
	d := datasources.New${dataSourceName}DataSource()
	testutil.ValidateDataSourceConfigure(t, d)
}
"@

    # Write the test file
    $testPath = $testFile.FullName
    Set-Content -Path $testPath -Value $testContent -Encoding UTF8
    Write-Host "Updated: $($testFile.Name)"
}

Write-Host "Done! Updated $($testFiles.Count) test files."
