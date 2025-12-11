# Test: Contact Assignment resource
# This tests creating contact assignments that link contacts to Netbox objects

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

# Create dependencies
resource "netbox_site" "test" {
  name   = "Test Site for Contact Assignment"
  slug   = "test-site-ca"
  status = "active"
}

resource "netbox_contact" "test" {
  name  = "Test Contact for Assignment"
  email = "contact-assignment-test@example.com"
}

resource "netbox_contact_role" "test" {
  name = "Test Contact Role"
  slug = "test-contact-role-ca"
}

resource "netbox_contact_role" "secondary" {
  name = "Secondary Contact Role"
  slug = "secondary-contact-role-ca"
}

# Test 1: Contact assignment with role only
resource "netbox_contact_assignment" "basic" {
  object_type = "dcim.site"
  object_id   = netbox_site.test.id
  contact_id  = netbox_contact.test.id
  role_id     = netbox_contact_role.secondary.id
}

# Test 2: Contact assignment with role and priority
resource "netbox_contact_assignment" "with_role" {
  object_type = "dcim.site"
  object_id   = netbox_site.test.id
  contact_id  = netbox_contact.test.id
  role_id     = netbox_contact_role.test.id
  priority    = "primary"
}
