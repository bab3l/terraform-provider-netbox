# Power Outlet Data Source Test

terraform {
  required_version = ">= 1.0"
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
  name   = "Test Site for Power Outlet DS"
  slug   = "test-site-power-outlet-ds"
  status = "active"
}

resource "netbox_manufacturer" "test" {
  name = "Test Manufacturer for Power Outlet DS"
  slug = "test-manufacturer-power-outlet-ds"
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = "Test Model for Power Outlet DS"
  slug         = "test-model-power-outlet-ds"
}

resource "netbox_device_role" "test" {
  name  = "Test Role for Power Outlet DS"
  slug  = "test-role-power-outlet-ds"
  color = "aabbcc"
}

resource "netbox_device" "test" {
  name        = "Test Device for Power Outlet DS"
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
  status      = "active"
}

resource "netbox_power_outlet" "test" {
  name        = "Outlet1-DS"
  device      = netbox_device.test.id
  description = "Test power outlet for data source"
}

# Test: Lookup power outlet by ID
data "netbox_power_outlet" "by_id" {
  id = netbox_power_outlet.test.id
}
