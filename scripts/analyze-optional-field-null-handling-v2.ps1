# Enhanced Optional Field Null Handling Analysis Script
# This script performs comprehensive analysis of all resources

param(
    [switch]$Detailed,
    [switch]$ExportJson
)

Write-Host "`n╔════════════════════════════════════════════════════════════════╗" -ForegroundColor Cyan
Write-Host "║  OPTIONAL FIELD NULL HANDLING - COMPREHENSIVE ANALYSIS        ║" -ForegroundColor Cyan
Write-Host "╚════════════════════════════════════════════════════════════════╝" -ForegroundColor Cyan
Write-Host "Date: $(Get-Date -Format 'yyyy-MM-dd HH:mm:ss')`n" -ForegroundColor Gray

$resourcePath = ".\internal\resources"
$resources = Get-ChildItem -Path "$resourcePath\*_resource.go" | Where-Object { $_.Name -notmatch 'provider_resource' }

Write-Host "Scanning $($resources.Count) resources..." -ForegroundColor Yellow
Write-Host "Looking for optional fields without null handling...`n" -ForegroundColor Gray

$results = @()
$totalResources = 0
$resourcesWithIssues = 0
$totalFields = 0
$fieldsWithIssues = 0

foreach ($file in $resources | Sort-Object Name) {
    $totalResources++
    $resourceName = $file.BaseName -replace '_resource$', ''
    $content = Get-Content $file.FullName -Raw

    # Find all build request functions
    $buildMatches = [regex]::Matches($content, '(?s)func \(r \*\w+Resource\) (build\w+(?:Request|CreateRequest|UpdateRequest))\([^)]*\).*?\{(.*?)(?=\n(?:func \(r \*\w+Resource\)|// map))')

    if ($buildMatches.Count -eq 0) {
        continue
    }

    $resourceIssues = @()

    foreach ($buildMatch in $buildMatches) {
        $funcName = $buildMatch.Groups[1].Value
        $funcBody = $buildMatch.Groups[2].Value

        # Find all if !data.FieldName.IsNull patterns
        $fieldMatches = [regex]::Matches($funcBody, 'if\s+!data\.(\w+)\.IsNull\(\)')

        foreach ($fieldMatch in $fieldMatches) {
            $fieldName = $fieldMatch.Groups[1].Value
            $totalFields++

            # Check if there's a corresponding else if data.FieldName.IsNull handler
            $hasElseHandler = $funcBody -match "else\s+if\s+data\.$fieldName\.IsNull\(\)"

            if (-not $hasElseHandler) {
                $fieldsWithIssues++

                # Don't report duplicate fields
                if ($resourceIssues.fieldName -notcontains $fieldName) {
                    $resourceIssues += [PSCustomObject]@{
                        Field = $fieldName
                        Function = $funcName
                    }
                }
            }
        }
    }

    if ($resourceIssues.Count -gt 0) {
        $resourcesWithIssues++

        # Determine priority
        $priority = 'LOW'
        if ($resourceName -match '^(ip_address|prefix|device|virtual_machine|vlan|site|asn|vrf|aggregate)$') {
            $priority = 'HIGH'
        } elseif ($resourceName -match '^(interface|rack|tenant|circuit|cluster|vm_interface|device_type|device_role)$') {
            $priority = 'MEDIUM'
        }

        $results += [PSCustomObject]@{
            Priority = $priority
            Resource = $resourceName
            File = $file.Name
            FieldCount = $resourceIssues.Count
            Fields = ($resourceIssues.Field | Sort-Object -Unique) -join ', '
            FieldsList = $resourceIssues.Field | Sort-Object -Unique
        }
    }
}

# Display Summary
Write-Host "╔════════════════════════════════════════════════════════════════╗" -ForegroundColor Cyan
Write-Host "║                        SUMMARY                                 ║" -ForegroundColor Cyan
Write-Host "╚════════════════════════════════════════════════════════════════╝" -ForegroundColor Cyan
Write-Host "Total Resources Scanned:        $totalResources" -ForegroundColor White
Write-Host "Resources with Issues:          $resourcesWithIssues" -ForegroundColor Yellow
Write-Host "Total Optional Fields Checked:  $totalFields" -ForegroundColor White
Write-Host "Fields Missing Null Handling:   $fieldsWithIssues`n" -ForegroundColor Red

# Display by Priority
$highPriority = $results | Where-Object Priority -eq 'HIGH' | Sort-Object Resource
$mediumPriority = $results | Where-Object Priority -eq 'MEDIUM' | Sort-Object Resource
$lowPriority = $results | Where-Object Priority -eq 'LOW' | Sort-Object Resource

if ($highPriority.Count -gt 0) {
    Write-Host "╔════════════════════════════════════════════════════════════════╗" -ForegroundColor Red
    Write-Host "║            HIGH PRIORITY ($($highPriority.Count) resources)                            ║" -ForegroundColor Red
    Write-Host "╚════════════════════════════════════════════════════════════════╝" -ForegroundColor Red
    foreach ($item in $highPriority) {
        Write-Host "`n  📌 $($item.Resource)" -ForegroundColor Yellow
        Write-Host "     Fields ($($item.FieldCount)): $($item.Fields)" -ForegroundColor Gray
    }
    Write-Host ""
}

if ($mediumPriority.Count -gt 0) {
    Write-Host "╔════════════════════════════════════════════════════════════════╗" -ForegroundColor Yellow
    Write-Host "║           MEDIUM PRIORITY ($($mediumPriority.Count) resources)                        ║" -ForegroundColor Yellow
    Write-Host "╚════════════════════════════════════════════════════════════════╝" -ForegroundColor Yellow
    foreach ($item in $mediumPriority) {
        Write-Host "`n  📋 $($item.Resource)" -ForegroundColor White
        Write-Host "     Fields ($($item.FieldCount)): $($item.Fields)" -ForegroundColor Gray
    }
    Write-Host ""
}

if ($lowPriority.Count -gt 0) {
    Write-Host "╔════════════════════════════════════════════════════════════════╗" -ForegroundColor Gray
    Write-Host "║            LOW PRIORITY ($($lowPriority.Count) resources)                            ║" -ForegroundColor Gray
    Write-Host "╚════════════════════════════════════════════════════════════════╝" -ForegroundColor Gray

    if ($Detailed) {
        foreach ($item in $lowPriority) {
            Write-Host "`n  📄 $($item.Resource)" -ForegroundColor DarkGray
            Write-Host "     Fields ($($item.FieldCount)): $($item.Fields)" -ForegroundColor DarkGray
        }
    } else {
        Write-Host "  (Use -Detailed flag to see all low priority resources)" -ForegroundColor DarkGray
        Write-Host "  Resources: $($lowPriority.Resource -join ', ')" -ForegroundColor DarkGray
    }
    Write-Host ""
}

# Export Results
$csvPath = ".\BUGFIX_ANALYSIS_optional_fields_detailed.csv"
$results | Export-Csv -Path $csvPath -NoTypeInformation
Write-Host "✓ Detailed results exported to: $csvPath" -ForegroundColor Green

if ($ExportJson) {
    $jsonPath = ".\BUGFIX_ANALYSIS_optional_fields_detailed.json"
    $results | ConvertTo-Json -Depth 10 | Out-File $jsonPath
    Write-Host "✓ JSON export saved to: $jsonPath" -ForegroundColor Green
}

# Generate Work Batches
Write-Host "`n╔════════════════════════════════════════════════════════════════╗" -ForegroundColor Cyan
Write-Host "║                  SUGGESTED WORK BATCHES                        ║" -ForegroundColor Cyan
Write-Host "╚════════════════════════════════════════════════════════════════╝`n" -ForegroundColor Cyan

$batchSize = 8
$allPrioritized = @($highPriority) + @($mediumPriority) + @($lowPriority)

$batchNum = 1
for ($i = 0; $i -lt $allPrioritized.Count; $i += $batchSize) {
    $batch = $allPrioritized[$i..[Math]::Min($i + $batchSize - 1, $allPrioritized.Count - 1)]
    $batchPriority = $batch[0].Priority

    Write-Host "Batch $batchNum ($batchPriority Priority - $($batch.Count) resources):" -ForegroundColor $(if($batchPriority -eq 'HIGH'){'Red'}elseif($batchPriority -eq 'MEDIUM'){'Yellow'}else{'Gray'})
    foreach ($item in $batch) {
        $fieldSummary = if ($item.FieldCount -le 3) { $item.Fields } else { "$($item.FieldCount) fields" }
        Write-Host "  • $($item.Resource) - $fieldSummary" -ForegroundColor White
    }
    Write-Host ""
    $batchNum++
}

# Summary Statistics
Write-Host "╔════════════════════════════════════════════════════════════════╗" -ForegroundColor Cyan
Write-Host "║                   COMPLETION ESTIMATES                         ║" -ForegroundColor Cyan
Write-Host "╚════════════════════════════════════════════════════════════════╝" -ForegroundColor Cyan
Write-Host "Estimated time per resource: ~30 minutes (code + tests)" -ForegroundColor Gray
Write-Host "High Priority:   $($highPriority.Count) resources × 30min = $([Math]::Ceiling($highPriority.Count * 0.5)) hours" -ForegroundColor Red
Write-Host "Medium Priority: $($mediumPriority.Count) resources × 30min = $([Math]::Ceiling($mediumPriority.Count * 0.5)) hours" -ForegroundColor Yellow
Write-Host "Low Priority:    $($lowPriority.Count) resources × 30min = $([Math]::Ceiling($lowPriority.Count * 0.5)) hours" -ForegroundColor Gray
Write-Host "─────────────────────────────────────────────────────────────────" -ForegroundColor Gray
$totalHours = [Math]::Ceiling($resourcesWithIssues * 0.5)
$totalDays = [Math]::Ceiling($totalHours / 8)
Write-Host "TOTAL ESTIMATED: $totalHours hours (~$totalDays working days)`n" -ForegroundColor White

Write-Host "╔════════════════════════════════════════════════════════════════╗" -ForegroundColor Green
Write-Host "║              BATCH 1A ANALYSIS COMPLETE ✓                      ║" -ForegroundColor Green
Write-Host "╚════════════════════════════════════════════════════════════════╝" -ForegroundColor Green
Write-Host "Next: Batch 1B - Create test infrastructure`n" -ForegroundColor Cyan
