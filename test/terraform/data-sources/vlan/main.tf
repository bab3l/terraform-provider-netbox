# VLAN Data Source Test

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
resource "netbox_vlan" "test" {
  vid         = 100
  name        = "Test VLAN DS"
  status      = "active"
  description = "Test VLAN for data source"
}

# Test: Lookup VLAN by ID
data "netbox_vlan" "by_id" {
  id = netbox_vlan.test.id
}
