# Outputs to verify tenant data source lookups

# Source resource outputs
output "source_id" {
  description = "ID of the source tenant"
  value       = netbox_tenant.source.id
}

output "source_name" {
  description = "Name of the source tenant"
  value       = netbox_tenant.source.name
}

output "source_slug" {
  description = "Slug of the source tenant"
  value       = netbox_tenant.source.slug
}

output "source_group" {
  description = "Group of the source tenant"
  value       = netbox_tenant.source.group
}

# Lookup by ID outputs
output "by_id_name" {
  description = "Name from ID lookup"
  value       = data.netbox_tenant.by_id.name
}

output "by_id_slug" {
  description = "Slug from ID lookup"
  value       = data.netbox_tenant.by_id.slug
}

output "by_id_group" {
  description = "Group from ID lookup"
  value       = data.netbox_tenant.by_id.group
}

# Lookup by name outputs
output "by_name_id" {
  description = "ID from name lookup"
  value       = data.netbox_tenant.by_name.id
}

output "by_name_slug" {
  description = "Slug from name lookup"
  value       = data.netbox_tenant.by_name.slug
}

# Lookup by slug outputs
output "by_slug_id" {
  description = "ID from slug lookup"
  value       = data.netbox_tenant.by_slug.id
}

output "by_slug_name" {
  description = "Name from slug lookup"
  value       = data.netbox_tenant.by_slug.name
}

# Verification
output "all_ids_match" {
  description = "Verification that all lookup methods return the same tenant"
  value       = data.netbox_tenant.by_id.id == data.netbox_tenant.by_name.id && data.netbox_tenant.by_name.id == data.netbox_tenant.by_slug.id
}

output "group_preserved" {
  description = "Verification that group is correctly returned"
  value       = data.netbox_tenant.by_id.group == netbox_tenant.source.group
}
