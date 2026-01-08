resource "netbox_cluster_type" "test" {
  name = "Test Cluster Type"
  slug = "test-cluster-type"
}

resource "netbox_cluster" "test" {
  name = "Test Cluster"
  type = netbox_cluster_type.test.id
}

resource "netbox_virtual_machine" "test" {
  name        = "test-vm-1"
  cluster     = netbox_cluster.test.name
  vcpus       = 2
  memory      = 4096
  disk        = 50
  status      = "active"
  description = "Test virtual machine"
  comments    = "Production web server VM"

  # Partial custom fields management
  # Only specified custom fields are managed, others in NetBox preserved
  custom_fields = [
    {
      name  = "operating_system"
      value = "Ubuntu 22.04 LTS"
    },
    {
      name  = "environment"
      value = "production"
    },
    {
      name  = "backup_policy"
      value = "daily"
    },
    {
      name  = "vm_owner"
      value = "web-team"
    }
  ]

  tags = [
    "virtual-machine",
    "production",
    "web-server"
  ]
}
