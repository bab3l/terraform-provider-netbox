# Provider Factories Cleanup Plan

## Overview
Replace inline ProtoV6ProviderFactories declarations with the centralized `testutil.TestAccProtoV6ProviderFactories` variable.

**Total: 259 occurrences across 56 test files**

## Pattern Replacement

### From (inline, repeated everywhere):
```go
ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
    "netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
},
```

### To (centralized, DRY):
```go
ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
```

## Benefits
- **DRY Principle**: Single source of truth in `testutil.provider_factories.go`
- **Maintainability**: Changes to provider initialization only need to be made once
- **Readability**: Tests are more concise and focused on test logic
- **Consistency**: All tests use the exact same provider factory setup
- **Less Import Clutter**: No need to import `providerserver` and `provider` in every test file

## Files to Update (56 files, 259 occurrences)

| File | Occurrences |
|------|-------------|
| aggregate_resource_test.go | 4 |
| asn_range_resource_test.go | 5 |
| asn_resource_test.go | 3 |
| circuit_group_assignment_resource_test.go | 8 |
| circuit_group_resource_test.go | 7 |
| circuit_resource_test.go | 6 |
| circuit_termination_resource_test.go | 6 |
| circuit_type_resource_test.go | 7 |
| cluster_group_resource_test.go | 3 |
| cluster_resource_test.go | 6 |
| cluster_type_resource_test.go | 6 |
| config_context_resource_test.go | 3 |
| config_template_resource_test.go | 4 |
| console_port_resource_test.go | 3 |
| console_port_template_resource_test.go | 2 |
| console_server_port_resource_test.go | 3 |
| console_server_port_template_resource_test.go | 2 |
| fhrp_group_resource_test.go | 1 |
| l2vpn_termination_resource_test.go | 5 |
| location_resource_test.go | 7 |
| module_bay_resource_test.go | 5 |
| module_bay_template_resource_test.go | 4 |
| module_resource_test.go | 6 |
| module_type_resource_test.go | 6 |
| power_feed_resource_test.go | 3 |
| power_port_resource_test.go | 3 |
| prefix_resource_test.go | 6 |
| provider_account_resource_test.go | 7 |
| provider_network_resource_test.go | 7 |
| provider_resource_test.go | 7 |
| rack_reservation_resource_test.go | 6 |
| rack_resource_test.go | 1 |
| rack_type_resource_test.go | 1 |
| rear_port_resource_test.go | 3 |
| rear_port_template_resource_test.go | 2 |
| region_resource_test.go | 8 |
| rir_resource_test.go | 2 |
| role_resource_test.go | 4 |
| route_target_resource_test.go | 5 |
| service_resource_test.go | 5 |
| service_template_resource_test.go | 1 |
| site_group_resource_test.go | 7 |
| site_resource_test.go | 6 |
| tag_resource_test.go | 5 |
| tenant_group_resource_test.go | 7 |
| tenant_resource_test.go | 7 |
| tunnel_group_resource_test.go | 6 |
| tunnel_termination_resource_test.go | 6 |
| virtual_chassis_resource_test.go | 3 |
| virtual_device_context_resource_test.go | 4 |
| virtual_disk_resource_test.go | 5 |
| vlan_group_resource_test.go | 6 |
| vrf_resource_test.go | 5 |
| webhook_resource_test.go | 3 |
| wireless_lan_group_resource_test.go | 3 |
| wireless_lan_resource_test.go | 3 |

## Batch Approach

### Strategy
Use PowerShell to perform regex replacement across all files at once:

```powershell
$files = Get-ChildItem "internal/resources_acceptance_tests/*_test.go"
foreach ($file in $files) {
    $content = Get-Content $file.FullName -Raw
    $pattern = 'ProtoV6ProviderFactories:\s*map\[string\]func\(\)\s*\(tfprotov6\.ProviderServer,\s*error\)\s*\{\s*"netbox":\s*providerserver\.NewProtocol6WithError\(provider\.New\("test"\)\(\)\),\s*\}'
    $replacement = 'ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories'
    $newContent = $content -replace $pattern, $replacement
    if ($content -ne $newContent) {
        Set-Content $file.FullName -Value $newContent -NoNewline
        Write-Host "Updated: $($file.Name)"
    }
}
```

### Verification Steps
1. Run PowerShell replacement script
2. Build: `go build .`
3. Grep verify: `Select-String -Path "internal/resources_acceptance_tests/*.go" -Pattern 'ProtoV6ProviderFactories:\s*map\[string\]'`
4. Run a sample acceptance test to ensure provider initialization still works
5. Commit changes

### Potential Issues & Solutions
- **Import cleanup**: May need to remove unused imports (`providerserver`, `provider`, `tfprotov6`) after replacement
- **Multiline formatting**: Pattern may need adjustment for different whitespace/formatting variations
- **Already updated files**: Some files already use `testutil.TestAccProtoV6ProviderFactories` - skip these

## Post-Cleanup
After all inline patterns are replaced:
- All 259 + existing 409 = **668 total test cases** will use the centralized pattern
- Future tests can simply copy the pattern from any existing test
- Provider initialization changes only need to be made in one place
