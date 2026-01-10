# Look up a console server port template by ID
data "netbox_console_server_port_template" "by_id" {
  id = "123"
}

# Look up a console server port template by name and device type
data "netbox_console_server_port_template" "by_device_type_and_name" {
  device_type = "456"
  name        = "CSP0"
}

# Individual attribute outputs
output "console_server_port_template_id" {
  value       = data.netbox_console_server_port_template.by_id.id
  description = "The unique ID of the console server port template"
}

output "console_server_port_template_name" {
  value       = data.netbox_console_server_port_template.by_device_type_and_name.name
  description = "The name of the console server port template"
}

output "console_server_port_template_type" {
  value       = data.netbox_console_server_port_template.by_device_type_and_name.type
  description = "The type of console server port (DB9, RJ45, etc.)"
}

output "console_server_port_template_connector" {
  value       = data.netbox_console_server_port_template.by_device_type_and_name.connector
  description = "The connector type for this console server port"
}

output "console_server_port_template_device_type" {
  value       = data.netbox_console_server_port_template.by_device_type_and_name.device_type
  description = "The device type this template belongs to"
}

# Note: Console server port templates do not support custom fields in NetBox API
output "console_server_port_template_note" {
  value       = "Console server port templates are read-only device design metadata"
  description = "Template datasources do not support custom fields"
}
