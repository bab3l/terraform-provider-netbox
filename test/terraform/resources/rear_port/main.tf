# Rear Port Resource Test
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
  name   = "Test Site for Rear Port"
  slug   = "test-site-rear-port"
  status = "active"
}

resource "netbox_manufacturer" "test" {
  name = "Test Manufacturer for Rear Port"
  slug = "test-mfg-rear-port"
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = "Test Device Type Rear Port"
  slug         = "test-dt-rear-port"
}

resource "netbox_device_role" "test" {
  name  = "Test Role for Rear Port"
  slug  = "test-role-rear-port"
  color = "aa1409"
}

resource "netbox_device" "test" {
  name        = "test-device-rear-port"
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
}

# Test 1: Basic rear port creation
resource "netbox_rear_port" "basic" {
  device = netbox_device.test.id
  name   = "rear0"
  type   = "8p8c"
}

# Test 2: Rear port with all optional fields
resource "netbox_rear_port" "complete" {
  device         = netbox_device.test.id
  name           = "rear1"
  type           = "lc"
  label          = "Rear Port 1"
  color          = "aa1409"
  positions      = 4
  description    = "Rear port for testing"
  mark_connected = true
}
