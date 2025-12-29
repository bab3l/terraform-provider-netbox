# Fix-Formatting.ps1
# Removes excessive blank lines from Go source files

param(
    [Parameter(Mandatory=$true)]
    [string]$Path,

    [switch]$DryRun
)

function Fix-GoFile {
    param([string]$FilePath)

    Write-Host "Processing: $FilePath"

    $content = Get-Content $FilePath -Raw
    $original = $content

    # Remove blank lines immediately after opening braces
    $content = $content -creplace '(?m)(\{)\r?\n\r?\n', ('$1' + [Environment]::NewLine)

    # Remove blank lines immediately before closing braces
    $content = $content -creplace '(?m)\r?\n\r?\n(\s*\})', ([Environment]::NewLine + '$1')

    # Remove excessive blank lines in comment blocks
    $content = $content -creplace '(?m)(//[^\n]*)\r?\n\r?\n(\s*//)', ('$1' + [Environment]::NewLine + '$2')

    # Reduce multiple blank lines (3+) to max 2
    $content = $content -creplace '(\r?\n){4,}', ([Environment]::NewLine + [Environment]::NewLine + [Environment]::NewLine)

    if ($content -ne $original) {
        $linesBefore = ($original -split "\n").Count
        $linesAfter = ($content -split "\n").Count
        $blankBefore = (($original -split "\n") | Where-Object { $_.Trim() -eq "" }).Count
        $blankAfter = (($content -split "\n") | Where-Object { $_.Trim() -eq "" }).Count

        $pctBefore = if ($linesBefore -gt 0) { [math]::Round($blankBefore/$linesBefore*100) } else { 0 }
        $pctAfter = if ($linesAfter -gt 0) { [math]::Round($blankAfter/$linesAfter*100) } else { 0 }

        Write-Host "  Before: $linesBefore lines, $blankBefore blank ($pctBefore%)" -ForegroundColor Yellow
        Write-Host "  After:  $linesAfter lines, $blankAfter blank ($pctAfter%)" -ForegroundColor Green

        if (-not $DryRun) {
            Set-Content -Path $FilePath -Value $content -NoNewline
            Write-Host "  Saved" -ForegroundColor Green
        } else {
            Write-Host "  [DRY RUN] Would save changes" -ForegroundColor Cyan
        }

        return $true
    } else {
        Write-Host "  No changes needed" -ForegroundColor Gray
        return $false
    }
}

# Resolve the path
$files = Get-Item $Path -ErrorAction SilentlyContinue

if (-not $files) {
    Write-Error "No files found matching: $Path"
    exit 1
}

$changedCount = 0
$totalCount = 0

foreach ($file in $files) {
    if ($file.Extension -eq ".go") {
        $totalCount++
        if (Fix-GoFile -FilePath $file.FullName) {
            $changedCount++
        }
    }
}

Write-Host ""
Write-Host "========================================"  -ForegroundColor Cyan
Write-Host "Summary: $changedCount of $totalCount files modified" -ForegroundColor Cyan
if ($DryRun) {
    Write-Host "DRY RUN - No files were actually modified" -ForegroundColor Yellow
    Write-Host "Run without -DryRun to apply changes" -ForegroundColor Yellow
}
Write-Host "========================================" -ForegroundColor Cyan

# Run gofmt on all modified files if not a dry run
if (-not $DryRun -and $changedCount -gt 0) {
    Write-Host ""
    Write-Host "Running gofmt -s -w on modified files..." -ForegroundColor Cyan
    foreach ($file in $files) {
        if ($file.Extension -eq ".go") {
            gofmt -s -w $file.FullName
        }
    }
    Write-Host "gofmt complete" -ForegroundColor Green
}
