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
  vid    = 100
  name   = "Test VLAN"
  site   = netbox_site.test.slug
  group  = netbox_vlan_group.test.slug
  status = "active"
}
