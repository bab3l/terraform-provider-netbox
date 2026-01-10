# Look up a device bay template by ID
data "netbox_device_bay_template" "by_id" {
  id = 1
}

# Look up a device bay template by name (requires device_type for uniqueness)
data "netbox_device_bay_template" "by_name" {
  name        = "Bay 1"
  device_type = "123"
}

# Individual attribute outputs
output "device_bay_template_id" {
  value       = data.netbox_device_bay_template.by_id.id
  description = "The unique ID of the device bay template"
}

output "device_bay_template_name" {
  value       = data.netbox_device_bay_template.by_name.name
  description = "The name of the device bay template"
}

output "device_bay_template_device_type" {
  value       = data.netbox_device_bay_template.by_name.device_type
  description = "The device type this template belongs to"
}

output "device_bay_template_device_type_name" {
  value       = data.netbox_device_bay_template.by_name.device_type_name
  description = "The name of the device type"
}

output "device_bay_template_label" {
  value       = data.netbox_device_bay_template.by_name.label
  description = "The label or display name for this bay"
}

output "device_bay_template_description" {
  value       = data.netbox_device_bay_template.by_name.description
  description = "Description of the device bay template"
}

# Note: Device bay templates do not support custom fields in NetBox API
output "device_bay_template_note" {
  value       = "Device bay templates are read-only device design metadata"
  description = "Template datasources do not support custom fields"
}
