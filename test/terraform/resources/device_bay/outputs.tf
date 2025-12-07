# Outputs to verify resource creation

output "basic_id" {
  description = "ID of the basic device bay"
  value       = netbox_device_bay.basic.id
}

output "basic_name" {
  description = "Name of the basic device bay"
  value       = netbox_device_bay.basic.name
}

output "basic_id_valid" {
  description = "Basic device bay has valid ID"
  value       = netbox_device_bay.basic.id != ""
}

output "complete_id" {
  description = "ID of the complete device bay"
  value       = netbox_device_bay.complete.id
}

output "complete_name" {
  description = "Name of the complete device bay"
  value       = netbox_device_bay.complete.name
}

output "complete_label" {
  description = "Label of the complete device bay"
  value       = netbox_device_bay.complete.label
}

output "complete_description" {
  description = "Description of the complete device bay"
  value       = netbox_device_bay.complete.description
}

output "device_id" {
  description = "ID of the parent device"
  value       = netbox_device.test.id
}
