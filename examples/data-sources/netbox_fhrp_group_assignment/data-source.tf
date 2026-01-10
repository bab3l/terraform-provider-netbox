# Example: Look up an FHRP group assignment by ID (only supported lookup method)
data "netbox_fhrp_group_assignment" "by_id" {
  id = "1"
}

# Example: Use FHRP group assignment data in other resources
output "assignment_id" {
  value = data.netbox_fhrp_group_assignment.by_id.id
}

output "assignment_group_id" {
  value = data.netbox_fhrp_group_assignment.by_id.group_id
}

output "assignment_interface_type" {
  value = data.netbox_fhrp_group_assignment.by_id.interface_type
}

output "assignment_interface_id" {
  value = data.netbox_fhrp_group_assignment.by_id.interface_id
}

output "assignment_priority" {
  value = data.netbox_fhrp_group_assignment.by_id.priority
}

# Access all custom fields
output "assignment_custom_fields" {
  value       = data.netbox_fhrp_group_assignment.by_id.custom_fields
  description = "All custom fields defined in NetBox for this FHRP group assignment"
}

# Access specific custom field by name
output "assignment_tracking_interface" {
  value       = try([for cf in data.netbox_fhrp_group_assignment.by_id.custom_fields : cf.value if cf.name == "tracking_interface"][0], null)
  description = "Example: accessing a text custom field for tracking interface"
}

output "assignment_weight" {
  value       = try([for cf in data.netbox_fhrp_group_assignment.by_id.custom_fields : cf.value if cf.name == "weight"][0], null)
  description = "Example: accessing a numeric custom field for weight"
}
