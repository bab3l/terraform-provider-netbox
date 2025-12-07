# Power Panel Resource Test

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
  name   = "Test Site for Power Panel"
  slug   = "test-site-power-panel"
  status = "active"
}

resource "netbox_location" "test" {
  name   = "Test Location for Power Panel"
  slug   = "test-location-power-panel"
  site   = netbox_site.test.id
  status = "active"
}

# Test 1: Basic power panel creation
resource "netbox_power_panel" "basic" {
  name = "Test Power Panel Basic"
  site = netbox_site.test.id
}

# Test 2: Power panel with all optional fields
resource "netbox_power_panel" "complete" {
  name        = "Test Power Panel Complete"
  site        = netbox_site.test.id
  location    = netbox_location.test.id
  description = "A power panel for testing"
  comments    = "This power panel was created for integration testing."
}
