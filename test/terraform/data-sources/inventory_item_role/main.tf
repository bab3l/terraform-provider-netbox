# Inventory Item Role Data Source Test

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
resource "netbox_inventory_item_role" "test" {
  name        = "Test Inventory Role DS"
  slug        = "test-inventory-role-ds"
  color       = "aabbcc"
  description = "Test inventory item role for data source"
}

# Test: Lookup inventory item role by ID
data "netbox_inventory_item_role" "by_id" {
  id = netbox_inventory_item_role.test.id
}

# Test: Lookup inventory item role by name
data "netbox_inventory_item_role" "by_name" {
  name = netbox_inventory_item_role.test.name

  depends_on = [netbox_inventory_item_role.test]
}
