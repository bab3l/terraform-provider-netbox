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

resource "netbox_interface_template" "test" {
  name        = "GigabitEthernet1/0/1"
  device_type = netbox_device_type.test.id
  type        = "1000base-t"

  # Partial custom fields management
  # Only specified custom fields are managed, others preserved
  custom_fields = [
    {
      name  = "port_profile"
      value = "access"
    },
    {
      name  = "poe_enabled"
      value = "true"
    }
  ]
}

# Optional: seed owned custom fields during import
import {
  to = netbox_interface_template.test
  id = "123"

  identity = {
    custom_fields = [
      "port_profile:text",
      "poe_enabled:boolean",
    ]
  }
}
