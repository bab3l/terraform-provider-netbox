# Virtual Chassis Data Source Test

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
resource "netbox_virtual_chassis" "test" {
  name        = "Test Virtual Chassis DS"
  description = "Test virtual chassis for data source"
}

# Test: Lookup virtual chassis by ID
data "netbox_virtual_chassis" "by_id" {
  id = netbox_virtual_chassis.test.id
}
