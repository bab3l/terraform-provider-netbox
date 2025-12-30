# Front Port Data Source Test

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
  name   = "Test Site for Front Port DS"
  slug   = "test-site-front-port-ds"
  status = "active"
}

resource "netbox_manufacturer" "test" {
  name = "Test Manufacturer for Front Port DS"
  slug = "test-manufacturer-front-port-ds"
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = "Test Model for Front Port DS"
  slug         = "test-model-front-port-ds"
}

resource "netbox_device_role" "test" {
  name  = "Test Role for Front Port DS"
  slug  = "test-role-front-port-ds"
  color = "aa1409"
}

resource "netbox_device" "test" {
  name        = "test-device-front-port-ds"
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
}

resource "netbox_rear_port" "test" {
  name      = "RearPort0-DS"
  device    = netbox_device.test.id
  type      = "8p8c"
  positions = 2
}

resource "netbox_front_port" "test" {
  name        = "FrontPort0-DS"
  device      = netbox_device.test.id
  type        = "8p8c"
  rear_port   = netbox_rear_port.test.id
  description = "Test front port for data source"
}

# Test: Lookup front port by ID
data "netbox_front_port" "by_id" {
  id = netbox_front_port.test.id
}
