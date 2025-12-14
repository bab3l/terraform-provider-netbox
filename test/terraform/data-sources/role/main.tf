# Role (IPAM Role) Data Source Test

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
resource "netbox_role" "test" {
  name        = "Test IPAM Role DS"
  slug        = "test-ipam-role-ds"
  description = "Test IPAM role for data source"
}

# Test: Lookup role by ID
data "netbox_role" "by_id" {
  id = netbox_role.test.id
}

# Test: Lookup role by name
data "netbox_role" "by_name" {
  name = netbox_role.test.name

  depends_on = [netbox_role.test]
}
