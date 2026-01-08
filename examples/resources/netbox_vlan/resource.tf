resource "netbox_site" "test" {
  name   = "Test Site"
  slug   = "test-site"
  status = "active"
}

resource "netbox_vlan_group" "test" {
  name       = "Test VLAN Group"
  slug       = "test-vlan-group"
  scope_type = "dcim.site"
  scope_id   = netbox_site.test.id
}

resource "netbox_vlan" "test" {
  vid         = 100
  name        = "Test VLAN"
  site        = netbox_site.test.slug
  group       = netbox_vlan_group.test.slug
  status      = "active"
  description = "Production server VLAN"
  comments    = "Primary VLAN for datacenter servers"

  # Partial custom fields management
  # Only specified custom fields are managed, others preserved
  custom_fields = [
    {
      name  = "vlan_purpose"
      value = "servers"
    },
    {
      name  = "subnet"
      value = "10.0.100.0/24"
    },
    {
      name  = "security_zone"
      value = "trusted"
    }
  ]

  tags = [
    "production",
    "servers"
  ]
}
