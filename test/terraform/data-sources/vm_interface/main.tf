# VM Interface Data Source Test
# This test creates a VM interface resource, then looks it up using the data source

terraform {
  required_version = ">= 1.0"
  required_providers {
    netbox = {
      source = "bab3l/netbox"
    }
  }
}

provider "netbox" {
  # Uses NETBOX_SERVER_URL and NETBOX_API_TOKEN environment variables
}

# Dependencies: Cluster Type, Cluster, and Virtual Machine
resource "netbox_cluster_type" "source_type" {
  name        = "DS VMI Test Cluster Type"
  slug        = "ds-vmi-test-cluster-type"
  description = "Cluster type for VM interface data source testing"
}

resource "netbox_cluster" "source_cluster" {
  name        = "DS VMI Test Cluster"
  type        = netbox_cluster_type.source_type.slug
  status      = "active"
  description = "Cluster for VM interface data source testing"

  depends_on = [netbox_cluster_type.source_type]
}

resource "netbox_virtual_machine" "source_vm" {
  name        = "ds-vmi-test-vm"
  cluster     = netbox_cluster.source_cluster.name
  status      = "active"
  vcpus       = 2
  memory      = 4096
  disk        = 50
  description = "VM for interface data source testing"

  depends_on = [netbox_cluster.source_cluster]
}

# Create a VM interface to look up
resource "netbox_vm_interface" "source" {
  name            = "eth0"
  virtual_machine = netbox_virtual_machine.source_vm.name
  enabled         = true
  mtu             = 1500
  mac_address     = "00:50:56:AB:CD:EF"
  description     = "Interface created for data source testing"

  depends_on = [netbox_virtual_machine.source_vm]
}

# Test 1: Look up by ID
data "netbox_vm_interface" "by_id" {
  id = netbox_vm_interface.source.id

  depends_on = [netbox_vm_interface.source]
}

# Test 2: Look up by name and virtual machine
data "netbox_vm_interface" "by_name" {
  name            = netbox_vm_interface.source.name
  virtual_machine = netbox_virtual_machine.source_vm.name

  depends_on = [netbox_vm_interface.source]
}

# Outputs for verification
output "source_id" {
  value = netbox_vm_interface.source.id
}

output "by_id_name" {
  value = data.netbox_vm_interface.by_id.name
}

output "by_id_vm" {
  value = data.netbox_vm_interface.by_id.virtual_machine
}

output "by_id_mtu" {
  value = data.netbox_vm_interface.by_id.mtu
}

output "by_id_mac" {
  value = data.netbox_vm_interface.by_id.mac_address
}

output "by_name_enabled" {
  value = data.netbox_vm_interface.by_name.enabled
}

output "by_name_description" {
  value = data.netbox_vm_interface.by_name.description
}

# Verify all lookups return the same ID
output "ids_match" {
  value = data.netbox_vm_interface.by_id.id == data.netbox_vm_interface.by_name.id
}
