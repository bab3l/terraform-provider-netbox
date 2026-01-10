# Lookup by ID
data "netbox_power_port" "by_id" {
  id = "123"
}

# Lookup by device_id and name
data "netbox_power_port" "by_device_and_name" {
  device_id = "456"
  name      = "PWR1"
}

# Use power port data in other resources
output "power_port_name" {
  value = data.netbox_power_port.by_id.name
}

output "power_port_type" {
  value = data.netbox_power_port.by_device_and_name.type
}

output "power_port_device" {
  value = data.netbox_power_port.by_id.device
}

output "power_port_maximum_draw" {
  value = data.netbox_power_port.by_id.maximum_draw
}

output "power_port_allocated_draw" {
  value = data.netbox_power_port.by_id.allocated_draw
}

# Access all custom fields
output "power_port_custom_fields" {
  value       = data.netbox_power_port.by_id.custom_fields
  description = "All custom fields defined in NetBox for this power port"
}

# Access specific custom fields by name
output "power_port_pdu_outlet" {
  value       = try([for cf in data.netbox_power_port.by_id.custom_fields : cf.value if cf.name == "pdu_outlet_number"][0], null)
  description = "Example: accessing a text custom field"
}

output "power_port_monitored" {
  value       = try([for cf in data.netbox_power_port.by_id.custom_fields : cf.value if cf.name == "is_monitored"][0], null)
  description = "Example: accessing a boolean custom field"
}
