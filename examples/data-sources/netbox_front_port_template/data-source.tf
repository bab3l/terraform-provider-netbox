# Look up a front port template by ID
data "netbox_front_port_template" "by_id" {
  id = 1
}

# Look up a front port template by name and device type
data "netbox_front_port_template" "by_device_type" {
  name        = "GigabitEthernet"
  device_type = 5
}

# Look up a front port template by name and module type
data "netbox_front_port_template" "by_module_type" {
  name        = "SFP"
  module_type = 10
}

# Individual attribute outputs
output "front_port_template_id" {
  value       = data.netbox_front_port_template.by_id.id
  description = "The unique ID of the front port template"
}

output "front_port_template_name" {
  value       = data.netbox_front_port_template.by_device_type.name
  description = "The name of the front port template"
}

output "front_port_template_type" {
  value       = data.netbox_front_port_template.by_device_type.type
  description = "The port type (e.g., 1000base-t, sfp-plus)"
}

output "front_port_template_connector" {
  value       = data.netbox_front_port_template.by_device_type.connector
  description = "The connector type (e.g., RJ45, SFP)"
}

output "front_port_template_color" {
  value       = data.netbox_front_port_template.by_device_type.color
  description = "The color code for this port type"
}

output "front_port_template_device_type" {
  value       = data.netbox_front_port_template.by_device_type.device_type
  description = "The device type this template belongs to"
}

# Note: Front port templates do not support custom fields in NetBox API
output "front_port_template_note" {
  value       = "Front port templates are read-only device design metadata"
  description = "Template datasources do not support custom fields"
}
