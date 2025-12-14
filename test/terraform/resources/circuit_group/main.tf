# Circuit Group Resource Test

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

# Create a tenant to use for testing
resource "netbox_tenant" "test" {
  name = "test-circuit-group-tenant"
  slug = "test-circuit-group-tenant"
}

# Test 1: Basic circuit group with required fields only
resource "netbox_circuit_group" "basic" {
  name = "test-circuit-group-basic"
  slug = "test-circuit-group-basic"
}

# Test 2: Circuit group with all optional fields
resource "netbox_circuit_group" "complete" {
  name        = "test-circuit-group-complete"
  slug        = "test-circuit-group-complete"
  description = "Complete circuit group for integration testing"
  tenant      = netbox_tenant.test.id
}

# Test 3: Circuit group with tags
resource "netbox_tag" "test_tag" {
  name = "circuit-group-test"
  slug = "circuit-group-test"
}

resource "netbox_circuit_group" "with_tags" {
  name        = "test-circuit-group-tagged"
  slug        = "test-circuit-group-tagged"
  description = "Circuit group with tags"

  tags = [
    {
      name = netbox_tag.test_tag.name
      slug = netbox_tag.test_tag.slug
    }
  ]
}

# Test 4: Output values for verification
output "basic_circuit_group_id" {
  value = netbox_circuit_group.basic.id
}

output "basic_circuit_group_name" {
  value = netbox_circuit_group.basic.name
}

output "basic_circuit_group_slug" {
  value = netbox_circuit_group.basic.slug
}

output "complete_circuit_group_id" {
  value = netbox_circuit_group.complete.id
}

output "complete_circuit_group_description" {
  value = netbox_circuit_group.complete.description
}

output "complete_circuit_group_tenant" {
  value = netbox_circuit_group.complete.tenant
}

output "with_tags_circuit_group_id" {
  value = netbox_circuit_group.with_tags.id
}
