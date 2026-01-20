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

resource "netbox_console_port_template" "test" {
  name        = "Console Port Template"
  device_type = netbox_device_type.test.model
  type        = "rj-45"

  # Partial custom fields management
  # Only specified custom fields are managed, others preserved
  custom_fields = [
    {
      name  = "console_role"
      value = "primary"
    },
    {
      name  = "baud_rate"
      value = "9600"
    }
  ]
}

# Optional: seed owned custom fields during import
import {
  to = netbox_console_port_template.test
  id = "123"

  identity = {
    custom_fields = [
      "console_role:text",
      "baud_rate:integer",
    ]
  }
}
