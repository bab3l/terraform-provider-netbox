# Power Port Resource Test

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
  name   = "Test Site for Power Port"
  slug   = "test-site-power-port"
  status = "active"
}

resource "netbox_manufacturer" "test" {
  name = "Test Manufacturer for Power Port"
  slug = "test-manufacturer-power-port"
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = "Test Model for Power Port"
  slug         = "test-model-power-port"
}

resource "netbox_device_role" "test" {
  name  = "Test Role for Power Port"
  slug  = "test-role-power-port"
  color = "aabbcc"
}

resource "netbox_device" "test" {
  name        = "Test Device for Power Port"
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
  status      = "active"
}

# Test 1: Basic power port creation
resource "netbox_power_port" "basic" {
  name   = "PSU1"
  device = netbox_device.test.id
}

# Test 2: Power port with all optional fields
resource "netbox_power_port" "complete" {
  name             = "PSU2"
  device           = netbox_device.test.id
  type             = "iec-60320-c14"
  maximum_draw     = 500
  allocated_draw   = 400
  mark_connected   = true
  description      = "A power port for testing"
}
