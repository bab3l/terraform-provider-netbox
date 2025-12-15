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
