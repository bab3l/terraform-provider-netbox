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
