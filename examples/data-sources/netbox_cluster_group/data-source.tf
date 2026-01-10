# Example: Look up cluster group by ID
data "netbox_cluster_group" "by_id" {
  id = "1"
}

# Example: Look up cluster group by name
data "netbox_cluster_group" "by_name" {
  name = "Production Clusters"
}

# Example: Look up cluster group by slug
data "netbox_cluster_group" "by_slug" {
  slug = "production-clusters"
}

# Example: Use cluster group data with clusters
output "cluster_group_name" {
  value = data.netbox_cluster_group.by_name.name
}

output "cluster_group_slug" {
  value = data.netbox_cluster_group.by_name.slug
}

output "cluster_group_parent" {
  value = data.netbox_cluster_group.by_name.parent
}

output "cluster_group_description" {
  value = data.netbox_cluster_group.by_name.description
}

# Access all custom fields
output "cluster_group_custom_fields" {
  value       = data.netbox_cluster_group.by_id.custom_fields
  description = "All custom fields defined in NetBox for this cluster group"
}

# Access specific custom field by name
output "cluster_group_manager" {
  value       = try([for cf in data.netbox_cluster_group.by_id.custom_fields : cf.value if cf.name == "manager_name"][0], null)
  description = "Example: accessing a text custom field for manager name"
}

output "cluster_group_priority" {
  value       = try([for cf in data.netbox_cluster_group.by_id.custom_fields : cf.value if cf.name == "priority"][0], null)
  description = "Example: accessing a numeric custom field for priority"
}
