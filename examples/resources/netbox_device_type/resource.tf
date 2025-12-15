resource "netbox_manufacturer" "test" {
  name = "Test Manufacturer"
  slug = "test-manufacturer"
}

resource "netbox_device_type" "test" {
  model         = "Test Model"
  slug          = "test-model"
  manufacturer  = netbox_manufacturer.test.id
  u_height      = 1
  is_full_depth = true
}
