# Look up a wireless link by ID
data "netbox_wireless_link" "by_id" {
  id = "1"
}

# Use wireless link data in outputs
output "link_id" {
  value = data.netbox_wireless_link.by_id.id
}

output "link_interface_a" {
  value = data.netbox_wireless_link.by_id.interface_a
}

output "link_interface_b" {
  value = data.netbox_wireless_link.by_id.interface_b
}

output "link_ssid" {
  value = data.netbox_wireless_link.by_id.ssid
}

output "link_status" {
  value = data.netbox_wireless_link.by_id.status
}

output "link_auth_type" {
  value = data.netbox_wireless_link.by_id.auth_type
}

output "link_auth_cipher" {
  value = data.netbox_wireless_link.by_id.auth_cipher
}

output "link_distance" {
  value = data.netbox_wireless_link.by_id.distance
}

output "link_distance_unit" {
  value = data.netbox_wireless_link.by_id.distance_unit
}

# Access all custom fields
output "link_custom_fields" {
  value       = data.netbox_wireless_link.by_id.custom_fields
  description = "All custom fields defined in NetBox for this wireless link"
}

# Access specific custom field by name
output "link_frequency" {
  value       = try([for cf in data.netbox_wireless_link.by_id.custom_fields : cf.value if cf.name == "frequency_ghz"][0], null)
  description = "Example: accessing a numeric custom field for frequency"
}

output "link_is_active" {
  value       = try([for cf in data.netbox_wireless_link.by_id.custom_fields : cf.value if cf.name == "is_active"][0], null)
  description = "Example: accessing a boolean custom field for active status"
}
