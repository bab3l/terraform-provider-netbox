resource "netbox_l2vpn" "test" {
  name = "Test L2VPN"
  slug = "test-l2vpn"
  type = "vxlan"
}
