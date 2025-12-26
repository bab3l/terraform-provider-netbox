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

output "device_role_slug" {
  value = data.netbox_device_role.by_slug.slug
}
