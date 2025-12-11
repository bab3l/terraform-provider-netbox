output "basic_id" {
  description = "ID of the basic rear port"
  value       = netbox_rear_port.basic.id
}

output "basic_type" {
  description = "Type of the basic rear port"
  value       = netbox_rear_port.basic.type
}

output "basic_positions" {
  description = "Positions of the basic rear port"
  value       = netbox_rear_port.basic.positions
}

output "basic_id_valid" {
  description = "Basic rear port has valid ID"
  value       = netbox_rear_port.basic.id != ""
}

output "complete_id" {
  description = "ID of the complete rear port"
  value       = netbox_rear_port.complete.id
}

output "complete_label" {
  description = "Label of the complete rear port"
  value       = netbox_rear_port.complete.label
}

output "complete_color" {
  description = "Color of the complete rear port"
  value       = netbox_rear_port.complete.color
}

output "complete_positions" {
  description = "Positions count of the complete rear port"
  value       = netbox_rear_port.complete.positions
}

output "complete_description" {
  description = "Description of the complete rear port"
  value       = netbox_rear_port.complete.description
}

output "complete_mark_connected" {
  description = "Mark connected flag of the complete rear port"
  value       = netbox_rear_port.complete.mark_connected
}

output "device_id" {
  description = "ID of the parent device"
  value       = netbox_device.test.id
}

output "site_id" {
  description = "ID of the parent site"
  value       = netbox_site.test.id
}
