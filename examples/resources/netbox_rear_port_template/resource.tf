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

  # Partial custom fields management
  # Only specified custom fields are managed, others preserved
  custom_fields = [
    {
      name  = "panel_position"
      value = "rear"
    },
    {
      name  = "port_group"
      value = "row-a"
    }
  ]
}

# Optional: seed owned custom fields during import
import {
  to = netbox_rear_port_template.test
  id = "123"

  identity = {
    custom_fields = [
      "panel_position:text",
      "port_group:text",
    ]
  }
}
