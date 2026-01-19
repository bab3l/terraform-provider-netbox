resource "netbox_manufacturer" "test" {
  name = "Test Manufacturer"
  slug = "test-manufacturer"
}

resource "netbox_rack_type" "test" {
  model        = "Test Rack Type"
  slug         = "test-rack-type"
  manufacturer = netbox_manufacturer.test.name
  width        = 19
  u_height     = 42
}

# Optional: seed owned custom fields during import
import {
  to = netbox_rack_type.test
  id = "123"

  identity = {
    custom_fields = [
      "mounting_type:text",
      "shipping_class:text",
    ]
  }
}
