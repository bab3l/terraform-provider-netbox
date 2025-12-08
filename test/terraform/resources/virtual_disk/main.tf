# Virtual Disk Resource Test

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
resource "netbox_cluster_type" "vd_test" {
  name = "VD Test Cluster Type"
  slug = "vd-test-cluster-type"
}

resource "netbox_cluster" "vd_test" {
  name = "VD Test Cluster"
  type = netbox_cluster_type.vd_test.id
}

resource "netbox_virtual_machine" "vd_test" {
  name    = "VD Test VM"
  cluster = netbox_cluster.vd_test.id
}

# Test 1: Basic virtual disk creation
resource "netbox_virtual_disk" "basic" {
  virtual_machine = netbox_virtual_machine.vd_test.id
  name            = "disk0"
  size            = "100"
}

# Test 2: Virtual disk with description
resource "netbox_virtual_disk" "complete" {
  virtual_machine = netbox_virtual_machine.vd_test.id
  name            = "disk1"
  size            = "500"
  description     = "Primary data disk"
}
