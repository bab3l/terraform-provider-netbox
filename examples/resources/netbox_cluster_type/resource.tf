resource "netbox_cluster_type" "test" {
  name        = "Test Cluster Type"
  slug        = "test-cluster-type"
  description = "VMware vSphere cluster type"

  # Partial custom fields management
  # Only specified custom fields are managed, others in NetBox preserved
  custom_fields = [
    {
      name  = "platform_vendor"
      value = "VMware"
    },
    {
      name  = "platform_type"
      value = "vSphere"
    },
    {
      name  = "default_overcommit_ratio"
      value = "2.0"
    }
  ]

  tags = [
    "cluster-type",
    "vmware"
  ]
}
