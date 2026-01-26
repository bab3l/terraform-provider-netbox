resource "netbox_site" "test" {
  name   = "Test Site"
  slug   = "test-site"
  status = "active"
}

resource "netbox_device_role" "test" {
  name  = "PDU Role"
  slug  = "pdu-role"
  color = "00ff00"
}

resource "netbox_manufacturer" "test" {
  name = "Test Manufacturer"
  slug = "test-manufacturer"
}

resource "netbox_device_type" "test" {
  model        = "PDU Model"
  slug         = "pdu-model"
  manufacturer = netbox_manufacturer.test.id
  u_height     = 0
}

resource "netbox_device" "test" {
  name        = "test-pdu-1"
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
  status      = "active"
}

resource "netbox_power_outlet" "test" {
  name        = "Outlet 1"
  device      = netbox_device.test.id
  type        = "iec-60320-c13"
  description = "PDU power outlet"

  # Partial custom fields management
  # Only specified custom fields are managed, others preserved
  custom_fields = [
    {
      name  = "max_load_watts"
      value = "1200"
    },
    {
      name  = "breaker_size_amps"
      value = "15"
    },
    {
      name  = "metered"
      value = "true"
    }
  ]

  tags = [
    "pdu-outlet",
    "monitored"
  ]
}

# Optional: seed owned custom fields during import
import {
  to = netbox_power_outlet.test
  id = "123"

  identity = {
    custom_fields = [
      "max_load_watts:integer",
      "breaker_size_amps:integer",
      "metered:boolean",
    ]
  }
}
