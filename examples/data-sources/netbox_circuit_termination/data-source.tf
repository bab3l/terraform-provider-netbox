data "netbox_circuit_termination" "test" {
  id = "123"
}

output "termination_id" {
  value = data.netbox_circuit_termination.test.id
}

output "termination_circuit" {
  value = data.netbox_circuit_termination.test.circuit
}

output "termination_term_side" {
  value = data.netbox_circuit_termination.test.term_side
}

output "termination_site" {
  value = data.netbox_circuit_termination.test.site
}

output "termination_provider_network" {
  value = data.netbox_circuit_termination.test.provider_network
}

# Access all custom fields
output "termination_custom_fields" {
  value       = data.netbox_circuit_termination.test.custom_fields
  description = "All custom fields defined in NetBox for this circuit termination"
}

# Access specific custom field by name
output "termination_port_speed" {
  value       = try([for cf in data.netbox_circuit_termination.test.custom_fields : cf.value if cf.name == "port_speed_gbps"][0], null)
  description = "Example: accessing a numeric custom field for port speed"
}

output "termination_is_primary" {
  value       = try([for cf in data.netbox_circuit_termination.test.custom_fields : cf.value if cf.name == "is_primary"][0], null)
  description = "Example: accessing a boolean custom field for primary status"
}
