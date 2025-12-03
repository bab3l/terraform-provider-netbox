# Interface Data Source Test Outputs

# From lookup by ID
output "by_id_name" {
  value       = data.netbox_interface.by_id.name
  description = "Interface name from ID lookup"
}

output "by_id_type" {
  value       = data.netbox_interface.by_id.type
  description = "Interface type from ID lookup"
}

output "by_id_device" {
  value       = data.netbox_interface.by_id.device
  description = "Device ID from ID lookup"
}

output "by_id_enabled" {
  value       = data.netbox_interface.by_id.enabled
  description = "Enabled status from ID lookup"
}

output "by_id_mtu" {
  value       = data.netbox_interface.by_id.mtu
  description = "MTU from ID lookup"
}

output "by_id_label" {
  value       = data.netbox_interface.by_id.label
  description = "Label from ID lookup"
}

output "by_id_description" {
  value       = data.netbox_interface.by_id.description
  description = "Description from ID lookup"
}

output "by_id_mode" {
  value       = data.netbox_interface.by_id.mode
  description = "Mode from ID lookup"
}

# From lookup by device ID and name
output "by_device_id_and_name_id" {
  value       = data.netbox_interface.by_device_id_and_name.id
  description = "Interface ID from device ID + name lookup"
}

output "by_device_id_and_name_type" {
  value       = data.netbox_interface.by_device_id_and_name.type
  description = "Interface type from device ID + name lookup"
}

# From lookup by device name and interface name
output "by_device_name_and_name_id" {
  value       = data.netbox_interface.by_device_name_and_name.id
  description = "Interface ID from device name + interface name lookup"
}

output "by_device_name_and_name_type" {
  value       = data.netbox_interface.by_device_name_and_name.type
  description = "Interface type from device name + interface name lookup"
}

# Verify all lookups return the same interface
output "all_lookups_match" {
  value = (
    data.netbox_interface.by_id.id == data.netbox_interface.by_device_id_and_name.id &&
    data.netbox_interface.by_id.id == data.netbox_interface.by_device_name_and_name.id
  )
  description = "True if all lookup methods return the same interface"
}
