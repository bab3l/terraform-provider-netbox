# Look up a VLAN group by ID
data "netbox_vlan_group" "by_id" {
  id = "1"
}

# Look up a VLAN group by slug
data "netbox_vlan_group" "by_slug" {
  slug = "test-vlan-group"
}

# Look up a VLAN group by name
data "netbox_vlan_group" "by_name" {
  name = "test-vlan-group"
}

# Use VLAN group data in other resources
output "vlan_group_name" {
  value = data.netbox_vlan_group.by_id.name
}

output "vlan_group_slug" {
  value = data.netbox_vlan_group.by_name.slug
}

output "vlan_group_description" {
  value = data.netbox_vlan_group.by_id.description
}

output "vlan_group_scope_type" {
  value = data.netbox_vlan_group.by_id.scope_type
}

output "vlan_group_scope_id" {
  value = data.netbox_vlan_group.by_id.scope_id
}

output "vlan_group_min_vid" {
  value = data.netbox_vlan_group.by_slug.min_vid
}

output "vlan_group_max_vid" {
  value = data.netbox_vlan_group.by_slug.max_vid
}

# Access all custom fields
output "vlan_group_custom_fields" {
  value       = data.netbox_vlan_group.by_id.custom_fields
  description = "All custom fields defined in NetBox for this VLAN group"
}

# Access specific custom fields by name
output "vlan_group_region" {
  value       = try([for cf in data.netbox_vlan_group.by_id.custom_fields : cf.value if cf.name == "region_name"][0], null)
  description = "Example: accessing a text custom field"
}

output "vlan_group_managed" {
  value       = try([for cf in data.netbox_vlan_group.by_id.custom_fields : cf.value if cf.name == "managed"][0], null)
  description = "Example: accessing a boolean custom field"
}
