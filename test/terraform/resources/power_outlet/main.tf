# Power Outlet Resource Test

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
  name   = "Test Site for Power Outlet"
  slug   = "test-site-power-outlet"
  status = "active"
}

resource "netbox_manufacturer" "test" {
  name = "Test Manufacturer for Power Outlet"
  slug = "test-mfg-power-outlet"
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = "Test Device Type Power Outlet"
  slug         = "test-dt-power-outlet"
}

resource "netbox_device_role" "test" {
  name  = "Test Role for Power Outlet"
  slug  = "test-role-power-outlet"
  color = "aa1409"
}

resource "netbox_device" "test" {
  name        = "test-device-power-outlet"
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
}

# Test 1: Basic power outlet creation
resource "netbox_power_outlet" "basic" {
  device = netbox_device.test.id
  name   = "Outlet 1"
}

# Test 2: Power outlet with all optional fields
resource "netbox_power_outlet" "complete" {
  device      = netbox_device.test.id
  name        = "Outlet 2"
  label       = "Power Outlet 2"
  type        = "iec-60320-c13"
  description = "Power outlet for testing"
}
