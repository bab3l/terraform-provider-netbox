# Outputs to verify resource creation

output "basic_id" {
  description = "ID of the basic config template"
  value       = netbox_config_template.basic.id
}

output "basic_name" {
  description = "Name of the basic config template"
  value       = netbox_config_template.basic.name
}

output "basic_id_valid" {
  description = "Basic config template has valid ID"
  value       = netbox_config_template.basic.id != ""
}

output "complete_id" {
  description = "ID of the complete config template"
  value       = netbox_config_template.complete.id
}

output "complete_name" {
  description = "Name of the complete config template"
  value       = netbox_config_template.complete.name
}

output "complete_description" {
  description = "Description of the complete config template"
  value       = netbox_config_template.complete.description
}
