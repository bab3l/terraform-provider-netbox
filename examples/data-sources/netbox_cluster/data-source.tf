data "netbox_cluster" "by_id" {
  id = "123"
}

data "netbox_cluster" "by_name" {
  name = "example-cluster"
}

output "by_id" {
  value = data.netbox_cluster.by_id.name
}

output "by_name" {
  value = data.netbox_cluster.by_name.id
}

output "cluster_type" {
  value = data.netbox_cluster.by_id.type
}

output "cluster_group" {
  value = data.netbox_cluster.by_id.group
}

output "cluster_site" {
  value = data.netbox_cluster.by_id.site
}

output "cluster_tenant" {
  value = data.netbox_cluster.by_id.tenant
}

# Access all custom fields
output "cluster_custom_fields" {
  value       = data.netbox_cluster.by_id.custom_fields
  description = "All custom fields defined in NetBox for this cluster"
}

# Access specific custom field by name
output "cluster_vcenter_id" {
  value       = try([for cf in data.netbox_cluster.by_id.custom_fields : cf.value if cf.name == "vcenter_id"][0], null)
  description = "Example: accessing a text custom field for vCenter ID"
}

output "cluster_ha_enabled" {
  value       = try([for cf in data.netbox_cluster.by_id.custom_fields : cf.value if cf.name == "ha_enabled"][0], null)
  description = "Example: accessing a boolean custom field for HA status"
}
