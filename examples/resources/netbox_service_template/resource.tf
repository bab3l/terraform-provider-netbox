resource "netbox_service_template" "test" {
  name     = "SSH"
  protocol = "tcp"
  ports    = [22]
}
