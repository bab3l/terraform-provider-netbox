# Module Bay Resource Test

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
  name   = "Test Site for Module Bay"
  slug   = "test-site-module-bay"
  status = "active"
}

resource "netbox_manufacturer" "test" {
  name = "Test Manufacturer for Module Bay"
  slug = "test-mfg-module-bay"
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = "Test Device Type Module Bay"
  slug         = "test-dt-module-bay"
}

resource "netbox_device_role" "test" {
  name  = "Test Role for Module Bay"
  slug  = "test-role-module-bay"
  color = "aa1409"
}

resource "netbox_device" "test" {
  name        = "test-device-module-bay"
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
}

# Test 1: Basic module bay creation
resource "netbox_module_bay" "basic" {
  device = netbox_device.test.id
  name   = "Slot 1"
}

# Test 2: Module bay with all optional fields
resource "netbox_module_bay" "complete" {
  device      = netbox_device.test.id
  name        = "Slot 2"
  label       = "Module Bay 2"
  position    = "2"
  description = "Module bay for testing"
}
