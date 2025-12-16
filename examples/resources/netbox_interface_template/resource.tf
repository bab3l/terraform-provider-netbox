resource "netbox_manufacturer" "test" {
  name = "Test Manufacturer"
  slug = "test-manufacturer"
}

resource "netbox_device_type" "test" {
  model           = "Test Model"
  slug            = "test-model"
  manufacturer_id = netbox_manufacturer.test.id
  u_height        = 1
}

resource "netbox_interface_template" "test" {
  name        = "GigabitEthernet1/0/1"
  device_type = netbox_device_type.test.model
  type        = "1000base-t"
}
