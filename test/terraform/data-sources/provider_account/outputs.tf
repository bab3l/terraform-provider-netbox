# Provider Account Data Source Outputs

output "by_id_name" {
  value = data.netbox_provider_account.by_id.name
}

output "by_id_circuit_provider" {
  value = data.netbox_provider_account.by_id.circuit_provider
}

output "by_id_description" {
  value = data.netbox_provider_account.by_id.description
}
