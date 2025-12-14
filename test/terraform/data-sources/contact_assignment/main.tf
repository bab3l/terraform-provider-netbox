# Test: Contact Assignment data source
# This tests looking up contact assignments by ID

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

# Create dependencies
resource "netbox_site" "test" {
  name   = "Test Site for CA DS"
  slug   = "test-site-ca-ds"
  status = "active"
}

resource "netbox_contact" "test" {
  name  = "Test Contact for CA DS"
  email = "contact-assignment-ds-test@example.com"
}

resource "netbox_contact_role" "test" {
  name = "Test CA DS Role"
  slug = "test-ca-ds-role"
}

# Create contact assignment to look up
resource "netbox_contact_assignment" "test" {
  object_type = "dcim.site"
  object_id   = netbox_site.test.id
  contact_id  = netbox_contact.test.id
  role_id     = netbox_contact_role.test.id
  priority    = "primary"
}

# Look up by ID
data "netbox_contact_assignment" "by_id" {
  id = netbox_contact_assignment.test.id
}
