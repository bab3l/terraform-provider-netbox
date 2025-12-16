resource "netbox_ip_address" "test_v4" {
  address  = "10.0.0.1/24"
  status   = "active"
  dns_name = "test.example.com"
}

resource "netbox_ip_address" "test_v6" {
  address = "2001:db8::1/64"
  status  = "active"
}

resource "netbox_vrf" "test" {
  name = "Test VRF"
  rd   = "65000:1"
}

resource "netbox_ip_address" "test_vrf" {
  address = "192.168.1.1/24"
  vrf     = netbox_vrf.test.name
  status  = "active"
}
