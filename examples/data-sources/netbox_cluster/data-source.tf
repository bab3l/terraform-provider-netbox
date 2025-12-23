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
