# Test Configuration for netbox_contact_group resource
# Creates contact groups with various configurations to test CRUD operations

terraform {
  required_version = ">= 1.0"
  required_providers {
    netbox = {
      source = "bab3l/netbox"
    }
  }
}

# Basic contact group with minimal fields
resource "netbox_contact_group" "basic" {
  name = "Test Contact Group Basic test"
  slug = "test-contact-group-basic-test"
}

# Contact group with description
resource "netbox_contact_group" "with_description" {
  name        = "Test Contact Group Full test"
  slug        = "test-contact-group-full-test"
  description = "Test contact group with full configuration"
}

# Parent contact group for hierarchy testing
resource "netbox_contact_group" "parent" {
  name        = "Test Contact Group Parent test"
  slug        = "test-contact-group-parent-test"
  description = "Parent group for hierarchy testing"
}

# Child contact group with parent reference
resource "netbox_contact_group" "child" {
  name        = "Test Contact Group Child test"
  slug        = "test-contact-group-child-test"
  parent      = netbox_contact_group.parent.id
  description = "Child group for hierarchy testing"
}

# Data source lookups
data "netbox_contact_group" "by_id" {
  id = netbox_contact_group.basic.id
}

data "netbox_contact_group" "by_name" {
  name = netbox_contact_group.with_description.name

  depends_on = [netbox_contact_group.with_description]
}

data "netbox_contact_group" "by_slug" {
  slug = netbox_contact_group.basic.slug

  depends_on = [netbox_contact_group.basic]
}

# Validation outputs
output "basic_id_valid" {
  value = netbox_contact_group.basic.id != ""
}

output "basic_name_valid" {
  value = netbox_contact_group.basic.name == "Test Contact Group Basic test"
}

output "with_description_valid" {
  value = netbox_contact_group.with_description.description == "Test contact group with full configuration"
}

output "hierarchy_valid" {
  value = netbox_contact_group.child.parent == netbox_contact_group.parent.id
}

output "datasource_by_id_matches" {
  value = data.netbox_contact_group.by_id.name == netbox_contact_group.basic.name
}

output "datasource_by_name_matches" {
  value = data.netbox_contact_group.by_name.id == netbox_contact_group.with_description.id
}

output "datasource_by_slug_matches" {
  value = data.netbox_contact_group.by_slug.id == netbox_contact_group.basic.id
}
