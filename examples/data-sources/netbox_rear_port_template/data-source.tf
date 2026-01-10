# Look up a rear port template by ID
data "netbox_rear_port_template" "by_id" {
  id = 123
}

# Look up a rear port template by device type and name
data "netbox_rear_port_template" "by_device_type" {
  device_type = 456
  name        = "RearPort1"
}

# Look up a rear port template by module type and name
data "netbox_rear_port_template" "by_module_type" {
  module_type = 789
  name        = "RearPort1"
}

# Individual attribute outputs
output "rear_port_template_id" {
  value       = data.netbox_rear_port_template.by_id.id
  description = "The unique ID of the rear port template"
}

output "rear_port_template_name" {
  value       = data.netbox_rear_port_template.by_device_type.name
  description = "The name of the rear port template"
}

output "rear_port_template_type" {
  value       = data.netbox_rear_port_template.by_device_type.type
  description = "The port type (e.g., 1000base-t, sfp-plus)"
}

output "rear_port_template_color" {
  value       = data.netbox_rear_port_template.by_device_type.color
  description = "The color code for this port type"
}

output "rear_port_template_positions" {
  value       = data.netbox_rear_port_template.by_device_type.positions
  description = "Number of positions or channels for this port"
}

output "rear_port_template_device_type" {
  value       = data.netbox_rear_port_template.by_device_type.device_type
  description = "The device type this template belongs to"
}

# Note: Rear port templates do not support custom fields in NetBox API
output "rear_port_template_note" {
  value       = "Rear port templates are read-only device design metadata"
  description = "Template datasources do not support custom fields"
}
