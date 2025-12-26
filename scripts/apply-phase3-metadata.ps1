# Apply Phase 3 CommonMetadataAttributes helper to resources
# This script adds maps import and replaces "tags" and "custom_fields" lines with CommonMetadataAttributes()

$resources = @(
    "circuit_type",
    "cluster_group",
    "cluster_type",
    "rir",
    "region"
)

foreach ($res in $resources) {
    $file = "c:\GitRoot\terraform-provider-netbox\internal\resources\${res}_resource.go"
    Write-Host "Processing $file..."

    # Read the file
    $content = Get-Content $file -Raw

    # Step 1: Add maps import if not present
    if ($content -notmatch 'import \([^)]*"maps"') {
        Write-Host "  Adding maps import..."
        $content = $content -replace '(import \(\s*"context"\s*"fmt"\s*)', '$1"maps"`n`t'
    }

    # Step 2: Replace tags and custom_fields with CommonMetadataAttributes
    Write-Host "  Replacing tags and custom_fields..."

    # Pattern: find the two lines with TagsAttribute and CustomFieldsAttribute and the closing braces
    $pattern = '(\t\t)"tags":\s+nbschema\.TagsAttribute\(\),\s+"custom_fields":\s+nbschema\.CustomFieldsAttribute\(\),\s+\},\s+\}\s+\}'
    $replacement = '$1},`n`t}`n`n`t// Add common metadata attributes (tags, custom_fields)`n`tmaps.Copy(resp.Schema.Attributes, nbschema.CommonMetadataAttributes())`n}'

    $content = $content -replace $pattern, $replacement

    # Save the file
    $content | Set-Content $file -NoNewline
    Write-Host "  Done!"
}

Write-Host "`nAll resources updated. Running go build..."
Set-Location "c:\GitRoot\terraform-provider-netbox"
go build .
