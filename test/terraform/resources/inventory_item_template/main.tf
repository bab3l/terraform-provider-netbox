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
  name = "Test Manufacturer for Inventory Item Template"
  slug = "test-mfg-iit"
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = "Test Device Type Inventory Item Template"
  slug         = "test-dt-iit"
}

resource "netbox_inventory_item_role" "test" {
  name  = "Test Inventory Role"
  slug  = "test-role-iit"
  color = "0066cc"
}

# Test 1: Basic inventory item template
resource "netbox_inventory_item_template" "basic" {
  device_type = netbox_device_type.test.id
  name        = "SFP Module"
}

# Test 2: Inventory item template with parent
resource "netbox_inventory_item_template" "parent" {
  device_type = netbox_device_type.test.id
  name        = "Parent Component"
}

resource "netbox_inventory_item_template" "child" {
  device_type = netbox_device_type.test.id
  name        = "Child Component"
  parent      = netbox_inventory_item_template.parent.id
}

# Test 3: Inventory item template with all optional fields
resource "netbox_inventory_item_template" "complete" {
  device_type  = netbox_device_type.test.id
  name         = "Power Supply"
  label        = "PSU-1"
  role         = netbox_inventory_item_role.test.id
  manufacturer = netbox_manufacturer.test.id
  part_id      = "PS-001"
  description  = "Test inventory item template with all fields"
}
