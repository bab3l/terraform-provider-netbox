# Device Bay Resource Test

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
  name   = "Test Site for Device Bay"
  slug   = "test-site-device-bay"
  status = "active"
}

resource "netbox_manufacturer" "test" {
  name = "Test Manufacturer for Device Bay"
  slug = "test-mfg-device-bay"
}

# Device type must support device bays (subdevice_role = parent)
resource "netbox_device_type" "test" {
  manufacturer   = netbox_manufacturer.test.id
  model          = "Test Device Type Device Bay"
  slug           = "test-dt-device-bay"
  subdevice_role = "parent"
}

resource "netbox_device_role" "test" {
  name  = "Test Role for Device Bay"
  slug  = "test-role-device-bay"
  color = "aa1409"
}

resource "netbox_device" "test" {
  name        = "test-device-device-bay"
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
}

# Test 1: Basic device bay creation
resource "netbox_device_bay" "basic" {
  device = netbox_device.test.id
  name   = "Bay 1"
}

# Test 2: Device bay with all optional fields
resource "netbox_device_bay" "complete" {
  device      = netbox_device.test.id
  name        = "Bay 2"
  label       = "Device Bay 2"
  description = "Device bay for testing"
}
