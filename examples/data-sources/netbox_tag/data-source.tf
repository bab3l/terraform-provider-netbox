# Look up a tag by ID
data "netbox_tag" "by_id" {
  id = "1"
}

# Look up a tag by name
data "netbox_tag" "by_name" {
  name = "Production"
}

# Look up a tag by slug
data "netbox_tag" "by_slug" {
  slug = "production"
}

# Output tag details
output "tag_id" {
  value = data.netbox_tag.by_id.id
}

output "tag_name" {
  value = data.netbox_tag.by_name.name
}

output "tag_slug" {
  value = data.netbox_tag.by_slug.slug
}

output "tag_color" {
  value = data.netbox_tag.by_name.color
}

output "tag_description" {
  value = data.netbox_tag.by_name.description
}

# Note: Tags do not support custom fields in NetBox API
output "tag_note" {
  value       = "Tags are read-only organizational metadata"
  description = "Tags provide simple labels for organizing and categorizing objects"
}
