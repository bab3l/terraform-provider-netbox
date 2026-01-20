resource "netbox_manufacturer" "test" {
  name = "Test Manufacturer"
  slug = "test-manufacturer"
}

resource "netbox_device_type" "test" {
  model         = "Test Model"
  slug          = "test-model"
  manufacturer  = netbox_manufacturer.test.id
  u_height      = 1
  is_full_depth = true
  description   = "1U network switch"
  comments      = "Standard datacenter ToR switch model"

  # Partial custom fields management
  # Only specified custom fields are managed, others in NetBox preserved
  custom_fields = [
    {
      name  = "power_supply_type"
      value = "redundant-ac"
    },
    {
      name  = "max_power_draw_watts"
      value = "150"
    },
    {
      name  = "port_count"
      value = "48"
    },
    {
      name  = "end_of_life_date"
      value = "2028-12-31"
    }
  ]

  tags = [
    "network-switch",
    "1u"
  ]
}

# Optional: seed owned custom fields during import
import {
  to = netbox_device_type.test
  id = "123"

  identity = {
    custom_fields = [
      "power_supply_type:text",
      "max_power_draw_watts:integer",
      "port_count:integer",
      "end_of_life_date:date",
    ]
  }
}
