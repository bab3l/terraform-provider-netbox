data "netbox_cluster_type" "test" {
  name = "test-cluster-type"
}

output "example" {
  value = data.netbox_cluster_type.test.id
}
