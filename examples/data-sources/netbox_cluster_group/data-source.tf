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

output "cluster_group_description" {
  value = data.netbox_cluster_group.by_name.description
}
