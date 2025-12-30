# Front Port Resource Test

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
  name   = "Test Site for Front Port"
  slug   = "test-site-front-port"
  status = "active"
}

resource "netbox_manufacturer" "test" {
  name = "Test Manufacturer for Front Port"
  slug = "test-mfg-front-port"
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = "Test Device Type Front Port"
  slug         = "test-dt-front-port"
}

resource "netbox_device_role" "test" {
  name  = "Test Role for Front Port"
  slug  = "test-role-front-port"
  color = "aa1409"
}

resource "netbox_device" "test" {
  name        = "test-device-front-port"
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
}

resource "netbox_rear_port" "test" {
  device    = netbox_device.test.id
  name      = "rear0"
  type      = "8p8c"
  positions = 4
}

# Test 1: Basic front port creation
resource "netbox_front_port" "basic" {
  device    = netbox_device.test.id
  name      = "front0"
  type      = "8p8c"
  rear_port = netbox_rear_port.test.id
}

# Test 2: Front port with all optional fields
resource "netbox_front_port" "complete" {
  device             = netbox_device.test.id
  name               = "front1"
  type               = "8p8c"
  rear_port          = netbox_rear_port.test.id
  rear_port_position = 2
  label              = "Front Port 1"
  color              = "aa1409"
  description        = "Front port for testing"
  mark_connected     = true
}
