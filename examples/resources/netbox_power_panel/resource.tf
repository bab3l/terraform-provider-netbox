resource "netbox_site" "test" {
  name   = "Test Site"
  slug   = "test-site"
  status = "active"
}

resource "netbox_power_panel" "test" {
  name = "Test Power Panel"
  site = netbox_site.test.slug

  # Partial custom fields management
  # Only specified custom fields are managed, others in NetBox preserved
  custom_fields = [
    {
      name  = "panel_voltage"
      value = "208"
    },
    {
      name  = "panel_amperage"
      value = "200"
    },
    {
      name  = "panel_type"
      value = "3-phase"
    },
    {
      name  = "breaker_count"
      value = "42"
    }
  ]

  tags = [
    "electrical-panel",
    "datacenter-power"
  ]
}

# Optional: seed owned custom fields during import
import {
  to = netbox_power_panel.test
  id = "123"

  identity = {
    custom_fields = [
      "panel_voltage:integer",
      "panel_amperage:integer",
      "panel_type:text",
      "breaker_count:integer",
    ]
  }
}
