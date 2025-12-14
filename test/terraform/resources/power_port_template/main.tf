# Power Port Template Resource Test

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
  name = "Test Manufacturer for Power Port Template"
  slug = "test-manufacturer-power-port-template"
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = "Test Model for Power Port Template"
  slug         = "test-model-power-port-template"
}

# Test 1: Basic power port template creation
resource "netbox_power_port_template" "basic" {
  name        = "PSU-Template1"
  device_type = netbox_device_type.test.id
}

# Test 2: Power port template with all optional fields
resource "netbox_power_port_template" "complete" {
  name           = "PSU-Template2"
  device_type    = netbox_device_type.test.id
  type           = "iec-60320-c14"
  maximum_draw   = 500
  allocated_draw = 400
  label          = "Power Supply Unit Template"
  description    = "A power port template for testing"
}
