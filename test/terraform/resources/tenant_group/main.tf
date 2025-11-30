# Tenant Group Resource Test

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

# Test 1: Basic tenant group (root level)
resource "netbox_tenant_group" "root" {
  name        = "Test Root Tenant Group"
  slug        = "test-root-tenant-group"
  description = "A root-level tenant group for testing"
}

# Test 2: Child tenant group (with parent)
resource "netbox_tenant_group" "child" {
  name        = "Test Child Tenant Group"
  slug        = "test-child-tenant-group"
  parent      = netbox_tenant_group.root.id
  description = "A child tenant group under the root"
}

# Test 3: Grandchild tenant group (deeper hierarchy)
resource "netbox_tenant_group" "grandchild" {
  name        = "Test Grandchild Tenant Group"
  slug        = "test-grandchild-tenant-group"
  parent      = netbox_tenant_group.child.id
  description = "A grandchild tenant group for testing deep hierarchies"
}

# Test 4: Another root-level group (sibling)
resource "netbox_tenant_group" "sibling" {
  name        = "Test Sibling Tenant Group"
  slug        = "test-sibling-tenant-group"
  description = "A sibling root-level tenant group"
}

# Test 5: Tenant group with tenant assigned
resource "netbox_tenant_group" "with_tenant" {
  name        = "Test Group With Tenant"
  slug        = "test-group-with-tenant"
  description = "A tenant group that will have a tenant"
}

resource "netbox_tenant" "in_group" {
  name        = "Tenant In Test Group"
  slug        = "tenant-in-test-group"
  group       = netbox_tenant_group.with_tenant.id
  description = "A tenant assigned to the test group"
}
