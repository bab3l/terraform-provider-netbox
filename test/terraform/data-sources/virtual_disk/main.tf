# Virtual Disk Data Source Test

terraform {
  required_providers {
    netbox = {
      source = "bab3l/netbox"
    }
  }
}

provider "netbox" {
  # Uses NETBOX_SERVER_URL and NETBOX_API_TOKEN environment variables
}

# Create required dependencies
resource "netbox_cluster_type" "ds_test" {
  name = "VD DS Test Cluster Type"
  slug = "vd-ds-test-cluster-type"
}

resource "netbox_cluster" "ds_test" {
  name = "VD DS Test Cluster"
  type = netbox_cluster_type.ds_test.id
}

resource "netbox_virtual_machine" "ds_test" {
  name    = "VD DS Test VM"
  cluster = netbox_cluster.ds_test.id
}

# Create a virtual disk to look up
resource "netbox_virtual_disk" "test" {
  virtual_machine = netbox_virtual_machine.ds_test.id
  name            = "ds-test-disk"
  size            = "250"
  description     = "Test disk for data source"
}

# Look up by ID
data "netbox_virtual_disk" "by_id" {
  id = netbox_virtual_disk.test.id
}

# Look up by name (with virtual_machine)
data "netbox_virtual_disk" "by_name" {
  name            = netbox_virtual_disk.test.name
  virtual_machine = netbox_virtual_machine.ds_test.id
}
