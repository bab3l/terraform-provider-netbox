# Outputs to verify resource creation

output "basic_id" {
  description = "ID of the basic power outlet"
  value       = netbox_power_outlet.basic.id
}

output "basic_name" {
  description = "Name of the basic power outlet"
  value       = netbox_power_outlet.basic.name
}

output "basic_id_valid" {
  description = "Basic power outlet has valid ID"
  value       = netbox_power_outlet.basic.id != ""
}

output "complete_id" {
  description = "ID of the complete power outlet"
  value       = netbox_power_outlet.complete.id
}

output "complete_name" {
  description = "Name of the complete power outlet"
  value       = netbox_power_outlet.complete.name
}

output "complete_label" {
  description = "Label of the complete power outlet"
  value       = netbox_power_outlet.complete.label
}

output "device_id" {
  description = "ID of the parent device"
  value       = netbox_device.test.id
}
