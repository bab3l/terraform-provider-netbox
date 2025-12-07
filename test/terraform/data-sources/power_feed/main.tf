# Power Feed Data Source Test

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
  name   = "Test Site for Power Feed DS"
  slug   = "test-site-power-feed-ds"
  status = "active"
}

resource "netbox_power_panel" "test" {
  name = "Test Power Panel for Power Feed DS"
  site = netbox_site.test.id
}

resource "netbox_power_feed" "test" {
  name        = "Feed-A-DS"
  power_panel = netbox_power_panel.test.id
  status      = "active"
  description = "Test power feed for data source"
}

# Test: Lookup power feed by ID
data "netbox_power_feed" "by_id" {
  id = netbox_power_feed.test.id
}
