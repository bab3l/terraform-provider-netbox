# Outputs to verify region data source lookups

# Outputs for lookup by ID
output "by_id_region_id" {
  description = "ID of the region looked up by ID"
  value       = data.netbox_region.by_id.id
}

output "by_id_region_name" {
  description = "Name of the region looked up by ID"
  value       = data.netbox_region.by_id.name
}

output "by_id_region_slug" {
  description = "Slug of the region looked up by ID"
  value       = data.netbox_region.by_id.slug
}

output "by_id_region_description" {
  description = "Description of the region looked up by ID"
  value       = data.netbox_region.by_id.description
}

# Outputs for lookup by name
output "by_name_region_id" {
  description = "ID of the region looked up by name"
  value       = data.netbox_region.by_name.id
}

output "by_name_region_name" {
  description = "Name of the region looked up by name"
  value       = data.netbox_region.by_name.name
}

output "by_name_region_slug" {
  description = "Slug of the region looked up by name"
  value       = data.netbox_region.by_name.slug
}

# Outputs for lookup by slug
output "by_slug_region_id" {
  description = "ID of the region looked up by slug"
  value       = data.netbox_region.by_slug.id
}

output "by_slug_region_name" {
  description = "Name of the region looked up by slug"
  value       = data.netbox_region.by_slug.name
}

output "by_slug_region_slug" {
  description = "Slug of the region looked up by slug"
  value       = data.netbox_region.by_slug.slug
}

# Outputs for child region lookup (verify parent relationship)
output "child_region_id" {
  description = "ID of the child region"
  value       = data.netbox_region.child_by_id.id
}

output "child_region_parent" {
  description = "Parent ID of the child region"
  value       = data.netbox_region.child_by_id.parent
}

# Verify consistency between lookup methods
output "id_matches_between_lookups" {
  description = "Verify all lookup methods return the same ID"
  value       = data.netbox_region.by_id.id == data.netbox_region.by_name.id && data.netbox_region.by_name.id == data.netbox_region.by_slug.id
}
