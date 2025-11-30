# Tenant Resource Test

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

# Test 1: Basic tenant
resource "netbox_tenant" "basic" {
  name        = "Test Tenant Basic"
  slug        = "test-tenant-basic"
  description = "A basic test tenant"
}

# Test 2: Tenant with all fields
resource "netbox_tenant" "complete" {
  name        = "Test Tenant Complete"
  slug        = "test-tenant-complete"
  description = "A complete test tenant with all optional fields"
  comments    = "This tenant was created for integration testing."
}

# Test 3: Tenant with tenant group (created inline)
resource "netbox_tenant_group" "for_tenant" {
  name        = "Tenant Test Group"
  slug        = "tenant-test-group"
  description = "Group created for tenant testing"
}

resource "netbox_tenant" "with_group" {
  name        = "Test Tenant With Group"
  slug        = "test-tenant-with-group"
  group       = netbox_tenant_group.for_tenant.id
  description = "A tenant with a tenant group"
}

# Test 4: Multiple tenants in same group
resource "netbox_tenant" "sibling1" {
  name        = "Test Tenant Sibling 1"
  slug        = "test-tenant-sibling-1"
  group       = netbox_tenant_group.for_tenant.id
  description = "First sibling tenant"
}

resource "netbox_tenant" "sibling2" {
  name        = "Test Tenant Sibling 2"
  slug        = "test-tenant-sibling-2"
  group       = netbox_tenant_group.for_tenant.id
  description = "Second sibling tenant"
}
