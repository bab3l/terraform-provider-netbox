resource "netbox_manufacturer" "test" {
  name = "Test Manufacturer"
  slug = "test-manufacturer"
}

resource "netbox_device_type" "test" {
  model        = "Test Model"
  slug         = "test-model"
  manufacturer = netbox_manufacturer.test.id
  u_height     = 1
}

resource "netbox_rear_port_template" "test" {
  name        = "Rear Port Template"
  device_type = netbox_device_type.test.model
  type        = "8p8c"
  positions   = 1
}
