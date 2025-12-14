# Test: Contact data source
# This tests looking up contacts by various identifiers

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

variable "test_id" {
  description = "Unique identifier for this test run"
  type        = string
  default     = "test"
}

# Create a contact to look up
resource "netbox_contact" "test" {
  name        = "Test Contact DS ${var.test_id}"
  title       = "Test Title"
  phone       = "+1-555-0111"
  email       = "ds-test-${var.test_id}@example.com"
  description = "Test contact for data source testing"
}

# Look up by ID
data "netbox_contact" "by_id" {
  id = netbox_contact.test.id
}

# Look up by name
data "netbox_contact" "by_name" {
  name = netbox_contact.test.name
}

# Look up by email
data "netbox_contact" "by_email" {
  email = netbox_contact.test.email
}

# Verification outputs
output "lookup_by_id_name" {
  value = data.netbox_contact.by_id.name
}

output "lookup_by_name_email" {
  value = data.netbox_contact.by_name.email
}

output "lookup_by_email_phone" {
  value = data.netbox_contact.by_email.phone
}

output "id_matches" {
  value = data.netbox_contact.by_id.id == netbox_contact.test.id
}

output "name_matches" {
  value = data.netbox_contact.by_name.name == netbox_contact.test.name
}

output "email_matches" {
  value = data.netbox_contact.by_email.email == netbox_contact.test.email
}
