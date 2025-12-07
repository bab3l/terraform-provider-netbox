# VLAN Group Data Source Test

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
resource "netbox_vlan_group" "test" {
  name        = "Test VLAN Group DS"
  slug        = "test-vlan-group-ds"
  description = "Test VLAN group for data source"
}

# Test: Lookup VLAN group by ID
data "netbox_vlan_group" "by_id" {
  id = netbox_vlan_group.test.id
}

# Test: Lookup VLAN group by name
data "netbox_vlan_group" "by_name" {
  name = netbox_vlan_group.test.name

  depends_on = [netbox_vlan_group.test]
}
