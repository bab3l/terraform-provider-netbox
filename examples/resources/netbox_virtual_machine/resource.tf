resource "netbox_cluster_type" "test" {
  name = "Test Cluster Type"
  slug = "test-cluster-type"
}

resource "netbox_cluster" "test" {
  name = "Test Cluster"
  type = netbox_cluster_type.test.id
}

resource "netbox_config_template" "test" {
  name          = "test-config-template"
  template_code = "hostname {{ device.name }}"
}

resource "netbox_virtual_machine" "test" {
  name        = "test-vm-1"
  cluster     = netbox_cluster.test.id
  vcpus       = 2
  memory      = 4096
  disk        = 50
  status      = "active"
  description = "Test virtual machine"
  comments    = "Production web server VM"
  serial      = "VM-serial-12345"

  config_template = netbox_config_template.test.id

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

# Optional: seed owned custom fields during import
import {
  to = netbox_virtual_machine.test
  id = "123"

  identity = {
    custom_fields = [
      "operating_system:text",
      "environment:text",
      "backup_policy:text",
      "vm_owner:text",
    ]
  }
}
