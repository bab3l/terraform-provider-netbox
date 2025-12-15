data "netbox_cluster" "test" {
  name = "test-cluster"
}

output "example" {
  value = data.netbox_cluster.test.id
}
