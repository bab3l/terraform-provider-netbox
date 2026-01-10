# Example 1: Lookup by ID
data "netbox_module" "by_id" {
  id = "1"
}

# Example 2: Lookup by device_id and module_bay_id
data "netbox_module" "by_device_and_bay" {
  device_id     = "5"
  module_bay_id = "10"
}

# Example 3: Lookup by device_id and serial
data "netbox_module" "by_device_and_serial" {
  device_id = "5"
  serial    = "SN123456"
}

# Use module data in other resources
output "module_serial" {
  value = data.netbox_module.by_id.serial
}

output "module_status" {
  value = data.netbox_module.by_device_and_bay.status
}

output "module_type" {
  value = data.netbox_module.by_id.module_type
}

output "module_device" {
  value = data.netbox_module.by_id.device
}

# Access all custom fields
output "module_custom_fields" {
  value       = data.netbox_module.by_id.custom_fields
  description = "All custom fields defined in NetBox for this module"
}

# Access specific custom fields by name
output "module_firmware_version" {
  value       = try([for cf in data.netbox_module.by_id.custom_fields : cf.value if cf.name == "firmware_version"][0], null)
  description = "Example: accessing a text custom field"
}

output "module_install_date" {
  value       = try([for cf in data.netbox_module.by_id.custom_fields : cf.value if cf.name == "install_date"][0], null)
  description = "Example: accessing a date custom field"
}
