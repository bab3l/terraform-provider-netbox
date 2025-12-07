# Power Panel Data Source Test

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
  name   = "Test Site for Power Panel DS"
  slug   = "test-site-power-panel-ds"
  status = "active"
}

resource "netbox_power_panel" "test" {
  name        = "Test Power Panel DS"
  site        = netbox_site.test.id
  description = "Test power panel for data source"
}

# Test: Lookup power panel by ID
data "netbox_power_panel" "by_id" {
  id = netbox_power_panel.test.id
}

# Test: Lookup power panel by name
data "netbox_power_panel" "by_name" {
  name = netbox_power_panel.test.name

  depends_on = [netbox_power_panel.test]
}
