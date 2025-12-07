# Provider Network Data Source Outputs

output "by_id_name" {
  value = data.netbox_provider_network.by_id.name
}

output "by_id_circuit_provider" {
  value = data.netbox_provider_network.by_id.circuit_provider
}

output "by_id_description" {
  value = data.netbox_provider_network.by_id.description
}
