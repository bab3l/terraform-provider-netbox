resource "netbox_site" "test" {
  name   = "Test Site"
  slug   = "test-site"
  status = "active"
}

resource "netbox_device_role" "test" {
  name  = "Test Role"
  slug  = "test-role"
  color = "ff0000"
}

resource "netbox_device_type" "test" {
  model        = "Test Model"
  slug         = "test-model"
  manufacturer = netbox_manufacturer.test.id
  u_height     = 1
}

resource "netbox_manufacturer" "test" {
  name = "Test Manufacturer"
  slug = "test-manufacturer"
}

resource "netbox_device" "test" {
  name        = "test-device-1"
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
  status      = "active"
}

resource "netbox_console_port" "test" {
  name        = "Console Port 1"
  device      = netbox_device.test.id
  type        = "rj-45"
  description = "Management console port"

  # Partial custom fields management
  # Only specified custom fields are managed, others preserved
  custom_fields = [
    {
      name  = "console_server"
      value = "console-srv-01"
    },
    {
      name  = "port_number"
      value = "16"
    },
    {
      name  = "baud_rate"
      value = "9600"
    }
  ]

  tags = [
    "console-access",
    "management"
  ]
}

# Optional: seed owned custom fields during import
import {
  to = netbox_console_port.test
  id = "123"

  identity = {
    custom_fields = [
      "console_server:text",
      "port_number:integer",
      "baud_rate:integer",
    ]
  }
}
