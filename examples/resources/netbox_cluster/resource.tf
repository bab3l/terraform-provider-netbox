resource "netbox_cluster_type" "test" {
  name = "Test Cluster Type"
  slug = "test-cluster-type"
}

resource "netbox_cluster" "test" {
  name        = "Test Cluster"
  type        = netbox_cluster_type.test.name
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
