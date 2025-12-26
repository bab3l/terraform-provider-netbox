# Example 1: Lookup by ID
data "netbox_platform" "by_id" {
  id = 1
}

output "platform_by_id" {
  value       = data.netbox_platform.by_id.name
  description = "Platform name when looked up by ID"
}

# Example 2: Lookup by slug
data "netbox_platform" "by_slug" {
  slug = "linux"
}

output "platform_by_slug" {
  value       = data.netbox_platform.by_slug.display_name
  description = "Platform display name when looked up by slug"
}

# Example 3: Lookup by name
data "netbox_platform" "by_name" {
  name = "Linux"
}

output "platform_by_name" {
  value       = data.netbox_platform.by_name.description
  description = "Platform description when looked up by name"
}
