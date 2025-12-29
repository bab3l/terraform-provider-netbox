#!/usr/bin/env pwsh

# Process files one at a time to avoid PowerShell array issues
$fileList = "circuit_resource.go|circuit_termination_resource.go|cluster_group_resource.go|config_template_resource.go|console_port_resource.go|console_port_template_resource.go|console_server_port_resource.go|console_server_port_template_resource.go|contact_resource.go|custom_field_choice_set_resource.go|custom_field_resource.go|custom_link_resource.go|device_bay_resource.go|device_bay_template_resource.go|device_resource.go|event_rule_resource.go|export_template_resource.go|fhrp_group_assignment_resource.go|fhrp_group_resource.go|front_port_resource.go|front_port_template_resource.go|ike_policy_resource.go|ike_proposal_resource.go|interface_template_resource.go|inventory_item_resource.go|inventory_item_role_resource.go|inventory_item_template_resource.go|ip_address_resource.go|ip_range_resource.go|journal_entry_resource.go|l2vpn_resource.go|l2vpn_termination_resource.go|module_bay_template_resource.go|module_resource.go|notification_group_resource.go|power_feed_resource.go|power_outlet_resource.go|power_outlet_template_resource.go|power_panel_resource.go|power_port_resource.go|power_port_template_resource.go|prefix_resource.go|provider_account_resource.go|provider_network_resource.go|rack_resource.go|rack_type_resource.go|rear_port_resource.go|rear_port_template_resource.go|route_target_resource.go|service_template_resource.go|tag_resource.go|tunnel_group_resource.go|tunnel_resource.go|virtual_chassis_resource.go|virtual_device_context_resource.go|virtual_disk_resource.go|virtual_machine_resource.go|vm_interface_resource.go|vrf_resource.go|wireless_lan_group_resource.go|wireless_link_resource.go"

$files = $fileList -split '\|'
$fixed = 0

foreach ($f in $files) {
    $path = "internal\resources\$f"
    $content = [System.IO.File]::ReadAllText($path)

    $match = [regex]::Match($content, '\tif (\w+)\.(?:Get)?Display\(\)? != ""')
    if ($match.Success) {
        $varName = $match.Groups[1].Value
        $found = $false

        # Try DisplayName first
        $pattern = "`n`t// DisplayName`n`tif $varName.Display != `"`" {`n`t} else {`n`t}`n"
        if ($content.Contains($pattern)) {
            $content = $content.Replace($pattern, "`n")
            $found = $true
        } else {
            # Try "Display name"
            $pattern = "`n`t// Display name`n`tif $varName.Display != `"`" {`n`t} else {`n`t}`n"
            if ($content.Contains($pattern)) {
                $content = $content.Replace($pattern, "`n")
                $found = $true
            } else {
                # Try GetDisplay()
                $pattern = "`n`t// DisplayName`n`tif $varName.GetDisplay() != `"`" {`n`t} else {`n`t}`n"
                if ($content.Contains($pattern)) {
                    $content = $content.Replace($pattern, "`n")
                    $found = $true
                }
            }
        }

        if ($found) {
            [System.IO.File]::WriteAllText($path, $content)
            $fixed++
            Write-Host "Fixed $f"
        }
    }
}

Write-Host ""
Write-Host "Total: $fixed files fixed"
