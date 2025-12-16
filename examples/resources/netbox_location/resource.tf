resource "netbox_site" "test" {
  name   = "Test Site"
  slug   = "test-site"
  status = "active"
}

resource "netbox_location" "test" {
  name = "Test Location"
  slug = "test-location"
  site = netbox_site.test.slug
}
