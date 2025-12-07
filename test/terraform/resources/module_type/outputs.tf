# Outputs to verify resource creation

output "basic_id" {
  description = "ID of the basic module type"
  value       = netbox_module_type.basic.id
}

output "basic_model" {
  description = "Model of the basic module type"
  value       = netbox_module_type.basic.model
}

output "basic_id_valid" {
  description = "Basic module type has valid ID"
  value       = netbox_module_type.basic.id != ""
}

output "complete_id" {
  description = "ID of the complete module type"
  value       = netbox_module_type.complete.id
}

output "complete_model" {
  description = "Model of the complete module type"
  value       = netbox_module_type.complete.model
}

output "complete_description" {
  description = "Description of the complete module type"
  value       = netbox_module_type.complete.description
}

output "manufacturer_id" {
  description = "ID of the manufacturer"
  value       = netbox_manufacturer.test.id
}
