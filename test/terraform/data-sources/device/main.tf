# Device Data Source Test

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
  name   = "Test Site for Device DS"
  slug   = "test-site-device-ds"
  status = "active"
}

resource "netbox_manufacturer" "test" {
  name = "Test Manufacturer for Device DS"
  slug = "test-manufacturer-device-ds"
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = "Test Model for Device DS"
  slug         = "test-model-device-ds"
}

resource "netbox_device_role" "test" {
  name  = "Test Role for Device DS"
  slug  = "test-role-device-ds"
  color = "aabbcc"
}

resource "netbox_device" "test" {
  name        = "Test Device DS"
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
  status      = "active"
  description = "Test device for data source"
}

# Test: Lookup device by ID
data "netbox_device" "by_id" {
  id = netbox_device.test.id
}

# Test: Lookup device by name
data "netbox_device" "by_name" {
  name = netbox_device.test.name

  depends_on = [netbox_device.test]
}
