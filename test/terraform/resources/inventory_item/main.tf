# Inventory Item Resource Test

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
  name   = "Test Site for Inventory Item"
  slug   = "test-site-inventory-item"
  status = "active"
}

resource "netbox_manufacturer" "test" {
  name = "Test Manufacturer for Inventory Item"
  slug = "test-mfg-inventory-item"
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = "Test Device Type Inventory Item"
  slug         = "test-dt-inventory-item"
}

resource "netbox_device_role" "test" {
  name  = "Test Role for Inventory Item"
  slug  = "test-role-inventory-item"
  color = "aa1409"
}

resource "netbox_device" "test" {
  name        = "test-device-inventory-item"
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
}

# Test 1: Basic inventory item creation
resource "netbox_inventory_item" "basic" {
  device = netbox_device.test.id
  name   = "SFP Module"
}

# Test 2: Inventory item with all optional fields
resource "netbox_inventory_item" "complete" {
  device       = netbox_device.test.id
  name         = "Power Supply"
  label        = "PSU-1"
  description  = "Inventory item for testing"
  manufacturer = netbox_manufacturer.test.id
  part_id      = "PS-001"
  serial       = "SN123456"
  asset_tag    = "ASSET001"
}
