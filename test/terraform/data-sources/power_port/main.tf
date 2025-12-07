# Power Port Data Source Test

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
  name   = "Test Site for Power Port DS"
  slug   = "test-site-power-port-ds"
  status = "active"
}

resource "netbox_manufacturer" "test" {
  name = "Test Manufacturer for Power Port DS"
  slug = "test-manufacturer-power-port-ds"
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = "Test Model for Power Port DS"
  slug         = "test-model-power-port-ds"
}

resource "netbox_device_role" "test" {
  name  = "Test Role for Power Port DS"
  slug  = "test-role-power-port-ds"
  color = "aabbcc"
}

resource "netbox_device" "test" {
  name        = "Test Device for Power Port DS"
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
  status      = "active"
}

resource "netbox_power_port" "test" {
  name        = "PSU1-DS"
  device      = netbox_device.test.id
  description = "Test power port for data source"
}

# Test: Lookup power port by ID
data "netbox_power_port" "by_id" {
  id = netbox_power_port.test.id
}
