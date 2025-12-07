# VRF Data Source Test

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
resource "netbox_vrf" "test" {
  name        = "Test VRF DS"
  rd          = "65000:1"
  description = "Test VRF for data source"
}

# Test: Lookup VRF by ID
data "netbox_vrf" "by_id" {
  id = netbox_vrf.test.id
}

# Test: Lookup VRF by name
data "netbox_vrf" "by_name" {
  name = netbox_vrf.test.name

  depends_on = [netbox_vrf.test]
}
