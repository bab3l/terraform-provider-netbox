# Interface Template Resource Test

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
  name = "Test Manufacturer for Interface Template"
  slug = "test-mfg-interface-tpl"
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = "Test Device Type Interface Template"
  slug         = "test-dt-interface-tpl"
}

# Test 1: Basic interface template creation
resource "netbox_interface_template" "basic" {
  device_type = netbox_device_type.test.id
  name        = "eth0"
  type        = "1000base-t"
}

# Test 2: Interface template with all optional fields
resource "netbox_interface_template" "complete" {
  device_type = netbox_device_type.test.id
  name        = "eth1"
  type        = "10gbase-t"
  label       = "Ethernet 1"
  description = "Interface template for testing"
  mgmt_only   = false
}
