# Example 1: Lookup by ID
data "netbox_module" "by_id" {
  id = 1
}

output "module_by_id" {
  value       = data.netbox_module.by_id.display_name
  description = "Module display name when looked up by ID"
}

# Example 2: Lookup by device_id and module_bay_id
data "netbox_module" "by_device_and_bay" {
  device_id     = 5
  module_bay_id = 10
}

output "module_by_device_and_bay" {
  value       = data.netbox_module.by_device_and_bay.serial
  description = "Module serial when looked up by device and module bay"
}

# Example 3: Lookup by device_id and serial
data "netbox_module" "by_device_and_serial" {
  device_id = 5
  serial    = "SN123456"
}

output "module_by_device_and_serial" {
  value       = data.netbox_module.by_device_and_serial.status
  description = "Module status when looked up by device and serial"
}
