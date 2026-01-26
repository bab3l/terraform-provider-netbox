resource "netbox_site" "test" {
  name   = "Test Site"
  slug   = "test-site"
  status = "active"
}

resource "netbox_rack_role" "test" {
  name  = "Test Rack Role"
  slug  = "test-rack-role"
  color = "ff0000"
}

resource "netbox_rack" "test" {
  name        = "test-rack-1"
  site        = netbox_site.test.id
  status      = "active"
  role        = netbox_rack_role.test.id
  facility_id = "FAC-01"
  u_height    = 42
  width       = 19
  description = "Primary datacenter rack"
  comments    = "Main equipment rack in server room"

  # Partial custom fields management
  # Only specified custom fields are managed, others preserved
  custom_fields = [
    {
      name  = "power_circuit_a"
      value = "PDU-A-001"
    },
    {
      name  = "power_circuit_b"
      value = "PDU-B-001"
    },
    {
      name  = "max_power_draw_watts"
      value = "5000"
    },
    {
      name  = "cooling_zone"
      value = "hot-aisle-1"
    }
  ]

  tags = [
    "production",
    "server-rack"
  ]
}

# Optional: seed owned custom fields during import
import {
  to = netbox_rack.test
  id = "123"

  identity = {
    custom_fields = [
      "power_circuit_a:text",
      "power_circuit_b:text",
      "max_power_draw_watts:integer",
      "cooling_zone:text",
    ]
  }
}
