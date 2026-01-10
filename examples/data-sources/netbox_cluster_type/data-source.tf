data "netbox_cluster_type" "by_id" {
  id = "123"
}

data "netbox_cluster_type" "by_name" {
  name = "VMware vSphere"
}

data "netbox_cluster_type" "by_slug" {
  slug = "vmware-vsphere"
}

output "by_id" {
  value = data.netbox_cluster_type.by_id.name
}

output "by_name" {
  value = data.netbox_cluster_type.by_name.id
}

output "by_slug" {
  value = data.netbox_cluster_type.by_slug.id
}

output "cluster_type_slug" {
  value = data.netbox_cluster_type.by_id.slug
}

output "cluster_type_description" {
  value = data.netbox_cluster_type.by_id.description
}

# Access all custom fields
output "cluster_type_custom_fields" {
  value       = data.netbox_cluster_type.by_id.custom_fields
  description = "All custom fields defined in NetBox for this cluster type"
}

# Access specific custom field by name
output "cluster_type_hypervisor" {
  value       = try([for cf in data.netbox_cluster_type.by_id.custom_fields : cf.value if cf.name == "hypervisor_version"][0], null)
  description = "Example: accessing a text custom field for hypervisor version"
}

output "cluster_type_requires_license" {
  value       = try([for cf in data.netbox_cluster_type.by_id.custom_fields : cf.value if cf.name == "requires_license"][0], null)
  description = "Example: accessing a boolean custom field for license requirement"
}
