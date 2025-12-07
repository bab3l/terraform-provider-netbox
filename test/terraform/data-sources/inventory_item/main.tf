# Inventory Item Data Source Test

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
  name   = "Test Site for Inventory Item DS"
  slug   = "test-site-inventory-item-ds"
  status = "active"
}

resource "netbox_manufacturer" "test" {
  name = "Test Manufacturer for Inventory Item DS"
  slug = "test-manufacturer-inventory-item-ds"
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = "Test Model for Inventory Item DS"
  slug         = "test-model-inventory-item-ds"
}

resource "netbox_device_role" "test" {
  name  = "Test Role for Inventory Item DS"
  slug  = "test-role-inventory-item-ds"
  color = "aabbcc"
}

resource "netbox_device" "test" {
  name        = "Test Device for Inventory Item DS"
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
  status      = "active"
}

resource "netbox_inventory_item" "test" {
  name        = "SFP-DS"
  device      = netbox_device.test.id
  description = "Test inventory item for data source"
}

# Test: Lookup inventory item by ID
data "netbox_inventory_item" "by_id" {
  id = netbox_inventory_item.test.id
}
