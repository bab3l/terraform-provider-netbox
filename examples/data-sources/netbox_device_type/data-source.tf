# Example: Look up a device type by ID
data "netbox_device_type" "by_id" {
  id = "1"
}

# Example: Look up a device type by slug
data "netbox_device_type" "by_slug" {
  slug = "catalyst-3750"
}

# Example: Look up a device type by model name
data "netbox_device_type" "by_model" {
  model = "Catalyst 3750"
}

# Example: Use device type data in other resources
output "device_type_id" {
  value = data.netbox_device_type.by_id.id
}

output "device_type_model" {
  value = data.netbox_device_type.by_model.model
}

output "device_type_manufacturer" {
  value = data.netbox_device_type.by_model.manufacturer
}

output "device_type_u_height" {
  value = data.netbox_device_type.by_id.u_height
}

# Access all custom fields
output "device_type_custom_fields" {
  value       = data.netbox_device_type.by_id.custom_fields
  description = "All custom fields defined in NetBox for this device type"
}

# Access specific custom fields by name
output "device_type_power_draw" {
  value       = try([for cf in data.netbox_device_type.by_id.custom_fields : cf.value if cf.name == "power_draw"][0], null)
  description = "Example: accessing a numeric custom field"
}

output "device_type_end_of_life" {
  value       = try([for cf in data.netbox_device_type.by_id.custom_fields : cf.value if cf.name == "end_of_life"][0], null)
  description = "Example: accessing a date custom field"
}
