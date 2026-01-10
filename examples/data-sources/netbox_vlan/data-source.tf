# Look up a VLAN by ID
data "netbox_vlan" "by_id" {
  id = "1"
}

# Look up a VLAN by VID
data "netbox_vlan" "by_vid" {
  vid = "100"
}

# Look up a VLAN by name
data "netbox_vlan" "by_name" {
  name = "test-vlan"
}

# Look up a VLAN by VID with optional name filter
data "netbox_vlan" "by_vid_and_name" {
  vid  = "100"
  name = "test-vlan"
}

# Use VLAN data in other resources
output "vlan_name" {
  value = data.netbox_vlan.by_id.name
}

output "vlan_vid" {
  value = data.netbox_vlan.by_name.vid
}

output "vlan_site" {
  value = data.netbox_vlan.by_id.site
}

output "vlan_group" {
  value = data.netbox_vlan.by_id.group
}

output "vlan_status" {
  value = data.netbox_vlan.by_id.status
}

output "vlan_role" {
  value = data.netbox_vlan.by_vid.role
}

output "vlan_tenant" {
  value = data.netbox_vlan.by_id.tenant
}

output "vlan_description" {
  value = data.netbox_vlan.by_name.description
}

# Access all custom fields
output "vlan_custom_fields" {
  value       = data.netbox_vlan.by_id.custom_fields
  description = "All custom fields defined in NetBox for this VLAN"
}

# Access specific custom fields by name
output "vlan_subnet" {
  value       = try([for cf in data.netbox_vlan.by_id.custom_fields : cf.value if cf.name == "subnet"][0], null)
  description = "Example: accessing a text custom field"
}

output "vlan_dhcp_enabled" {
  value       = try([for cf in data.netbox_vlan.by_id.custom_fields : cf.value if cf.name == "dhcp_enabled"][0], null)
  description = "Example: accessing a boolean custom field"
}
