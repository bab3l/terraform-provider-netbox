output "provider_basic_id" {
  description = "ID of the basic provider"
  value       = netbox_provider.basic.id
}

output "provider_basic_name" {
  description = "Name of the basic provider"
  value       = netbox_provider.basic.name
}

output "provider_basic_slug" {
  description = "Slug of the basic provider"
  value       = netbox_provider.basic.slug
}

output "provider_complete_id" {
  description = "ID of the complete provider"
  value       = netbox_provider.complete.id
}

output "provider_complete_description" {
  description = "Description of the complete provider"
  value       = netbox_provider.complete.description
}

output "basic_provider_valid" {
  description = "Validates basic provider was created correctly"
  value       = netbox_provider.basic.id != "" && netbox_provider.basic.slug == "basic-test-provider"
}

output "complete_provider_valid" {
  description = "Validates complete provider was created correctly"
  value       = netbox_provider.complete.id != "" && netbox_provider.complete.description == "Complete provider for integration testing"
}
