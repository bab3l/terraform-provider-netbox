resource "netbox_cluster_type" "test" {
  name = "Test Cluster Type"
  slug = "test-cluster-type"
}

resource "netbox_cluster" "test" {
  name            = "Test Cluster"
  cluster_type_id = netbox_cluster_type.test.id
}

resource "netbox_virtual_machine" "test" {
  name       = "Test VM"
  cluster_id = netbox_cluster.test.id
}

resource "netbox_vm_interface" "test" {
  name               = "eth0"
  virtual_machine_id = netbox_virtual_machine.test.id
}
