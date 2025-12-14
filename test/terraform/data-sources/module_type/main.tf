# Module Type Data Source Test

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
  name = "Test Manufacturer for Module Type DS"
  slug = "test-manufacturer-module-type-ds"
}

resource "netbox_module_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = "Test Module Type DS"
  description  = "Test module type for data source"
}

# Test: Lookup module type by ID
data "netbox_module_type" "by_id" {
  id = netbox_module_type.test.id
}
