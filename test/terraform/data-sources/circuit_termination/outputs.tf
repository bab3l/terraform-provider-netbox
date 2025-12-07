# Circuit Termination Data Source Outputs

output "by_id_circuit" {
  value = data.netbox_circuit_termination.by_id.circuit
}

output "by_id_term_side" {
  value = data.netbox_circuit_termination.by_id.term_side
}

output "by_id_site" {
  value = data.netbox_circuit_termination.by_id.site
}

output "by_id_description" {
  value = data.netbox_circuit_termination.by_id.description
}
