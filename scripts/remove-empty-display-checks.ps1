#!/usr/bin/env pwsh
# Script to remove empty if/else blocks left over from display_name removal
# Pattern to remove:
#   if <something>.Display != "" {
#   } else {
#   }

param(
    [Parameter(Mandatory=$true)]
    [string]$Path
)

$content = Get-Content $Path -Raw

# Pattern 1: Match with comment above
# \t// Display name\r\n\tif X.Display != "" {\r\n\t} else {\r\n\t}\r\n
$pattern1 = '(?m)^\t// Display.*?\r?\n\tif\s+\w+\.(?:Get)?Display\(\)?\s*!=\s*""\s*\{\r?\n\t\}\s*else\s*\{\r?\n\t\}\r?\n'

# Pattern 2: Without comment
$pattern2 = '(?m)^\tif\s+\w+\.(?:Get)?Display\(\)?\s*!=\s*""\s*\{\r?\n\t\}\s*else\s*\{\r?\n\t\}\r?\n'

# Remove both patterns
$newContent = $content -replace $pattern1, ""
$newContent = $newContent -replace $pattern2, ""

# Count removals
$beforeLines = ($content -split "`n").Count
$afterLines = ($newContent -split "`n").Count
$removed = $beforeLines - $afterLines

if ($removed -gt 0) {
    Set-Content $Path -Value $newContent -NoNewline
    Write-Host "Removed $removed lines from $Path"
} else {
    Write-Host "No empty Display checks found in $Path"
}
