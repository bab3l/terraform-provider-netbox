# Inventory Item Role Resource Test

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

# Test 1: Basic inventory item role creation
resource "netbox_inventory_item_role" "basic" {
  name  = "Test Inventory Item Role Basic"
  slug  = "test-inventory-item-role-basic"
  color = "00ff00"
}

# Test 2: Inventory item role with all optional fields
resource "netbox_inventory_item_role" "complete" {
  name        = "Test Inventory Item Role Complete"
  slug        = "test-inventory-item-role-complete"
  color       = "ff0000"
  description = "An inventory item role for testing"
}
