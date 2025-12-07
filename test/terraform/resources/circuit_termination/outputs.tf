# Outputs to verify resource creation

output "basic_id" {
  description = "ID of the basic circuit termination"
  value       = netbox_circuit_termination.basic.id
}

output "basic_term_side" {
  description = "Term side of the basic termination"
  value       = netbox_circuit_termination.basic.term_side
}

output "basic_id_valid" {
  description = "Basic termination has valid ID"
  value       = netbox_circuit_termination.basic.id != ""
}

output "complete_id" {
  description = "ID of the complete circuit termination"
  value       = netbox_circuit_termination.complete.id
}

output "complete_term_side" {
  description = "Term side of the complete termination"
  value       = netbox_circuit_termination.complete.term_side
}

output "complete_port_speed" {
  description = "Port speed of the complete termination"
  value       = netbox_circuit_termination.complete.port_speed
}

output "complete_description" {
  description = "Description of the complete termination"
  value       = netbox_circuit_termination.complete.description
}

output "circuit_id" {
  description = "ID of the parent circuit"
  value       = netbox_circuit.test.id
}
