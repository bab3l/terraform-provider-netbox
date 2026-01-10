# Example: Look up a device bay by ID
data "netbox_device_bay" "by_id" {
  id = "1"
}

# Example: Look up a device bay by device and name
# This is useful when you have the parent device and bay name
data "netbox_device_bay" "by_device_and_name" {
  device = "5" # Device ID
  name   = "Bay 1"
}

# Example: Use device bay data in other resources
output "device_bay_id" {
  value = data.netbox_device_bay.by_id.id
}

output "device_bay_name" {
  value = data.netbox_device_bay.by_device_and_name.name
}

output "device_bay_device" {
  value = data.netbox_device_bay.by_device_and_name.device
}

output "device_bay_installed_device" {
  value = data.netbox_device_bay.by_id.installed_device
}

# Access all custom fields
output "device_bay_custom_fields" {
  value       = data.netbox_device_bay.by_id.custom_fields
  description = "All custom fields defined in NetBox for this device bay"
}

# Access a specific custom field by name
output "device_bay_slot_type" {
  value       = try([for cf in data.netbox_device_bay.by_id.custom_fields : cf.value if cf.name == "slot_type"][0], null)
  description = "Example: accessing a specific custom field value"
}
