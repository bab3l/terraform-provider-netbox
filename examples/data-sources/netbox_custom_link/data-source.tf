# Look up by ID
data "netbox_custom_link" "by_id" {
  id = "123"
}

# Look up by name
data "netbox_custom_link" "by_name" {
  name = "device_documentation"
}

output "custom_link_id" {
  value = data.netbox_custom_link.by_id.id
}

output "custom_link_name" {
  value = data.netbox_custom_link.by_name.name
}

output "custom_link_url" {
  value = data.netbox_custom_link.by_name.url
}

output "custom_link_text" {
  value = data.netbox_custom_link.by_name.text
}

output "custom_link_object_types" {
  value = data.netbox_custom_link.by_name.object_types
}

# Note: Custom links are metadata definitions and do not support custom fields in NetBox API
output "custom_link_note" {
  value       = "Custom link definitions are read-only metadata"
  description = "Custom links define external links to be displayed on objects"
}
