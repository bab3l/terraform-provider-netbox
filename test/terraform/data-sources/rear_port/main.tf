# Rear Port Data Source Test

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
  name   = "Test Site for Rear Port DS"
  slug   = "test-site-rear-port-ds"
  status = "active"
}

resource "netbox_manufacturer" "test" {
  name = "Test Manufacturer for Rear Port DS"
  slug = "test-manufacturer-rear-port-ds"
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = "Test Model for Rear Port DS"
  slug         = "test-model-rear-port-ds"
}

resource "netbox_device_role" "test" {
  name  = "Test Role for Rear Port DS"
  slug  = "test-role-rear-port-ds"
  color = "aa1409"
}

resource "netbox_device" "test" {
  name        = "test-device-rear-port-ds"
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
}

resource "netbox_rear_port" "test" {
  name        = "RearPort0-DS"
  device      = netbox_device.test.id
  type        = "8p8c"
  positions   = 2
  description = "Test rear port for data source"
}

# Test: Lookup rear port by ID
data "netbox_rear_port" "by_id" {
  id = netbox_rear_port.test.id
}
