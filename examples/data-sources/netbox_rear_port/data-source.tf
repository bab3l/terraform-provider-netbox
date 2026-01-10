# Lookup by ID
data "netbox_rear_port" "by_id" {
  id = "123"
}

# Lookup by device_id and name
data "netbox_rear_port" "by_device_and_name" {
  device_id = "456"
  name      = "RearPort1"
}

# Use rear port data in other resources
output "rear_port_name" {
  value = data.netbox_rear_port.by_id.name
}

output "rear_port_type" {
  value = data.netbox_rear_port.by_device_and_name.type
}

output "rear_port_device" {
  value = data.netbox_rear_port.by_id.device
}

output "rear_port_positions" {
  value = data.netbox_rear_port.by_id.positions
}

output "rear_port_label" {
  value = data.netbox_rear_port.by_id.label
}

# Access all custom fields
output "rear_port_custom_fields" {
  value       = data.netbox_rear_port.by_id.custom_fields
  description = "All custom fields defined in NetBox for this rear port"
}

# Access specific custom fields by name
output "rear_port_patch_panel" {
  value       = try([for cf in data.netbox_rear_port.by_id.custom_fields : cf.value if cf.name == "patch_panel_id"][0], null)
  description = "Example: accessing a text custom field"
}

output "rear_port_managed" {
  value       = try([for cf in data.netbox_rear_port.by_id.custom_fields : cf.value if cf.name == "is_managed"][0], null)
  description = "Example: accessing a boolean custom field"
}
