# Look up IP range by ID
data "netbox_ip_range" "by_id" {
  id = "1"
}

# Look up IP range by start and end addresses
data "netbox_ip_range" "by_addresses" {
  start_address = "10.0.0.0"
  end_address   = "10.0.0.255"
}

# Use IP range data in other resources
output "range_start" {
  value = data.netbox_ip_range.by_id.start_address
}

output "range_end" {
  value = data.netbox_ip_range.by_addresses.end_address
}

output "range_size" {
  value = data.netbox_ip_range.by_id.size
}

output "range_vrf" {
  value = data.netbox_ip_range.by_id.vrf
}

output "range_tenant" {
  value = data.netbox_ip_range.by_id.tenant
}

output "range_status" {
  value = data.netbox_ip_range.by_id.status
}

output "range_description" {
  value = data.netbox_ip_range.by_addresses.description
}

# Access all custom fields
output "range_custom_fields" {
  value       = data.netbox_ip_range.by_id.custom_fields
  description = "All custom fields defined in NetBox for this IP range"
}

# Access specific custom fields by name
output "range_dhcp_enabled" {
  value       = try([for cf in data.netbox_ip_range.by_id.custom_fields : cf.value if cf.name == "dhcp_enabled"][0], null)
  description = "Example: accessing a boolean custom field"
}

output "range_pool_name" {
  value       = try([for cf in data.netbox_ip_range.by_id.custom_fields : cf.value if cf.name == "pool_name"][0], null)
  description = "Example: accessing a text custom field"
}
