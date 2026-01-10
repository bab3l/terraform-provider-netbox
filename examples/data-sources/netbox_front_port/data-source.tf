# Example: Look up a front port by ID
data "netbox_front_port" "by_id" {
  id = "1"
}

# Example: Look up a front port by device_id and name
data "netbox_front_port" "by_device_and_name" {
  device_id = "5"
  name      = "eth0"
}

# Example: Use front port data in other resources
output "front_port_id" {
  value = data.netbox_front_port.by_id.id
}

output "front_port_name" {
  value = data.netbox_front_port.by_device_and_name.name
}

output "front_port_type" {
  value = data.netbox_front_port.by_id.type
}

output "front_port_rear_port" {
  value = data.netbox_front_port.by_device_and_name.rear_port
}

# Access all custom fields
output "front_port_custom_fields" {
  value       = data.netbox_front_port.by_id.custom_fields
  description = "All custom fields defined in NetBox for this front port"
}

# Access a specific custom field by name
output "front_port_cable_type" {
  value       = try([for cf in data.netbox_front_port.by_id.custom_fields : cf.value if cf.name == "cable_type"][0], null)
  description = "Example: accessing a specific custom field value"
}
