data "netbox_circuit_group_assignment" "example" {
  id = "456"
}

output "assignment_circuit" {
  value = data.netbox_circuit_group_assignment.example.circuit_cid
}

output "assignment_group" {
  value = data.netbox_circuit_group_assignment.example.group
}

output "assignment_priority" {
  value = data.netbox_circuit_group_assignment.example.priority
}

# Access all custom fields
output "assignment_custom_fields" {
  value       = data.netbox_circuit_group_assignment.example.custom_fields
  description = "All custom fields defined in NetBox for this circuit group assignment"
}

# Access specific custom field by name
output "assignment_failover_mode" {
  value       = try([for cf in data.netbox_circuit_group_assignment.example.custom_fields : cf.value if cf.name == "failover_mode"][0], null)
  description = "Example: accessing a select custom field for failover mode"
}

output "assignment_is_active" {
  value       = try([for cf in data.netbox_circuit_group_assignment.example.custom_fields : cf.value if cf.name == "is_active"][0], null)
  description = "Example: accessing a boolean custom field for active status"
}
