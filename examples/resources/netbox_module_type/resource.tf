resource "netbox_manufacturer" "test" {
  name = "Test Manufacturer"
  slug = "test-manufacturer"
}

resource "netbox_module_type" "test" {
  model        = "Test Module Type"
  manufacturer = netbox_manufacturer.test.id
}

import {
  to = netbox_module_type.test
  id = "123"
}
