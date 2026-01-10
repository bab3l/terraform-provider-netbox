# Look up a console port template by ID
data "netbox_console_port_template" "by_id" {
  id = "123"
}

# Look up a console port template by name and device type
data "netbox_console_port_template" "by_device_type_and_name" {
  device_type = "456"
  name        = "Con0"
}

# Individual attribute outputs
output "console_port_template_id" {
  value       = data.netbox_console_port_template.by_id.id
  description = "The unique ID of the console port template"
}

output "console_port_template_name" {
  value       = data.netbox_console_port_template.by_device_type_and_name.name
  description = "The name of the console port template"
}

output "console_port_template_type" {
  value       = data.netbox_console_port_template.by_device_type_and_name.type
  description = "The type of console port (DB9, RJ45, etc.)"
}

output "console_port_template_connector" {
  value       = data.netbox_console_port_template.by_device_type_and_name.connector
  description = "The connector type for this console port"
}

output "console_port_template_device_type" {
  value       = data.netbox_console_port_template.by_device_type_and_name.device_type
  description = "The device type this template belongs to"
}

# Note: Console port templates do not support custom fields in NetBox API
output "console_port_template_note" {
  value       = "Console port templates are read-only device design metadata"
  description = "Template datasources do not support custom fields"
}
