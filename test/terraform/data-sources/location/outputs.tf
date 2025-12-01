# Outputs to verify location data source lookups

# Outputs for lookup by ID
output "by_id_location_id" {
  description = "ID of the location looked up by ID"
  value       = data.netbox_location.by_id.id
}

output "by_id_location_name" {
  description = "Name of the location looked up by ID"
  value       = data.netbox_location.by_id.name
}

output "by_id_location_slug" {
  description = "Slug of the location looked up by ID"
  value       = data.netbox_location.by_id.slug
}

output "by_id_location_site" {
  description = "Site ID of the location looked up by ID"
  value       = data.netbox_location.by_id.site
}

output "by_id_location_description" {
  description = "Description of the location looked up by ID"
  value       = data.netbox_location.by_id.description
}

output "by_id_location_status" {
  description = "Status of the location looked up by ID"
  value       = data.netbox_location.by_id.status
}

# Outputs for lookup by name
output "by_name_location_id" {
  description = "ID of the location looked up by name"
  value       = data.netbox_location.by_name.id
}

output "by_name_location_name" {
  description = "Name of the location looked up by name"
  value       = data.netbox_location.by_name.name
}

output "by_name_location_slug" {
  description = "Slug of the location looked up by name"
  value       = data.netbox_location.by_name.slug
}

# Outputs for lookup by slug
output "by_slug_location_id" {
  description = "ID of the location looked up by slug"
  value       = data.netbox_location.by_slug.id
}

output "by_slug_location_name" {
  description = "Name of the location looked up by slug"
  value       = data.netbox_location.by_slug.name
}

output "by_slug_location_slug" {
  description = "Slug of the location looked up by slug"
  value       = data.netbox_location.by_slug.slug
}

# Outputs for child location lookup (verify parent relationship)
output "child_location_id" {
  description = "ID of the child location"
  value       = data.netbox_location.child_by_id.id
}

output "child_location_parent" {
  description = "Parent ID of the child location"
  value       = data.netbox_location.child_by_id.parent
}

output "child_location_site" {
  description = "Site ID of the child location"
  value       = data.netbox_location.child_by_id.site
}

# Verify consistency between lookup methods
output "id_matches_between_lookups" {
  description = "Verify all lookup methods return the same ID"
  value       = data.netbox_location.by_id.id == data.netbox_location.by_name.id && data.netbox_location.by_name.id == data.netbox_location.by_slug.id
}
