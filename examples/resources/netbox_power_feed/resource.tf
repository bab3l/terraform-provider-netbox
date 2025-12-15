resource "netbox_site" "test" {
  name   = "Test Site"
  slug   = "test-site"
  status = "active"
}

resource "netbox_power_panel" "test" {
  name = "Test Power Panel"
  site = netbox_site.test.id
}

resource "netbox_rack" "test" {
  name   = "Test Rack"
  site   = netbox_site.test.id
  status = "active"
  width  = 19
}

resource "netbox_power_feed" "test" {
  name        = "Test Power Feed"
  power_panel = netbox_power_panel.test.id
  rack        = netbox_rack.test.id
  status      = "active"
  type        = "primary"
  supply      = "ac"
  phase       = "single-phase"
  voltage     = 230
  amperage    = 32
}
