# Power Outlet Template Resource Test

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
resource "netbox_manufacturer" "test" {
  name = "Test Manufacturer for Power Outlet Template"
  slug = "test-mfg-power-outlet-tpl"
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = "Test Device Type Power Outlet Template"
  slug         = "test-dt-power-outlet-tpl"
}

# Test 1: Basic power outlet template creation
resource "netbox_power_outlet_template" "basic" {
  device_type = netbox_device_type.test.id
  name        = "Outlet 1"
}

# Test 2: Power outlet template with all optional fields
resource "netbox_power_outlet_template" "complete" {
  device_type = netbox_device_type.test.id
  name        = "Outlet 2"
  label       = "Power Outlet Template 2"
  type        = "iec-60320-c13"
  description = "Power outlet template for testing"
}
