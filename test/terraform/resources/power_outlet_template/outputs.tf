# Outputs to verify resource creation

output "basic_id" {
  description = "ID of the basic power outlet template"
  value       = netbox_power_outlet_template.basic.id
}

output "basic_name" {
  description = "Name of the basic power outlet template"
  value       = netbox_power_outlet_template.basic.name
}

output "basic_id_valid" {
  description = "Basic power outlet template has valid ID"
  value       = netbox_power_outlet_template.basic.id != ""
}

output "complete_id" {
  description = "ID of the complete power outlet template"
  value       = netbox_power_outlet_template.complete.id
}

output "complete_name" {
  description = "Name of the complete power outlet template"
  value       = netbox_power_outlet_template.complete.name
}

output "complete_label" {
  description = "Label of the complete power outlet template"
  value       = netbox_power_outlet_template.complete.label
}

output "device_type_id" {
  description = "ID of the parent device type"
  value       = netbox_device_type.test.id
}
