# Example 1: Lookup by ID
data "netbox_power_outlet" "by_id" {
  id = "1"
}

# Example 2: Lookup by device_id and name
data "netbox_power_outlet" "by_device_and_name" {
  device_id = "5"
  name      = "PSU-1"
}

# Use power outlet data in other resources
output "power_outlet_name" {
  value = data.netbox_power_outlet.by_id.name
}

output "power_outlet_type" {
  value = data.netbox_power_outlet.by_device_and_name.type
}

output "power_outlet_device" {
  value = data.netbox_power_outlet.by_id.device
}

output "power_outlet_power_port" {
  value = data.netbox_power_outlet.by_id.power_port
}

# Access all custom fields
output "power_outlet_custom_fields" {
  value       = data.netbox_power_outlet.by_id.custom_fields
  description = "All custom fields defined in NetBox for this power outlet"
}

# Access specific custom fields by name
output "power_outlet_max_draw" {
  value       = try([for cf in data.netbox_power_outlet.by_id.custom_fields : cf.value if cf.name == "max_draw_watts"][0], null)
  description = "Example: accessing a numeric custom field"
}

output "power_outlet_monitored" {
  value       = try([for cf in data.netbox_power_outlet.by_id.custom_fields : cf.value if cf.name == "is_monitored"][0], null)
  description = "Example: accessing a boolean custom field"
}
