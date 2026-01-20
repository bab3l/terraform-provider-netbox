resource "netbox_site" "test" {
  name   = "Test Site"
  slug   = "test-site"
  status = "active"
}

resource "netbox_device_role" "test" {
  name  = "Console Server Role"
  slug  = "console-server-role"
  color = "0000ff"
}

resource "netbox_manufacturer" "test" {
  name = "Test Manufacturer"
  slug = "test-manufacturer"
}

resource "netbox_device_type" "test" {
  model        = "Console Server Model"
  slug         = "console-server-model"
  manufacturer = netbox_manufacturer.test.name
  u_height     = 1
}

resource "netbox_device" "test" {
  name        = "test-console-server-1"
  device_type = netbox_device_type.test.model
  role        = netbox_device_role.test.name
  site        = netbox_site.test.name
  status      = "active"
}

resource "netbox_console_server_port" "test" {
  name        = "Port 1"
  device      = netbox_device.test.name
  type        = "rj-45"
  description = "Console server port for device management"

  # Partial custom fields management
  # Only specified custom fields are managed, others in NetBox preserved
  custom_fields = [
    {
      name  = "connected_device_hostname"
      value = "switch-01.example.com"
    },
    {
      name  = "access_method"
      value = "ssh"
    },
    {
      name  = "baud_rate"
      value = "9600"
    }
  ]

  tags = [
    "console-server",
    "management-access"
  ]
}

# Optional: seed owned custom fields during import
import {
  to = netbox_console_server_port.test
  id = "123"

  identity = {
    custom_fields = [
      "connected_device_hostname:text",
      "access_method:text",
      "baud_rate:integer",
    ]
  }
}
