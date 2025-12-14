# VM Interface Resource Test

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
resource "netbox_cluster_type" "vmware" {
  name        = "VMware vSphere (VMI Test)"
  slug        = "vmware-vsphere-vmi-test"
  description = "VMware cluster type for VM interface testing"
}

resource "netbox_cluster" "test" {
  name        = "Test Cluster (VMI Test)"
  type        = netbox_cluster_type.vmware.slug
  status      = "active"
  description = "Test cluster for VM interface testing"

  depends_on = [netbox_cluster_type.vmware]
}

resource "netbox_virtual_machine" "webserver" {
  name        = "webserver-vmi-test"
  cluster     = netbox_cluster.test.name
  status      = "active"
  vcpus       = 2
  memory      = 4096
  disk        = 40
  description = "Web server for interface testing"

  depends_on = [netbox_cluster.test]
}

resource "netbox_virtual_machine" "database" {
  name        = "database-vmi-test"
  cluster     = netbox_cluster.test.name
  status      = "active"
  vcpus       = 4
  memory      = 8192
  disk        = 100
  description = "Database server for interface testing"

  depends_on = [netbox_cluster.test]
}

# Test 1: Basic VM interface creation
resource "netbox_vm_interface" "basic" {
  name            = "eth0"
  virtual_machine = netbox_virtual_machine.webserver.name

  depends_on = [netbox_virtual_machine.webserver]
}

# Test 2: VM interface with all optional fields
resource "netbox_vm_interface" "complete" {
  name            = "eth1"
  virtual_machine = netbox_virtual_machine.webserver.name
  enabled         = true
  mtu             = 1500
  mac_address     = "00:50:56:A1:B2:C3"
  description     = "Primary network interface"

  depends_on = [netbox_virtual_machine.webserver]
}

# Test 3: VM interface with custom MTU (jumbo frames)
resource "netbox_vm_interface" "jumbo" {
  name            = "eth2"
  virtual_machine = netbox_virtual_machine.webserver.name
  enabled         = true
  mtu             = 9000
  description     = "Storage network with jumbo frames"

  depends_on = [netbox_virtual_machine.webserver]
}

# Test 4: Disabled VM interface
resource "netbox_vm_interface" "disabled" {
  name            = "eth3"
  virtual_machine = netbox_virtual_machine.webserver.name
  enabled         = false
  description     = "Disabled backup interface"

  depends_on = [netbox_virtual_machine.webserver]
}

# Test 5: Database server primary interface
resource "netbox_vm_interface" "db_primary" {
  name            = "eth0"
  virtual_machine = netbox_virtual_machine.database.name
  enabled         = true
  mtu             = 1500
  mac_address     = "00:50:56:DB:00:01"
  description     = "Database primary network"

  depends_on = [netbox_virtual_machine.database]
}

# Test 6: Database server replication interface
resource "netbox_vm_interface" "db_replication" {
  name            = "eth1"
  virtual_machine = netbox_virtual_machine.database.name
  enabled         = true
  mtu             = 9000
  mac_address     = "00:50:56:DB:00:02"
  description     = "Database replication network"

  depends_on = [netbox_virtual_machine.database]
}

# Outputs for verification
output "basic_id" {
  value = netbox_vm_interface.basic.id
}

output "complete_name" {
  value = netbox_vm_interface.complete.name
}

output "complete_mtu" {
  value = netbox_vm_interface.complete.mtu
}

output "complete_mac" {
  value = netbox_vm_interface.complete.mac_address
}

output "jumbo_mtu" {
  value = netbox_vm_interface.jumbo.mtu
}

output "disabled_enabled" {
  value = netbox_vm_interface.disabled.enabled
}
