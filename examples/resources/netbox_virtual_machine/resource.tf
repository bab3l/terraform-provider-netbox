resource "netbox_cluster_type" "test" {
  name = "Test Cluster Type"
  slug = "test-cluster-type"
}

resource "netbox_cluster" "test" {
  name = "Test Cluster"
  type = netbox_cluster_type.test.id
}

resource "netbox_virtual_machine" "test" {
  name    = "test-vm-1"
  cluster = netbox_cluster.test.name
  vcpus   = 2
  memory  = 4096
  disk    = 50
  status  = "active"
}
