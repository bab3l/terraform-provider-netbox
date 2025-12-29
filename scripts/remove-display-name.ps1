#!/usr/bin/env pwsh
# Script to remove display_name field from resource files
# This removes:
# 1. DisplayName types.String from model struct
# 2. "display_name": nbschema.DisplayNameAttribute(...) from Schema()
# 3. data.DisplayName = ... assignments from mapTo functions

param(
    [Parameter(Mandatory=$true)]
    [string]$Path
)

$content = Get-Content $Path -Raw

# Count occurrences before
$beforeCount = ([regex]::Matches($content, "DisplayName")).Count

if ($beforeCount -eq 0) {
    Write-Host "No DisplayName references found in $Path"
    return
}

# 1. Remove DisplayName from model struct
$content = $content -replace '(?m)^\s*DisplayName\s+types\.String\s+`tfsdk:"display_name"`\r?\n', ''

# 2. Remove display_name from schema
$content = $content -replace '(?m)^\s*"display_name":\s*nbschema\.DisplayNameAttribute\([^)]+\),?\r?\n', ''

# 3. Remove DisplayName assignment blocks (various patterns)
# Pattern 1: Simple assignment
$content = $content -replace '(?m)^\s*data\.DisplayName\s*=\s*types\.String.*\r?\n', ''

# Pattern 2: If block with Display field
$content = $content -replace '(?s)(\r?\n)\s*//\s*DisplayName.*?(?:data\.DisplayName\s*=\s*types\.StringNull\(\).*?\r?\n)', '$1'

# Count occurrences after
$afterCount = ([regex]::Matches($content, "DisplayName")).Count

# Save if changes were made
if ($beforeCount -ne $afterCount) {
    $content | Set-Content $Path -NoNewline
    $removed = $beforeCount - $afterCount
    Write-Host "Removed $removed DisplayName references from $Path"
} else {
    Write-Host "Found $beforeCount DisplayName references in $Path but could not remove them"
}
