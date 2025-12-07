# Outputs to verify resource creation

output "basic_id" {
  description = "ID of the basic console port template"
  value       = netbox_console_port_template.basic.id
}

output "basic_name" {
  description = "Name of the basic console port template"
  value       = netbox_console_port_template.basic.name
}

output "basic_id_valid" {
  description = "Basic console port template has valid ID"
  value       = netbox_console_port_template.basic.id != ""
}

output "complete_id" {
  description = "ID of the complete console port template"
  value       = netbox_console_port_template.complete.id
}

output "complete_name" {
  description = "Name of the complete console port template"
  value       = netbox_console_port_template.complete.name
}

output "complete_label" {
  description = "Label of the complete console port template"
  value       = netbox_console_port_template.complete.label
}

output "device_type_id" {
  description = "ID of the parent device type"
  value       = netbox_device_type.test.id
}
