resource "netbox_site" "test" {
  name   = "Test Site"
  slug   = "test-site"
  status = "active"
}

resource "netbox_power_panel" "test" {
  name = "Test Power Panel"
  site = netbox_site.test.slug
}
