resource "netbox_cluster_type" "test" {
  name = "Primary IP Cluster Type"
  slug = "primary-ip-cluster-type"
}

resource "netbox_cluster" "test" {
  name = "Primary IP Cluster"
  type = netbox_cluster_type.test.id
}

resource "netbox_virtual_machine" "test" {
  name    = "primary-ip-vm-1"
  cluster = netbox_cluster.test.name
  status  = "active"
}

resource "netbox_vm_interface" "test" {
  name            = "eth0"
  virtual_machine = netbox_virtual_machine.test.name
}

resource "netbox_ip_address" "test" {
  address              = "192.0.2.10/24"
  status               = "active"
  assigned_object_type = "virtualization.vminterface"
  assigned_object_id   = netbox_vm_interface.test.id
}

resource "netbox_virtual_machine_primary_ip" "test" {
  virtual_machine = netbox_virtual_machine.test.name
  primary_ip4     = netbox_ip_address.test.id
}

import {
  to = netbox_virtual_machine_primary_ip.test
  id = "123"
}
