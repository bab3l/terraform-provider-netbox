# Outputs to verify resource creation

output "basic_id" {
  description = "ID of the basic module bay"
  value       = netbox_module_bay.basic.id
}

output "basic_name" {
  description = "Name of the basic module bay"
  value       = netbox_module_bay.basic.name
}

output "basic_id_valid" {
  description = "Basic module bay has valid ID"
  value       = netbox_module_bay.basic.id != ""
}

output "complete_id" {
  description = "ID of the complete module bay"
  value       = netbox_module_bay.complete.id
}

output "complete_name" {
  description = "Name of the complete module bay"
  value       = netbox_module_bay.complete.name
}

output "complete_label" {
  description = "Label of the complete module bay"
  value       = netbox_module_bay.complete.label
}

output "device_id" {
  description = "ID of the parent device"
  value       = netbox_device.test.id
}
