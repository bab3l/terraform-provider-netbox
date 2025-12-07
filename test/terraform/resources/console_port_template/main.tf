# Console Port Template Resource Test

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
resource "netbox_manufacturer" "test" {
  name = "Test Manufacturer for Console Port Template"
  slug = "test-mfg-console-port-tpl"
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = "Test Device Type Console Port Template"
  slug         = "test-dt-console-port-tpl"
}

# Test 1: Basic console port template creation
resource "netbox_console_port_template" "basic" {
  device_type = netbox_device_type.test.id
  name        = "console0"
}

# Test 2: Console port template with all optional fields
resource "netbox_console_port_template" "complete" {
  device_type = netbox_device_type.test.id
  name        = "console1"
  label       = "Console Port Template 1"
  description = "Console port template for testing"
  type        = "rj-45"
}
