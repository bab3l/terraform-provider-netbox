# Console Server Port Resource Test

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
  name   = "Test Site for Console Server Port"
  slug   = "test-site-console-server-port"
  status = "active"
}

resource "netbox_manufacturer" "test" {
  name = "Test Manufacturer for Console Server Port"
  slug = "test-mfg-console-server-port"
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = "Test Device Type Console Server Port"
  slug         = "test-dt-console-server-port"
}

resource "netbox_device_role" "test" {
  name  = "Test Role for Console Server Port"
  slug  = "test-role-console-server-port"
  color = "aa1409"
}

resource "netbox_device" "test" {
  name        = "test-device-console-server-port"
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
}

# Test 1: Basic console server port creation
resource "netbox_console_server_port" "basic" {
  device = netbox_device.test.id
  name   = "ttyS0"
}

# Test 2: Console server port with all optional fields
resource "netbox_console_server_port" "complete" {
  device      = netbox_device.test.id
  name        = "ttyS1"
  label       = "Console Server Port 1"
  description = "Console server port for testing"
  type        = "rj-45"
}
