resource "netbox_site" "test" {
  name   = "Test Site"
  slug   = "test-site"
  status = "active"
}

resource "netbox_rack_role" "test" {
  name  = "Test Rack Role"
  slug  = "test-rack-role"
  color = "ff0000"
}

resource "netbox_rack" "test" {
  name     = "test-rack-1"
  site     = netbox_site.test.id
  status   = "active"
  role     = netbox_rack_role.test.id
  u_height = 42
  width    = 19
}
