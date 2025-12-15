resource "netbox_manufacturer" "test" {
  name = "Test Manufacturer"
  slug = "test-manufacturer"
}

resource "netbox_rack_type" "test" {
  model        = "Test Rack Type"
  slug         = "test-rack-type"
  manufacturer = netbox_manufacturer.test.id
  width        = 19
  u_height     = 42
}
