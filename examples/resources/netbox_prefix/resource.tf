resource "netbox_site" "test" {
  name   = "Test Site"
  slug   = "test-site"
  status = "active"
}

resource "netbox_vlan" "test" {
  vid    = 100
  name   = "Test VLAN"
  site   = netbox_site.test.id
  status = "active"
}

resource "netbox_prefix" "test" {
  prefix = "10.0.0.0/24"
  site   = netbox_site.test.id
  vlan   = netbox_vlan.test.id
  status = "active"
}
