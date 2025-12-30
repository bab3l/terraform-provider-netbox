# Rear Port Template Resource Test

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
  name = "Test Manufacturer for Rear Port Template"
  slug = "test-mfg-rear-port-tpl"
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = "Test Device Type Rear Port Template"
  slug         = "test-dt-rear-port-tpl"
}

# Test 1: Basic rear port template creation
resource "netbox_rear_port_template" "basic" {
  device_type = netbox_device_type.test.id
  name        = "rear0"
  type        = "8p8c"
}

# Test 2: Rear port template with all optional fields
resource "netbox_rear_port_template" "complete" {
  device_type = netbox_device_type.test.id
  name        = "rear1"
  type        = "lc"
  label       = "Rear Port Template 1"
  color       = "aa1409"
  positions   = 4
  description = "Rear port template for testing"
}
