# Test: Contact resource
# This creates a contact with various optional fields

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

# Basic contact
resource "netbox_contact" "basic" {
  name  = "Test Contact Basic ${var.test_id}"
  email = "basic-${var.test_id}@example.com"
}

# Contact with full details
resource "netbox_contact" "full" {
  name        = "Test Contact Full ${var.test_id}"
  title       = "Senior Network Engineer"
  phone       = "+1-555-0199"
  email       = "full-${var.test_id}@example.com"
  address     = "456 Test Avenue, Test City, TC 12345"
  link        = "https://example.com/contact"
  description = "Test contact with all fields"
  comments    = "This is a test contact created for integration testing"
}

# Data source tests
data "netbox_contact" "by_id" {
  id = netbox_contact.basic.id
}

data "netbox_contact" "by_name" {
  name = netbox_contact.full.name
}

data "netbox_contact" "by_email" {
  email = netbox_contact.basic.email
}

# Outputs for verification
output "basic_contact_id" {
  value = netbox_contact.basic.id
}

output "full_contact_id" {
  value = netbox_contact.full.id
}

output "datasource_by_id_name" {
  value = data.netbox_contact.by_id.name
}

output "datasource_by_name_title" {
  value = data.netbox_contact.by_name.title
}

output "datasource_by_email_name" {
  value = data.netbox_contact.by_email.name
}
