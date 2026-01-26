resource "netbox_cluster_type" "test" {
  name = "Test Cluster Type"
  slug = "test-cluster-type"
}

resource "netbox_cluster" "test" {
  name        = "Test Cluster"
  type        = netbox_cluster_type.test.id
  description = "Production virtualization cluster"
  comments    = "Primary VMware cluster"

  # Partial custom fields management
  # Only specified custom fields are managed, others preserved
  custom_fields = [
    {
      name  = "hypervisor_version"
      value = "vSphere 8.0"
    },
    {
      name  = "cluster_nodes"
      value = "4"
    },
    {
      name  = "management_url"
      value = "https://vcenter.example.com"
    },
    {
      name  = "ha_enabled"
      value = "true"
    }
  ]

  tags = [
    "vmware-cluster",
    "production"
  ]
}

# Optional: seed owned custom fields during import
import {
  to = netbox_cluster.test
  id = "123"

  identity = {
    custom_fields = [
      "hypervisor_version:text",
      "cluster_nodes:integer",
      "management_url:text",
      "ha_enabled:boolean",
    ]
  }
}
