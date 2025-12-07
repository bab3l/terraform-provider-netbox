# Provider Network Outputs

# Basic provider network outputs
output "basic_id" {
  value = netbox_provider_network.basic.id
}

output "basic_name" {
  value = netbox_provider_network.basic.name
}

# Complete provider network outputs
output "complete_id" {
  value = netbox_provider_network.complete.id
}

output "complete_name" {
  value = netbox_provider_network.complete.name
}

output "complete_service_id" {
  value = netbox_provider_network.complete.service_id
}
