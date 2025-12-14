# Console Port Data Source Test

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
  name   = "Test Site for Console Port DS"
  slug   = "test-site-console-port-ds"
  status = "active"
}

resource "netbox_manufacturer" "test" {
  name = "Test Manufacturer for Console Port DS"
  slug = "test-manufacturer-console-port-ds"
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = "Test Model for Console Port DS"
  slug         = "test-model-console-port-ds"
}

resource "netbox_device_role" "test" {
  name  = "Test Role for Console Port DS"
  slug  = "test-role-console-port-ds"
  color = "aabbcc"
}

resource "netbox_device" "test" {
  name        = "Test Device for Console Port DS"
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
  status      = "active"
}

resource "netbox_console_port" "test" {
  name        = "Console0-DS"
  device      = netbox_device.test.id
  description = "Test console port for data source"
}

# Test: Lookup console port by ID
data "netbox_console_port" "by_id" {
  id = netbox_console_port.test.id
}
