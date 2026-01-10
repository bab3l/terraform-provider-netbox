# Look up a console port by ID
data "netbox_console_port" "by_id" {
  id = "789"
}

# Look up a console port by device and name
data "netbox_console_port" "by_device_and_name" {
  device_id = "456"
  name      = "con0"
}

# Access standard attributes
output "console_port_name" {
  value = data.netbox_console_port.by_id.name
}

output "console_port_type" {
  value = data.netbox_console_port.by_id.type
}

output "console_port_device" {
  value = data.netbox_console_port.by_device_and_name.device_id
}

# Access all custom fields
output "console_port_custom_fields" {
  value       = data.netbox_console_port.by_id.custom_fields
  description = "All custom fields defined in NetBox for this console port"
}

# Access a specific custom field by name
output "console_port_management_vlan" {
  value       = try([for cf in data.netbox_console_port.by_id.custom_fields : cf.value if cf.name == "management_vlan"][0], null)
  description = "Example: accessing a specific custom field value"
}
