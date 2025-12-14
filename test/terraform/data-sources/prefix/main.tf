# Prefix Data Source Test

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
resource "netbox_prefix" "test" {
  prefix      = "10.0.0.0/8"
  status      = "active"
  description = "Test prefix for data source"
}

# Test: Lookup prefix by ID
data "netbox_prefix" "by_id" {
  id = netbox_prefix.test.id
}

# Test: Lookup prefix by prefix
data "netbox_prefix" "by_prefix" {
  prefix = netbox_prefix.test.prefix

  depends_on = [netbox_prefix.test]
}
