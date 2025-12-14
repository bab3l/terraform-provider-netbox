# Contact Role Resource Integration Tests
# Tests the netbox_contact_role resource CRUD operations

terraform {
  required_version = ">= 1.0"
  required_providers {
    netbox = {
      source = "bab3l/netbox"
    }
  }
}

provider "netbox" {}

# Test 1: Basic contact role with minimal required fields
resource "netbox_contact_role" "test_basic" {
  name = "Test Contact Role Basic"
  slug = "test-contact-role-basic"
}

# Test 2: Contact role with all optional fields
resource "netbox_contact_role" "test_full" {
  name        = "Test Contact Role Full"
  slug        = "test-contact-role-full"
  description = "A fully configured contact role for integration testing"
}

# Test 3: Contact role for operations
resource "netbox_contact_role" "test_operations" {
  name        = "Operations Contact"
  slug        = "operations-contact"
  description = "Contact role for operations team"
}

# Data source test - lookup by ID
data "netbox_contact_role" "by_id" {
  id = netbox_contact_role.test_basic.id
}

# Data source test - lookup by name
data "netbox_contact_role" "by_name" {
  name = netbox_contact_role.test_full.name

  depends_on = [netbox_contact_role.test_full]
}

# Data source test - lookup by slug
data "netbox_contact_role" "by_slug" {
  slug = netbox_contact_role.test_operations.slug

  depends_on = [netbox_contact_role.test_operations]
}

# Outputs for verification
output "basic_role_id_valid" {
  value = can(tonumber(netbox_contact_role.test_basic.id))
}

output "basic_role_name_valid" {
  value = netbox_contact_role.test_basic.name == "Test Contact Role Basic"
}

output "full_role_description_valid" {
  value = netbox_contact_role.test_full.description == "A fully configured contact role for integration testing"
}

output "data_by_id_matches" {
  value = data.netbox_contact_role.by_id.name == netbox_contact_role.test_basic.name
}

output "data_by_name_matches" {
  value = data.netbox_contact_role.by_name.slug == netbox_contact_role.test_full.slug
}

output "data_by_slug_matches" {
  value = data.netbox_contact_role.by_slug.description == netbox_contact_role.test_operations.description
}
