# Contact Group Data Source Test

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
resource "netbox_contact_group" "test" {
  name        = "Test Contact Group DS"
  slug        = "test-contact-group-ds"
  description = "Test contact group for data source"
}

# Test: Lookup contact group by ID
data "netbox_contact_group" "by_id" {
  id = netbox_contact_group.test.id
}

# Test: Lookup contact group by name
data "netbox_contact_group" "by_name" {
  name = netbox_contact_group.test.name

  depends_on = [netbox_contact_group.test]
}
