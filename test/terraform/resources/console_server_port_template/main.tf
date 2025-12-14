# Console Server Port Template Resource Test

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
resource "netbox_manufacturer" "test" {
  name = "Test Manufacturer for Console Server Port Tpl"
  slug = "test-mfg-console-server-port-tpl"
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = "Test Device Type Console Server Port Tpl"
  slug         = "test-dt-console-server-port-tpl"
}

# Test 1: Basic console server port template creation
resource "netbox_console_server_port_template" "basic" {
  device_type = netbox_device_type.test.id
  name        = "ttyS0"
}

# Test 2: Console server port template with all optional fields
resource "netbox_console_server_port_template" "complete" {
  device_type = netbox_device_type.test.id
  name        = "ttyS1"
  label       = "Console Server Port Template 1"
  description = "Console server port template for testing"
  type        = "rj-45"
}
