# Example 1: Lookup by ID
data "netbox_module_bay" "by_id" {
  id = "1"
}

# Example 2: Lookup by device_id and name
data "netbox_module_bay" "by_device_and_name" {
  device_id = "5"
  name      = "Bay-1"
}

# Use module bay data in other resources
output "module_bay_name" {
  value = data.netbox_module_bay.by_id.name
}

output "module_bay_label" {
  value = data.netbox_module_bay.by_device_and_name.label
}

output "module_bay_position" {
  value = data.netbox_module_bay.by_id.position
}

output "module_bay_device" {
  value = data.netbox_module_bay.by_device_and_name.device_id
}

# Access all custom fields
output "module_bay_custom_fields" {
  value       = data.netbox_module_bay.by_id.custom_fields
  description = "All custom fields defined in NetBox for this module bay"
}

# Access a specific custom field by name
output "module_bay_slot_type" {
  value       = try([for cf in data.netbox_module_bay.by_id.custom_fields : cf.value if cf.name == "slot_type"][0], null)
  description = "Example: accessing a specific custom field value"
}
