resource "netbox_site" "test" {
  name   = "Test Site"
  slug   = "test-site"
  status = "active"
}

resource "netbox_vlan" "test" {
  vid    = 100
  name   = "Test VLAN"
  site   = netbox_site.test.name
  status = "active"
}

resource "netbox_l2vpn" "test" {
  name = "Test L2VPN"
  slug = "test-l2vpn"
  type = "vxlan"
}

resource "netbox_l2vpn_termination" "test" {
  l2vpn                = netbox_l2vpn.test.id
  assigned_object_type = "ipam.vlan"
  assigned_object_id   = netbox_vlan.test.id
}
