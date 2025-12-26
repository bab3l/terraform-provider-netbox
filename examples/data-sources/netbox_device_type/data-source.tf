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
