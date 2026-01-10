# Example: Look up a device role by ID
data "netbox_device_role" "by_id" {
  id = "1"
}

# Example: Look up a device role by name
data "netbox_device_role" "by_name" {
  name = "Router"
}

# Example: Look up a device role by slug
data "netbox_device_role" "by_slug" {
  slug = "router"
}

# Example: Use device role data in other resources
output "device_role_id" {
  value = data.netbox_device_role.by_id.id
}

output "device_role_name" {
  value = data.netbox_device_role.by_name.name
}

output "device_role_color" {
  value = data.netbox_device_role.by_slug.color
}

# Access all custom fields
output "device_role_custom_fields" {
  value       = data.netbox_device_role.by_id.custom_fields
  description = "All custom fields defined in NetBox for this device role"
}

# Access a specific custom field by name
output "device_role_priority" {
  value       = try([for cf in data.netbox_device_role.by_id.custom_fields : cf.value if cf.name == "priority"][0], null)
  description = "Example: accessing a specific custom field value"
}
