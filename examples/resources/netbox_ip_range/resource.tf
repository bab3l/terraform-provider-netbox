resource "netbox_ip_range" "test" {
  start_address = "10.0.0.1/24"
  end_address   = "10.0.0.10/24"
  status        = "active"
}
