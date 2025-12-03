# Virtual Machine Data Source Test
# This test creates a virtual machine resource, then looks it up using the data source

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

# Dependencies: Cluster Type and Cluster
resource "netbox_cluster_type" "source_type" {
  name        = "DS VM Test Cluster Type"
  slug        = "ds-vm-test-cluster-type"
  description = "Cluster type for VM data source testing"
}

resource "netbox_cluster" "source_cluster" {
  name        = "DS VM Test Cluster"
  type        = netbox_cluster_type.source_type.slug
  status      = "active"
  description = "Cluster for VM data source testing"

  depends_on = [netbox_cluster_type.source_type]
}

# Create a virtual machine to look up
resource "netbox_virtual_machine" "source" {
  name        = "ds-test-vm"
  cluster     = netbox_cluster.source_cluster.name
  status      = "active"
  vcpus       = 4
  memory      = 8192
  disk        = 100
  description = "VM created for data source testing"
  comments    = "Test VM for DS tests"

  depends_on = [netbox_cluster.source_cluster]
}

# Test 1: Look up by ID
data "netbox_virtual_machine" "by_id" {
  id = netbox_virtual_machine.source.id

  depends_on = [netbox_virtual_machine.source]
}

# Test 2: Look up by name
data "netbox_virtual_machine" "by_name" {
  name = netbox_virtual_machine.source.name

  depends_on = [netbox_virtual_machine.source]
}

# Outputs for verification
output "source_id" {
  value = netbox_virtual_machine.source.id
}

output "by_id_name" {
  value = data.netbox_virtual_machine.by_id.name
}

output "by_id_cluster" {
  value = data.netbox_virtual_machine.by_id.cluster
}

output "by_id_vcpus" {
  value = data.netbox_virtual_machine.by_id.vcpus
}

output "by_id_memory" {
  value = data.netbox_virtual_machine.by_id.memory
}

output "by_name_status" {
  value = data.netbox_virtual_machine.by_name.status
}

output "by_name_disk" {
  value = data.netbox_virtual_machine.by_name.disk
}

output "by_name_description" {
  value = data.netbox_virtual_machine.by_name.description
}

# Verify all lookups return the same ID
output "ids_match" {
  value = data.netbox_virtual_machine.by_id.id == data.netbox_virtual_machine.by_name.id
}
