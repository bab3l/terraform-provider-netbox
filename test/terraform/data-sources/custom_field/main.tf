# Custom Field Data Source Test

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
resource "netbox_custom_field" "test" {
  name          = "test_cf_ds"
  type          = "text"
  object_types  = ["dcim.site"]
  description   = "Test custom field for data source"
}

# Test: Lookup custom field by ID
data "netbox_custom_field" "by_id" {
  id = netbox_custom_field.test.id
}

# Test: Lookup custom field by name
data "netbox_custom_field" "by_name" {
  name = netbox_custom_field.test.name

  depends_on = [netbox_custom_field.test]
}
