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
  device_type = netbox_device_type.test.id
  type        = "8p8c"
  positions   = 1
}

resource "netbox_front_port_template" "test" {
  name               = "Front Port Template"
  device_type        = netbox_device_type.test.model
  type               = "8p8c"
  rear_port          = netbox_rear_port_template.test.name
  rear_port_position = 1

  # Partial custom fields management
  # Only specified custom fields are managed, others preserved
  custom_fields = [
    {
      name  = "patch_label"
      value = "PP-TPL-01"
    },
    {
      name  = "port_assignment"
      value = "distribution"
    }
  ]
}

# Optional: seed owned custom fields during import
import {
  to = netbox_front_port_template.test
  id = "123"

  identity = {
    custom_fields = [
      "patch_label:text",
      "port_assignment:text",
    ]
  }
}
