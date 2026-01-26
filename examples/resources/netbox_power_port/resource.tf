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

resource "netbox_power_port" "test" {
  name        = "PSU1"
  device      = netbox_device.test.id
  type        = "iec-60320-c14"
  description = "Primary power supply inlet"

  # Partial custom fields management
  # Only specified custom fields are managed, others in NetBox preserved
  custom_fields = [
    {
      name  = "pdu_connection"
      value = "PDU-A-001-PORT-12"
    },
    {
      name  = "rated_voltage"
      value = "120"
    },
    {
      name  = "max_amperage"
      value = "15"
    }
  ]

  tags = [
    "power-input",
    "primary"
  ]
}

# Optional: seed owned custom fields during import
import {
  to = netbox_power_port.test
  id = "123"

  identity = {
    custom_fields = [
      "pdu_connection:text",
      "rated_voltage:integer",
      "max_amperage:integer",
    ]
  }
}
