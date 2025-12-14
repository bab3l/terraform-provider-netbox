# Module Data Source Test

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
  name   = "Test Site for Module DS"
  slug   = "test-site-module-ds"
  status = "active"
}

resource "netbox_manufacturer" "test" {
  name = "Test Manufacturer for Module DS"
  slug = "test-manufacturer-module-ds"
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = "Test Model for Module DS"
  slug         = "test-model-module-ds"
}

resource "netbox_device_role" "test" {
  name  = "Test Role for Module DS"
  slug  = "test-role-module-ds"
  color = "aabbcc"
}

resource "netbox_device" "test" {
  name        = "Test Device for Module DS"
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
  status      = "active"
}

resource "netbox_module_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = "Test Module Type DS"
}

resource "netbox_module_bay" "test" {
  name   = "Slot1-DS"
  device = netbox_device.test.id
}

resource "netbox_module" "test" {
  device      = netbox_device.test.id
  module_bay  = netbox_module_bay.test.id
  module_type = netbox_module_type.test.id
  status      = "active"
}

# Test: Lookup module by ID
data "netbox_module" "by_id" {
  id = netbox_module.test.id
}
