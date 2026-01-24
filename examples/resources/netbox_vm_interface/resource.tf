resource "netbox_cluster_type" "test" {
  name = "Test Cluster Type"
  slug = "test-cluster-type"
}

resource "netbox_cluster" "test" {
  name = "Test Cluster"
  type = netbox_cluster_type.test.name
}

resource "netbox_virtual_machine" "test" {
  name    = "Test VM"
  cluster = netbox_cluster.test.name
}

resource "netbox_vm_interface" "test" {
  name            = "eth0"
  virtual_machine = netbox_virtual_machine.test.name
}

# VM interface linked using VM ID
resource "netbox_vm_interface" "test_by_id" {
  name            = "eth1"
  virtual_machine = netbox_virtual_machine.test.id
}

# Optional: seed owned custom fields during import
import {
  to = netbox_vm_interface.test
  id = "123"

  identity = {
    custom_fields = [
      "interface_role:text",
      "mtu_override:integer",
    ]
  }
}
