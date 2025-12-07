# Module Resource Test

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
  name   = "Test Site for Module"
  slug   = "test-site-module"
  status = "active"
}

resource "netbox_manufacturer" "test" {
  name = "Test Manufacturer for Module"
  slug = "test-mfg-module"
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = "Test Device Type Module"
  slug         = "test-dt-module"
}

resource "netbox_module_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = "Test Module Type"
}

resource "netbox_device_role" "test" {
  name  = "Test Role for Module"
  slug  = "test-role-module"
  color = "aa1409"
}

resource "netbox_device" "test" {
  name        = "test-device-module"
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
}

resource "netbox_module_bay" "test" {
  device = netbox_device.test.id
  name   = "Slot 1"
}

# Test 1: Basic module creation
resource "netbox_module" "basic" {
  device      = netbox_device.test.id
  module_bay  = netbox_module_bay.test.id
  module_type = netbox_module_type.test.id
}
