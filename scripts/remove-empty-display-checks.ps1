#!/usr/bin/env pwsh
# Script to remove empty if/else blocks left over from display_name removal

param(
    [Parameter(Mandatory=$true)]
    [string]$Path
)

$content = Get-Content $Path -Raw
$originalLength = $content.Length

# Pattern: Match empty if/else blocks with Display checks
# The files use LF (\n) line endings, not CRLF
$pattern = '\n\t// Display.*?\n\tif\s+\w+\.(?:Get)?Display\(\)?\s*!=\s*""\s*\{\n\t\}\s*else\s*\{\n\t\}\n'
$content = $content -replace $pattern, "`n"

# Count changes
if ($content.Length -ne $originalLength) {
    Set-Content $Path -Value $content -NoNewline
    $linesBefore = ($originalLength -split "`n").Count
    $linesAfter = ($content.Length -split "`n").Count
    $removed = $linesBefore - $linesAfter
    Write-Host "Removed $removed lines from $(Split-Path $Path -Leaf)"
} else {
    # Silently skip files with no matches
}
