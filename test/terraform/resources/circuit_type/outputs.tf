output "circuit_type_basic_id" {
  description = "ID of the basic circuit type"
  value       = netbox_circuit_type.basic.id
}

output "circuit_type_basic_name" {
  description = "Name of the basic circuit type"
  value       = netbox_circuit_type.basic.name
}

output "circuit_type_basic_slug" {
  description = "Slug of the basic circuit type"
  value       = netbox_circuit_type.basic.slug
}

output "circuit_type_complete_id" {
  description = "ID of the complete circuit type"
  value       = netbox_circuit_type.complete.id
}

output "circuit_type_complete_color" {
  description = "Color of the complete circuit type"
  value       = netbox_circuit_type.complete.color
}

output "basic_circuit_type_valid" {
  description = "Validates basic circuit type was created correctly"
  value       = netbox_circuit_type.basic.id != "" && netbox_circuit_type.basic.slug == "basic-test-circuit-type"
}

output "complete_circuit_type_valid" {
  description = "Validates complete circuit type was created correctly"
  value       = netbox_circuit_type.complete.id != "" && netbox_circuit_type.complete.color == "ff5722"
}

output "internet_circuit_type_valid" {
  description = "Validates internet circuit type was created correctly"
  value       = netbox_circuit_type.internet.id != "" && netbox_circuit_type.internet.slug == "internet"
}
