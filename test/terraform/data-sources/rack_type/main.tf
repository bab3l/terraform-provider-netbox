# Rack Type Data Source Test

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
  name = "Test Manufacturer for Rack Type DS"
  slug = "test-manufacturer-rack-type-ds"
}

resource "netbox_rack_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = "Test Rack Type DS"
  slug         = "test-rack-type-ds"
  form_factor  = "4-post-cabinet"
  width        = 19
  u_height     = 42
  description  = "Test rack type for data source"
}

# Test: Lookup rack type by ID
data "netbox_rack_type" "by_id" {
  id = netbox_rack_type.test.id
}
