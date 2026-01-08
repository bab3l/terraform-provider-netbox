resource "netbox_site" "test" {
  name   = "Test Site"
  slug   = "test-site"
  status = "active"
}

resource "netbox_vlan_group" "test" {
  name        = "Test VLAN Group"
  slug        = "test-vlan-group"
  scope_type  = "dcim.site"
  scope_id    = netbox_site.test.id
  description = "Primary VLAN group for datacenter"

  # Partial custom fields management
  # Only specified custom fields are managed, others preserved
  custom_fields = [
    {
      name  = "vlan_range"
      value = "100-199"
    },
    {
      name  = "group_purpose"
      value = "production"
    },
    {
      name  = "max_vlans"
      value = "100"
    }
  ]

  tags = [
    "vlan-group",
    "datacenter"
  ]
}
