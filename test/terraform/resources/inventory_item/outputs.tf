# Outputs to verify resource creation

output "basic_id" {
  description = "ID of the basic inventory item"
  value       = netbox_inventory_item.basic.id
}

output "basic_name" {
  description = "Name of the basic inventory item"
  value       = netbox_inventory_item.basic.name
}

output "basic_id_valid" {
  description = "Basic inventory item has valid ID"
  value       = netbox_inventory_item.basic.id != ""
}

output "complete_id" {
  description = "ID of the complete inventory item"
  value       = netbox_inventory_item.complete.id
}

output "complete_name" {
  description = "Name of the complete inventory item"
  value       = netbox_inventory_item.complete.name
}

output "complete_label" {
  description = "Label of the complete inventory item"
  value       = netbox_inventory_item.complete.label
}

output "complete_serial" {
  description = "Serial of the complete inventory item"
  value       = netbox_inventory_item.complete.serial
}

output "device_id" {
  description = "ID of the parent device"
  value       = netbox_device.test.id
}
