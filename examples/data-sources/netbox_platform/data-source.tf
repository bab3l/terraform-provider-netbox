# Example 1: Lookup by ID
data "netbox_platform" "by_id" {
  id = "1"
}

# Example 2: Lookup by slug
data "netbox_platform" "by_slug" {
  slug = "linux"
}

# Example 3: Lookup by name
data "netbox_platform" "by_name" {
  name = "Linux"
}

# Use platform data in other resources
output "platform_name" {
  value = data.netbox_platform.by_id.name
}

output "platform_slug" {
  value = data.netbox_platform.by_slug.slug
}

output "platform_description" {
  value = data.netbox_platform.by_name.description
}

output "platform_manufacturer" {
  value = data.netbox_platform.by_id.manufacturer
}

# Access all custom fields
output "platform_custom_fields" {
  value       = data.netbox_platform.by_id.custom_fields
  description = "All custom fields defined in NetBox for this platform"
}

# Access specific custom fields by name
output "platform_config_template" {
  value       = try([for cf in data.netbox_platform.by_id.custom_fields : cf.value if cf.name == "config_template"][0], null)
  description = "Example: accessing a text custom field"
}

output "platform_support_level" {
  value       = try([for cf in data.netbox_platform.by_id.custom_fields : cf.value if cf.name == "support_level"][0], null)
  description = "Example: accessing a select custom field"
}
