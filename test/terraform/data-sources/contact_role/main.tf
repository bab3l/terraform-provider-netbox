# Contact Role Data Source Test

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
resource "netbox_contact_role" "test" {
  name        = "Test Contact Role DS"
  slug        = "test-contact-role-ds"
  description = "Test contact role for data source"
}

# Test: Lookup contact role by ID
data "netbox_contact_role" "by_id" {
  id = netbox_contact_role.test.id
}

# Test: Lookup contact role by name
data "netbox_contact_role" "by_name" {
  name = netbox_contact_role.test.name

  depends_on = [netbox_contact_role.test]
}
