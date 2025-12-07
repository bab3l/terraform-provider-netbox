# Power Feed Resource Test

terraform {
  required_providers {
    netbox = {
      source = "bab3l/netbox"
    }
  }
}

provider "netbox" {
  # Uses NETBOX_SERVER_URL and NETBOX_API_TOKEN environment variables
}

# Dependencies
resource "netbox_site" "test" {
  name   = "Test Site for Power Feed"
  slug   = "test-site-power-feed"
  status = "active"
}

resource "netbox_power_panel" "test" {
  name = "Test Power Panel"
  site = netbox_site.test.id
}

resource "netbox_location" "test" {
  name   = "Test Location for Power Feed"
  slug   = "test-location-power-feed"
  site   = netbox_site.test.id
  status = "active"
}

resource "netbox_rack" "test" {
  name     = "Test Rack for Power Feed"
  site     = netbox_site.test.id
  location = netbox_location.test.id
}

# Test 1: Basic power feed creation
resource "netbox_power_feed" "basic" {
  name        = "Power Feed A"
  power_panel = netbox_power_panel.test.id
}

# Test 2: Power feed with all optional fields
resource "netbox_power_feed" "complete" {
  name        = "Power Feed B"
  power_panel = netbox_power_panel.test.id
  rack        = netbox_rack.test.id
  status      = "active"
  type        = "primary"
  supply      = "ac"
  phase       = "single-phase"
  voltage     = 220
  amperage    = 32
  description = "Power feed for testing"
}
