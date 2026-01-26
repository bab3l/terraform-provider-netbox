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

resource "netbox_module_bay_template" "test" {
  name        = "Module Bay Template"
  device_type = netbox_device_type.test.id
}

import {
  to = netbox_module_bay_template.test
  id = "123"
}
