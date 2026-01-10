# Look up a power outlet template by ID
data "netbox_power_outlet_template" "by_id" {
  id = 1
}

# Look up a power outlet template by device type and name
data "netbox_power_outlet_template" "by_device_type" {
  device_type = 10
  name        = "PSU"
}

# Look up a power outlet template by module type and name
data "netbox_power_outlet_template" "by_module_type" {
  module_type = 5
  name        = "PSU"
}

# Individual attribute outputs
output "power_outlet_template_id" {
  value       = data.netbox_power_outlet_template.by_id.id
  description = "The unique ID of the power outlet template"
}

output "power_outlet_template_name" {
  value       = data.netbox_power_outlet_template.by_device_type.name
  description = "Power outlet template name"
}

output "power_outlet_template_type" {
  value       = data.netbox_power_outlet_template.by_device_type.type
  description = "The outlet type (e.g., IEC 60320 C13, C19)"
}

output "power_outlet_template_feed_leg" {
  value       = data.netbox_power_outlet_template.by_device_type.feed_leg
  description = "The feed leg (e.g., A, B, C) for this outlet"
}

output "power_outlet_template_label" {
  value       = data.netbox_power_outlet_template.by_device_type.label
  description = "The label or display name for this outlet"
}

output "power_outlet_template_device_type" {
  value       = data.netbox_power_outlet_template.by_device_type.device_type
  description = "The device type this template belongs to"
}

# Note: Power outlet templates do not support custom fields in NetBox API
output "power_outlet_template_note" {
  value       = "Power outlet templates are read-only device design metadata"
  description = "Template datasources do not support custom fields"
}
