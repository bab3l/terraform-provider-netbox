# Look up a power port template by ID
data "netbox_power_port_template" "by_id" {
  id = 123
}

# Look up a power port template by device type and name
data "netbox_power_port_template" "by_device_type" {
  device_type = 456
  name        = "PWR1"
}

# Look up a power port template by module type and name
data "netbox_power_port_template" "by_module_type" {
  module_type = 789
  name        = "PWR1"
}

# Individual attribute outputs
output "power_port_template_id" {
  value       = data.netbox_power_port_template.by_id.id
  description = "The unique ID of the power port template"
}

output "power_port_template_name" {
  value       = data.netbox_power_port_template.by_device_type.name
  description = "The name of the power port template"
}

output "power_port_template_type" {
  value       = data.netbox_power_port_template.by_device_type.type
  description = "The port type (e.g., IEC 60320 C14, C20)"
}

output "power_port_template_feed_leg" {
  value       = data.netbox_power_port_template.by_device_type.feed_leg
  description = "The feed leg (e.g., A, B, C) for this port"
}

output "power_port_template_device_type" {
  value       = data.netbox_power_port_template.by_device_type.device_type
  description = "The device type this template belongs to"
}

# Note: Power port templates do not support custom fields in NetBox API
output "power_port_template_note" {
  value       = "Power port templates are read-only device design metadata"
  description = "Template datasources do not support custom fields"
}
