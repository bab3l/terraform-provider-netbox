# Testutil HTTP Status Code Cleanup Plan

## Overview
Replace HTTP status code literals (200, 404) with http package constants (http.StatusOK, http.StatusNotFound) across testutil files.

**Total occurrences: 152 matches across 4 files**

---

## Batch 1: cleanup_dcim.go (27 occurrences)
**File**: `internal/testutil/cleanup_dcim.go`

**Patterns to replace**:
- `resp.StatusCode != 200` → `resp.StatusCode != http.StatusOK` (26 occurrences)
- `resp.StatusCode == 404` → `resp.StatusCode == http.StatusNotFound` (1 occurrence at line 568)

**Lines with 200**:
- Lines 36, 86, 134, 182, 230, 278, 326, 374, 422, 470, 518, 558, 579, 621, 663, 703, 743, 791, 951, 999, 1047, 1095, 1169

**Lines with 404**:
- Line 568

**Strategy**: Use PowerShell regex replacement
```powershell
$content = Get-Content "internal/testutil/cleanup_dcim.go" -Raw
$content = $content -replace 'resp\.StatusCode != 200', 'resp.StatusCode != http.StatusOK'
$content = $content -replace 'resp\.StatusCode == 404', 'resp.StatusCode == http.StatusNotFound'
Set-Content "internal/testutil/cleanup_dcim.go" -Value $content -NoNewline
```

---

## Batch 2: cleanup_ipam_and_circuits.go (53 occurrences)
**File**: `internal/testutil/cleanup_ipam_and_circuits.go`

**Patterns to replace**:
- `resp.StatusCode != 200` → `resp.StatusCode != http.StatusOK` (51 occurrences)
- `resp.StatusCode == 404` → `resp.StatusCode == http.StatusNotFound` (1 occurrence at line 1878)
- Various compound conditions with != 200

**Lines with 200**:
- Lines 38, 86, 134, 182, 230, 284, 332, 380, 428, 476, 524, 572, 620, 668, 716, 756, 822, 870, 918, 972, 1020, 1068, 1116, 1158, 1228, 1276, 1324, 1372, 1420, 1468, 1516, 1566, 1630, 1678, 1732, 1780, 1828, 1868, 1889, 1937, 1985, 2067, 2115, 2253, 2301, 2349, 2397, 2445, 2479, 2527, 2575

**Lines with 404**:
- Line 1878

**Strategy**: Use PowerShell regex replacement
```powershell
$content = Get-Content "internal/testutil/cleanup_ipam_and_circuits.go" -Raw
$content = $content -replace 'resp\.StatusCode != 200', 'resp.StatusCode != http.StatusOK'
$content = $content -replace 'vmResp\.StatusCode != 200', 'vmResp.StatusCode != http.StatusOK'
$content = $content -replace 'deviceResp\.StatusCode != 200', 'deviceResp.StatusCode != http.StatusOK'
$content = $content -replace 'resp\.StatusCode == 404', 'resp.StatusCode == http.StatusNotFound'
Set-Content "internal/testutil/cleanup_ipam_and_circuits.go" -Value $content -NoNewline
```

---

## Batch 3: check_destroy_dcim.go (31 occurrences)
**File**: `internal/testutil/check_destroy_dcim.go`

**Patterns to replace**:
- `resp.StatusCode == 200` → `resp.StatusCode == http.StatusOK` (31 occurrences)

**Lines with 200**:
- Lines 60, 112, 164, 216, 268, 320, 372, 424, 476, 528, 578, 598, 648, 706, 760, 814, 868, 920, 974, 1028, 1082, 1136, 1188, 1240, 1292, 1344, 1395, 1428

**Strategy**: Use PowerShell regex replacement
```powershell
$content = Get-Content "internal/testutil/check_destroy_dcim.go" -Raw
$content = $content -replace 'resp\.StatusCode == 200', 'resp.StatusCode == http.StatusOK'
Set-Content "internal/testutil/check_destroy_dcim.go" -Value $content -NoNewline
```

---

## Batch 4: check_destroy_ipam_and_circuits.go (41 occurrences)
**File**: `internal/testutil/check_destroy_ipam_and_circuits.go`

**Patterns to replace**:
- `resp.StatusCode == 200` → `resp.StatusCode == http.StatusOK` (41 occurrences)

**Lines with 200**:
- Lines 54, 106, 158, 210, 262, 314, 366, 412, 462, 512, 568, 620, 672, 726, 778, 830, 882, 934, 986, 1038, 1090, 1144, 1196, 1248, 1300, 1352, 1404, 1456, 1508, 1562, 1614, 1666, 1718, 1770, 1816, 1872, 1924, 1978, 2030, 2082, 2136, 2190, 2244, 2296, 2348, 2400, 2452, 2504

**Strategy**: Use PowerShell regex replacement
```powershell
$content = Get-Content "internal/testutil/check_destroy_ipam_and_circuits.go" -Raw
$content = $content -replace 'resp\.StatusCode == 200', 'resp.StatusCode == http.StatusOK'
Set-Content "internal/testutil/check_destroy_ipam_and_circuits.go" -Value $content -NoNewline
```

---

## Execution Order

1. **Batch 1**: cleanup_dcim.go (simplest - 2 patterns)
2. **Batch 2**: cleanup_ipam_and_circuits.go (multiple patterns including special variables)
3. **Batch 3**: check_destroy_dcim.go (single pattern, many occurrences)
4. **Batch 4**: check_destroy_ipam_and_circuits.go (single pattern, many occurrences)

---

## Verification Steps

After each batch:
1. Run: `go build .` to ensure no syntax errors
2. Run: `grep -n "StatusCode.*200\|StatusCode.*404" internal/testutil/<filename>` to verify replacements
3. Run unit tests if applicable

After all batches:
1. Run: `grep -rn "StatusCode.*\b(200|404)\b" internal/testutil/` to confirm no literals remain
2. Run full acceptance test suite
3. Commit changes with descriptive message

---

## Notes
- All testutil files appear to follow consistent patterns
- cleanup_* files check for successful responses (200) to find resources to clean up
- check_destroy_* files verify resources are destroyed (should NOT find 200)
- Only 2 occurrences of 404 checks (both in cleanup files for fallback scenarios)
