# Test: Inventory Item Template data source
# This tests looking up inventory item templates by various identifiers

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
  name = "Test Manufacturer for IIT DS"
  slug = "test-mfg-iit-ds"
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = "Test Device Type IIT DS"
  slug         = "test-dt-iit-ds"
}

resource "netbox_inventory_item_role" "test" {
  name  = "Test IIT DS Role"
  slug  = "test-iit-ds-role"
  color = "0066cc"
}

# Create inventory item template to look up
resource "netbox_inventory_item_template" "test" {
  device_type  = netbox_device_type.test.id
  name         = "Test IIT for DS"
  label        = "IIT-DS-1"
  role         = netbox_inventory_item_role.test.id
  manufacturer = netbox_manufacturer.test.id
  part_id      = "PART-DS-001"
  description  = "Test inventory item template for data source testing"
}

# Look up by ID
data "netbox_inventory_item_template" "by_id" {
  id = netbox_inventory_item_template.test.id
}
