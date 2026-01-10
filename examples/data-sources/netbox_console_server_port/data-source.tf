# Look up a console server port by ID
data "netbox_console_server_port" "by_id" {
  id = "789"
}

# Look up a console server port by device and name
data "netbox_console_server_port" "by_device_and_name" {
  device_id = "456"
  name      = "csp0"
}

# Access standard attributes
output "console_server_port_name" {
  value = data.netbox_console_server_port.by_id.name
}

output "console_server_port_type" {
  value = data.netbox_console_server_port.by_id.type
}

output "console_server_port_device" {
  value = data.netbox_console_server_port.by_device_and_name.device_id
}

# Access all custom fields
output "console_server_port_custom_fields" {
  value       = data.netbox_console_server_port.by_id.custom_fields
  description = "All custom fields defined in NetBox for this console server port"
}

# Access a specific custom field by name
output "console_server_port_terminal_server" {
  value       = try([for cf in data.netbox_console_server_port.by_id.custom_fields : cf.value if cf.name == "terminal_server"][0], null)
  description = "Example: accessing a specific custom field value"
}
