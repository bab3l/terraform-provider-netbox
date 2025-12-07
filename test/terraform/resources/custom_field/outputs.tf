# Outputs to verify resource creation

output "basic_id" {
  description = "ID of the basic custom field"
  value       = netbox_custom_field.basic.id
}

output "basic_name" {
  description = "Name of the basic custom field"
  value       = netbox_custom_field.basic.name
}

output "basic_type" {
  description = "Type of the basic custom field"
  value       = netbox_custom_field.basic.type
}

output "basic_id_valid" {
  description = "Basic custom field has valid ID"
  value       = netbox_custom_field.basic.id != ""
}

output "complete_id" {
  description = "ID of the complete custom field"
  value       = netbox_custom_field.complete.id
}

output "complete_name" {
  description = "Name of the complete custom field"
  value       = netbox_custom_field.complete.name
}

output "complete_label" {
  description = "Label of the complete custom field"
  value       = netbox_custom_field.complete.label
}

output "complete_description" {
  description = "Description of the complete custom field"
  value       = netbox_custom_field.complete.description
}
