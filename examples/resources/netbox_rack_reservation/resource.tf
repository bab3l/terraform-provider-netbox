resource "netbox_site" "test" {
  name   = "Test Site"
  slug   = "test-site"
  status = "active"
}

resource "netbox_rack" "test" {
  name   = "Test Rack"
  site   = netbox_site.test.id
  status = "active"
  width  = 19
}

data "netbox_user" "admin" {
  username = "admin"
}

resource "netbox_rack_reservation" "test" {
  rack        = netbox_rack.test.id
  units       = [1, 2, 3]
  user        = data.netbox_user.admin.id
  description = "Reserved for testing"
}
