resource "netbox_vrf" "test" {
  name = "Test VRF"
  rd   = "65000:1"
}
