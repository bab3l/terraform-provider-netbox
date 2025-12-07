# Outputs to verify resource creation

output "basic_id" {
  description = "ID of the basic console server port"
  value       = netbox_console_server_port.basic.id
}

output "basic_name" {
  description = "Name of the basic console server port"
  value       = netbox_console_server_port.basic.name
}

output "basic_id_valid" {
  description = "Basic console server port has valid ID"
  value       = netbox_console_server_port.basic.id != ""
}

output "complete_id" {
  description = "ID of the complete console server port"
  value       = netbox_console_server_port.complete.id
}

output "complete_name" {
  description = "Name of the complete console server port"
  value       = netbox_console_server_port.complete.name
}

output "complete_label" {
  description = "Label of the complete console server port"
  value       = netbox_console_server_port.complete.label
}

output "device_id" {
  description = "ID of the parent device"
  value       = netbox_device.test.id
}
