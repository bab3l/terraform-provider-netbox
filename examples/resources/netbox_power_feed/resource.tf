resource "netbox_site" "test" {
  name   = "Test Site"
  slug   = "test-site"
  status = "active"
}

resource "netbox_power_panel" "test" {
  name = "Test Power Panel"
  site = netbox_site.test.slug
}

resource "netbox_rack" "test" {
  name   = "Test Rack"
  site   = netbox_site.test.slug
  status = "active"
  width  = 19
}

resource "netbox_power_feed" "test" {
  name        = "Test Power Feed"
  power_panel = netbox_power_panel.test.name
  rack        = netbox_rack.test.name
  status      = "active"
  type        = "primary"
  supply      = "ac"
  phase       = "single-phase"
  voltage     = 230
  amperage    = 32
  description = "Primary power feed to rack"
  comments    = "Main electrical feed from panel"

  # Partial custom fields management
  # Only specified custom fields are managed, others in NetBox preserved
  custom_fields = [
    {
      name  = "circuit_id"
      value = "FEED-001-A"
    },
    {
      name  = "upstream_breaker"
      value = "32A-BREAKER-12"
    },
    {
      name  = "redundant_feed"
      value = "FEED-001-B"
    }
  ]

  tags = [
    "power-feed",
    "primary"
  ]
}
