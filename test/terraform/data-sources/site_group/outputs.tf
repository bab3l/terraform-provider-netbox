# Outputs to verify site group data source lookups

# Source resource outputs
output "source_parent_id" {
  description = "ID of the source parent site group"
  value       = netbox_site_group.source_parent.id
}

output "source_parent_name" {
  description = "Name of the source parent site group"
  value       = netbox_site_group.source_parent.name
}

output "source_child_id" {
  description = "ID of the source child site group"
  value       = netbox_site_group.source_child.id
}

output "source_child_parent" {
  description = "Parent ID of the source child site group"
  value       = netbox_site_group.source_child.parent
}

# Lookup by ID outputs
output "by_id_name" {
  description = "Name from ID lookup"
  value       = data.netbox_site_group.by_id.name
}

output "by_id_slug" {
  description = "Slug from ID lookup"
  value       = data.netbox_site_group.by_id.slug
}

# Lookup by name outputs
output "by_name_id" {
  description = "ID from name lookup"
  value       = data.netbox_site_group.by_name.id
}

output "by_name_slug" {
  description = "Slug from name lookup"
  value       = data.netbox_site_group.by_name.slug
}

# Lookup by slug outputs
output "by_slug_id" {
  description = "ID from slug lookup"
  value       = data.netbox_site_group.by_slug.id
}

output "by_slug_name" {
  description = "Name from slug lookup"
  value       = data.netbox_site_group.by_slug.name
}

# Child lookup outputs
output "child_by_id_name" {
  description = "Name of child from ID lookup"
  value       = data.netbox_site_group.child_by_id.name
}

output "child_by_id_parent" {
  description = "Parent of child from ID lookup"
  value       = data.netbox_site_group.child_by_id.parent
}

# Verification
output "all_ids_match" {
  description = "Verification that all lookup methods return the same site group"
  value       = data.netbox_site_group.by_id.id == data.netbox_site_group.by_name.id && data.netbox_site_group.by_name.id == data.netbox_site_group.by_slug.id
}

output "child_parent_matches" {
  description = "Verification that child's parent matches source"
  value       = data.netbox_site_group.child_by_id.parent == netbox_site_group.source_parent.id
}
