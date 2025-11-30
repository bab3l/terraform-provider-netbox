# Outputs to verify data source lookups

# Source resource outputs (for comparison)
output "source_id" {
  description = "ID of the source site resource"
  value       = netbox_site.source.id
}

output "source_name" {
  description = "Name of the source site resource"
  value       = netbox_site.source.name
}

output "source_slug" {
  description = "Slug of the source site resource"
  value       = netbox_site.source.slug
}

# Lookup by ID outputs
output "by_id_name" {
  description = "Name from ID lookup"
  value       = data.netbox_site.by_id.name
}

output "by_id_slug" {
  description = "Slug from ID lookup"
  value       = data.netbox_site.by_id.slug
}

output "by_id_status" {
  description = "Status from ID lookup"
  value       = data.netbox_site.by_id.status
}

# Lookup by name outputs
output "by_name_id" {
  description = "ID from name lookup"
  value       = data.netbox_site.by_name.id
}

output "by_name_slug" {
  description = "Slug from name lookup"
  value       = data.netbox_site.by_name.slug
}

# Lookup by slug outputs
output "by_slug_id" {
  description = "ID from slug lookup"
  value       = data.netbox_site.by_slug.id
}

output "by_slug_name" {
  description = "Name from slug lookup"
  value       = data.netbox_site.by_slug.name
}

# Verification: All lookups should return the same ID
output "all_ids_match" {
  description = "Verification that all lookup methods return the same site"
  value       = data.netbox_site.by_id.id == data.netbox_site.by_name.id && data.netbox_site.by_name.id == data.netbox_site.by_slug.id
}
