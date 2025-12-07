# Console Port Resource Test

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
  name   = "Test Site for Console Port"
  slug   = "test-site-console-port"
  status = "active"
}

resource "netbox_manufacturer" "test" {
  name = "Test Manufacturer for Console Port"
  slug = "test-mfg-console-port"
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = "Test Device Type Console Port"
  slug         = "test-dt-console-port"
}

resource "netbox_device_role" "test" {
  name  = "Test Role for Console Port"
  slug  = "test-role-console-port"
  color = "aa1409"
}

resource "netbox_device" "test" {
  name        = "test-device-console-port"
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
}

# Test 1: Basic console port creation
resource "netbox_console_port" "basic" {
  device = netbox_device.test.id
  name   = "console0"
}

# Test 2: Console port with all optional fields
resource "netbox_console_port" "complete" {
  device      = netbox_device.test.id
  name        = "console1"
  label       = "Console Port 1"
  description = "Console port for testing"
  type        = "rj-45"
}
