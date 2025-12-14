# Device Bay Data Source Test

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
  name   = "Test Site for Device Bay DS"
  slug   = "test-site-device-bay-ds"
  status = "active"
}

resource "netbox_manufacturer" "test" {
  name = "Test Manufacturer for Device Bay DS"
  slug = "test-manufacturer-device-bay-ds"
}

resource "netbox_device_type" "test" {
  manufacturer  = netbox_manufacturer.test.id
  model         = "Test Model for Device Bay DS"
  slug          = "test-model-device-bay-ds"
  subdevice_role = "parent"
}

resource "netbox_device_role" "test" {
  name  = "Test Role for Device Bay DS"
  slug  = "test-role-device-bay-ds"
  color = "aabbcc"
}

resource "netbox_device" "test" {
  name        = "Test Device for Device Bay DS"
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
  status      = "active"
}

resource "netbox_device_bay" "test" {
  name        = "Bay0-DS"
  device      = netbox_device.test.id
  description = "Test device bay for data source"
}

# Test: Lookup device bay by ID
data "netbox_device_bay" "by_id" {
  id = netbox_device_bay.test.id
}
