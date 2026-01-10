# Example: Look up an L2VPN termination by ID (only supported lookup method)
data "netbox_l2vpn_termination" "by_id" {
  id = "123"
}

# Example: Use L2VPN termination data in other resources
output "termination_id" {
  value = data.netbox_l2vpn_termination.by_id.id
}

output "termination_l2vpn" {
  value = data.netbox_l2vpn_termination.by_id.l2vpn
}

output "termination_object_type" {
  value = data.netbox_l2vpn_termination.by_id.assigned_object_type
}

output "termination_object_id" {
  value = data.netbox_l2vpn_termination.by_id.assigned_object_id
}

# Access all custom fields
output "termination_custom_fields" {
  value       = data.netbox_l2vpn_termination.by_id.custom_fields
  description = "All custom fields defined in NetBox for this L2VPN termination"
}

# Access specific custom field by name
output "termination_interface_speed" {
  value       = try([for cf in data.netbox_l2vpn_termination.by_id.custom_fields : cf.value if cf.name == "interface_speed_mbps"][0], null)
  description = "Example: accessing a numeric custom field for interface speed"
}

output "termination_is_primary" {
  value       = try([for cf in data.netbox_l2vpn_termination.by_id.custom_fields : cf.value if cf.name == "is_primary"][0], null)
  description = "Example: accessing a boolean custom field for primary status"
}
