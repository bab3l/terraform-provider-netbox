# Optional Field Null Handling Analysis Script
# This script analyzes all resources to identify optional fields that may not handle null values correctly

Write-Host "`n=== OPTIONAL FIELD NULL HANDLING ANALYSIS ===" -ForegroundColor Cyan
Write-Host "Date: $(Get-Date -Format 'yyyy-MM-dd HH:mm:ss')`n" -ForegroundColor Gray

$resourcePath = ".\internal\resources"
$resources = Get-ChildItem -Path "$resourcePath\*_resource.go" | Where-Object { $_.Name -ne "provider_resource.go" }

Write-Host "Analyzing $($resources.Count) resources...`n" -ForegroundColor Yellow

$analysisResults = @()
$totalIssues = 0
$totalFieldsChecked = 0

foreach ($file in $resources) {
    $resourceName = $file.BaseName -replace '_resource$', ''
    $content = Get-Content $file.FullName -Raw

    # Find all optional fields from the model struct
    # Match types.String, types.Int64, types.Bool that have 'Optional: true' in schema or are in build functions
    $optionalFieldsFound = @{}

    # Look for fields that appear in IsNull() checks in build functions
    $buildFunctions = @()
    if ($content -match '(?s)(func \(r \*\w+\) build\w+(?:Request|CreateRequest|UpdateRequest)\([^)]+\)[^{]*\{.*?)(?=\nfunc \(r \*\w+\) (?:map|build|Create|Read|Update|Delete|Import))') {
        $buildFunctions += $Matches[1]
    }

    # Also try to capture the last build function before mapResponseToModel
    if ($content -match '(?s)(func \(r \*\w+\) build\w+(?:Request|CreateRequest|UpdateRequest)\([^)]+\)[^{]*\{.*?)(?=\n+// mapResponseToModel|\nfunc \(r \*\w+\) mapResponseToModel)') {
        $buildFunctions += $Matches[1]
    }

    if ($buildFunctions.Count -eq 0) {
        continue
    }

    $buildFunctionContent = $buildFunctions -join "`n`n"

    # Find all fields that have IsNull checks
    $isNullMatches = [regex]::Matches($buildFunctionContent, 'if !data\.(\w+)\.IsNull\(\)')

    foreach ($match in $isNullMatches) {
        $fieldName = $match.Groups[1].Value
        $totalFieldsChecked++

        # Check if this field has corresponding null handling (else if ... IsNull)
        $hasNullHandling = $false

        # Look for "else if data.FieldName.IsNull()" pattern
        if ($buildFunctionContent -match "else if data\.$fieldName\.IsNull\(\)") {
            $hasNullHandling = $true
        }

        # Track this field
        if (-not $optionalFieldsFound.ContainsKey($fieldName)) {
            $optionalFieldsFound[$fieldName] = @{
                HasSetCheck = $true
                HasNullHandling = $hasNullHandling
            }
        }
    }

    # Find fields with issues
    $fieldsWithIssues = @()
    foreach ($fieldName in $optionalFieldsFound.Keys) {
        if ($optionalFieldsFound[$fieldName].HasSetCheck -and -not $optionalFieldsFound[$fieldName].HasNullHandling) {
            $fieldsWithIssues += $fieldName
            $totalIssues++
        }
    }

    if ($fieldsWithIssues.Count -gt 0) {
        # Determine priority based on resource type
        $priority = 'LOW'
        if ($resourceName -in @('ip_address', 'prefix', 'device', 'virtual_machine', 'vlan', 'site', 'asn', 'vrf')) {
            $priority = 'HIGH'
        } elseif ($resourceName -in @('interface', 'rack', 'tenant', 'circuit', 'cluster', 'vm_interface', 'device_type')) {
            $priority = 'MEDIUM'
        }

        $analysisResults += [PSCustomObject]@{
            Resource = $resourceName
            File = $file.Name
            OptionalFieldsCount = $optionalFieldsFound.Keys.Count
            FieldsWithIssues = ($fieldsWithIssues | Sort-Object) -join ', '
            IssueCount = $fieldsWithIssues.Count
            Priority = $priority
        }
    }
}

Write-Host "=== SUMMARY ===" -ForegroundColor Cyan
Write-Host "Resources analyzed: $($resources.Count)" -ForegroundColor White
Write-Host "Resources with issues: $($analysisResults.Count)" -ForegroundColor Yellow
Write-Host "Total field issues found: $totalIssues`n" -ForegroundColor Red

Write-Host "=== DETAILED RESULTS ===" -ForegroundColor Cyan
Write-Host "(Showing resources that need fixes)`n" -ForegroundColor Gray

# Group by priority
$highPriority = $analysisResults | Where-Object { $_.Priority -eq 'HIGH' } | Sort-Object Resource
$mediumPriority = $analysisResults | Where-Object { $_.Priority -eq 'MEDIUM' } | Sort-Object Resource
$lowPriority = $analysisResults | Where-Object { $_.Priority -eq 'LOW' } | Sort-Object Resource

if ($highPriority.Count -gt 0) {
    Write-Host "`n--- HIGH PRIORITY ($($highPriority.Count) resources) ---" -ForegroundColor Red
    $highPriority | Format-Table Resource, FieldsWithIssues, IssueCount -AutoSize | Out-String -Width 200 | Write-Host
}

if ($mediumPriority.Count -gt 0) {
    Write-Host "`n--- MEDIUM PRIORITY ($($mediumPriority.Count) resources) ---" -ForegroundColor Yellow
    $mediumPriority | Format-Table Resource, FieldsWithIssues, IssueCount -AutoSize | Out-String -Width 200 | Write-Host
}

if ($lowPriority.Count -gt 0) {
    Write-Host "`n--- LOW PRIORITY ($($lowPriority.Count) resources) ---" -ForegroundColor Gray
    $lowPriority | Format-Table Resource, FieldsWithIssues, IssueCount -AutoSize | Out-String -Width 200 | Write-Host
}

# Export detailed results to CSV
$csvPath = ".\BUGFIX_ANALYSIS_optional_field_null_handling.csv"
$analysisResults | Export-Csv -Path $csvPath -NoTypeInformation
Write-Host "`nDetailed results exported to: $csvPath" -ForegroundColor Green

# Generate work batches
Write-Host "`n=== SUGGESTED WORK BATCHES ===" -ForegroundColor Cyan

$batchSize = 10
$batchNumber = 1

Write-Host "`nHIGH PRIORITY BATCHES:" -ForegroundColor Red
$highPriority | ForEach-Object -Begin { $count = 0; $batch = @() } -Process {
    $batch += $_
    $count++
    if ($count -eq $batchSize -or $_ -eq $highPriority[-1]) {
        Write-Host "`nBatch $batchNumber (HIGH):" -ForegroundColor Yellow
        $batch | ForEach-Object { Write-Host "  - $($_.Resource) ($($_.FieldsWithIssues))" }
        $batch = @()
        $batchNumber++
    }
}

Write-Host "`nMEDIUM PRIORITY BATCHES:" -ForegroundColor Yellow
$mediumPriority | ForEach-Object -Begin { $count = 0; $batch = @() } -Process {
    $batch += $_
    $count++
    if ($count -eq $batchSize -or $_ -eq $mediumPriority[-1]) {
        Write-Host "`nBatch $batchNumber (MEDIUM):" -ForegroundColor Yellow
        $batch | ForEach-Object { Write-Host "  - $($_.Resource) ($($_.FieldsWithIssues))" }
        $batch = @()
        $batchNumber++
    }
}

Write-Host "`n=== ANALYSIS COMPLETE ===" -ForegroundColor Green
Write-Host "Review the results and planning document (BUGFIX_PLAN_optional_field_null_handling.md)" -ForegroundColor White
Write-Host "Ready to begin Phase 1, Batch 1B: Test Infrastructure`n" -ForegroundColor Cyan
