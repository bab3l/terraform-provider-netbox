# Outputs to verify resource creation

output "basic_id" {
  description = "ID of the basic interface template"
  value       = netbox_interface_template.basic.id
}

output "basic_name" {
  description = "Name of the basic interface template"
  value       = netbox_interface_template.basic.name
}

output "basic_type" {
  description = "Type of the basic interface template"
  value       = netbox_interface_template.basic.type
}

output "basic_id_valid" {
  description = "Basic interface template has valid ID"
  value       = netbox_interface_template.basic.id != ""
}

output "complete_id" {
  description = "ID of the complete interface template"
  value       = netbox_interface_template.complete.id
}

output "complete_name" {
  description = "Name of the complete interface template"
  value       = netbox_interface_template.complete.name
}

output "complete_label" {
  description = "Label of the complete interface template"
  value       = netbox_interface_template.complete.label
}

output "device_type_id" {
  description = "ID of the parent device type"
  value       = netbox_device_type.test.id
}
