# Virtual Machine Resource Test

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
resource "netbox_cluster_type" "vmware" {
  name        = "VMware vSphere (VM Test)"
  slug        = "vmware-vsphere-vm-test"
  description = "VMware cluster type for VM testing"
}

resource "netbox_cluster" "production" {
  name        = "Production Cluster (VM Test)"
  type        = netbox_cluster_type.vmware.slug
  status      = "active"
  description = "Production cluster for VM testing"

  depends_on = [netbox_cluster_type.vmware]
}

resource "netbox_cluster" "development" {
  name        = "Development Cluster (VM Test)"
  type        = netbox_cluster_type.vmware.slug
  status      = "active"
  description = "Development cluster for VM testing"

  depends_on = [netbox_cluster_type.vmware]
}

# Test 1: Basic virtual machine creation
resource "netbox_virtual_machine" "basic" {
  name    = "test-vm-basic"
  cluster = netbox_cluster.production.name

  depends_on = [netbox_cluster.production]
}

# Test 2: Virtual machine with all optional fields
resource "netbox_virtual_machine" "complete" {
  name        = "test-vm-complete"
  cluster     = netbox_cluster.production.name
  status      = "active"
  vcpus       = 4
  memory      = 8192
  disk        = 100
  description = "A complete test VM with all fields"
  comments    = "This VM was created for integration testing purposes."

  depends_on = [netbox_cluster.production]
}

# Test 3: Staged VM
resource "netbox_virtual_machine" "staged" {
  name        = "test-vm-staged"
  cluster     = netbox_cluster.production.name
  status      = "staged"
  vcpus       = 2
  memory      = 4096
  disk        = 50
  description = "A VM in staging status"

  depends_on = [netbox_cluster.production]
}

# Test 4: Offline VM
resource "netbox_virtual_machine" "offline" {
  name        = "test-vm-offline"
  cluster     = netbox_cluster.production.name
  status      = "offline"
  vcpus       = 1
  memory      = 2048
  disk        = 20
  description = "An offline VM"

  depends_on = [netbox_cluster.production]
}

# Test 5: Web server VM
resource "netbox_virtual_machine" "webserver" {
  name        = "web-server-01"
  cluster     = netbox_cluster.production.name
  status      = "active"
  vcpus       = 2
  memory      = 4096
  disk        = 40
  description = "Production web server"
  comments    = "Nginx web server handling production traffic"

  depends_on = [netbox_cluster.production]
}

# Test 6: Database VM
resource "netbox_virtual_machine" "database" {
  name        = "db-server-01"
  cluster     = netbox_cluster.production.name
  status      = "active"
  vcpus       = 8
  memory      = 32768
  disk        = 500
  description = "Production database server"
  comments    = "PostgreSQL database server"

  depends_on = [netbox_cluster.production]
}

# Test 7: Development VM
resource "netbox_virtual_machine" "dev" {
  name        = "dev-server-01"
  cluster     = netbox_cluster.development.name
  status      = "active"
  vcpus       = 4
  memory      = 8192
  disk        = 100
  description = "Development server"
  comments    = "Shared development environment"

  depends_on = [netbox_cluster.development]
}

# Outputs for verification
output "basic_id" {
  value = netbox_virtual_machine.basic.id
}

output "complete_name" {
  value = netbox_virtual_machine.complete.name
}

output "complete_vcpus" {
  value = netbox_virtual_machine.complete.vcpus
}

output "complete_memory" {
  value = netbox_virtual_machine.complete.memory
}

output "database_disk" {
  value = netbox_virtual_machine.database.disk
}

output "webserver_status" {
  value = netbox_virtual_machine.webserver.status
}
