# Lookup by ID
data "netbox_prefix" "by_id" {
  id = "123"
}

# Lookup by prefix CIDR
data "netbox_prefix" "by_cidr" {
  prefix = "10.0.0.0/24"
}

# Use prefix data in other resources
output "prefix_cidr" {
  value = data.netbox_prefix.by_id.prefix
}

output "prefix_status" {
  value = data.netbox_prefix.by_cidr.status
}

output "prefix_vrf" {
  value = data.netbox_prefix.by_id.vrf
}

output "prefix_tenant" {
  value = data.netbox_prefix.by_id.tenant
}

output "prefix_site" {
  value = data.netbox_prefix.by_id.site
}

output "prefix_role" {
  value = data.netbox_prefix.by_id.role
}

output "prefix_is_pool" {
  value = data.netbox_prefix.by_cidr.is_pool
}

# Access all custom fields
output "prefix_custom_fields" {
  value       = data.netbox_prefix.by_id.custom_fields
  description = "All custom fields defined in NetBox for this prefix"
}

# Access specific custom fields by name
output "prefix_vlan_id" {
  value       = try([for cf in data.netbox_prefix.by_id.custom_fields : cf.value if cf.name == "vlan_id"][0], null)
  description = "Example: accessing a numeric custom field"
}

output "prefix_supernet" {
  value       = try([for cf in data.netbox_prefix.by_id.custom_fields : cf.value if cf.name == "supernet"][0], null)
  description = "Example: accessing a text custom field"
}
