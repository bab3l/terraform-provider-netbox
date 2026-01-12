# Comprehensive Null Handling Analysis
$resources = Get-ChildItem -Path ".\internal\resources\*_resource.go" | Where-Object { $_.Name -notmatch 'provider_resource' }
$results = @()

Write-Host "Analyzing $($resources.Count) resources..." -ForegroundColor Yellow

foreach ($file in $resources) {
    $content = Get-Content $file.FullName -Raw
    $resourceName = $file.BaseName -replace '_resource$', ''

    # Find if !data.X.IsNull checks
    $checks = [regex]::Matches($content, 'if\s+!data\.(\w+)\.IsNull\(\)')

    $issues = @()
    foreach ($check in $checks) {
        $field = $check.Groups[1].Value
        # Check for else if data.X.IsNull
        if ($content -notmatch "else\s+if\s+data\.$field\.IsNull\(\)") {
            $issues += $field
        }
    }

    if ($issues.Count -gt 0) {
        $priority = 'LOW'
        if ($resourceName -match '^(ip_address|prefix|device|virtual_machine|vlan|site|asn|vrf|aggregate)$') { $priority = 'HIGH' }
        elseif ($resourceName -match '^(interface|rack|tenant|circuit|cluster|vm_interface|device_type|device_role)$') { $priority = 'MEDIUM' }

        $results += [PSCustomObject]@{
            Priority = $priority
            Resource = $resourceName
            IssueCount = $issues.Count
            Fields = ($issues | Sort-Object -Unique) -join ', '
        }
    }
}

Write-Host "`n=== NULL HANDLING ANALYSIS COMPLETE ===`n" -ForegroundColor Cyan
Write-Host "Resources with issues: $($results.Count)" -ForegroundColor Yellow
Write-Host "Total field issues: $(($results | Measure-Object -Property IssueCount -Sum).Sum)`n" -ForegroundColor Red

$high = $results | Where-Object Priority -eq 'HIGH' | Sort-Object Resource
$med = $results | Where-Object Priority -eq 'MEDIUM' | Sort-Object Resource
$low = $results | Where-Object Priority -eq 'LOW' | Sort-Object Resource

if ($high) {
    Write-Host "HIGH PRIORITY ($($high.Count) resources):" -ForegroundColor Red
    $high | Format-Table Resource, IssueCount, Fields -AutoSize | Out-String -Width 250 | Write-Host
}

if ($med) {
    Write-Host "MEDIUM PRIORITY ($($med.Count) resources):" -ForegroundColor Yellow
    $med | Format-Table Resource, IssueCount, Fields -AutoSize | Out-String -Width 250 | Write-Host
}

if ($low) {
    Write-Host "LOW PRIORITY ($($low.Count) resources):" -ForegroundColor Gray
    $low | Format-Table Resource, IssueCount, Fields -AutoSize | Out-String -Width 250 | Write-Host
}

$results | Export-Csv -Path ".\BUGFIX_ANALYSIS_null_handling.csv" -NoTypeInformation
Write-Host "Results exported to: BUGFIX_ANALYSIS_null_handling.csv" -ForegroundColor Green
