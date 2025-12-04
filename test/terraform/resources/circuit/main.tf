# Circuit Integration Test
# Tests the netbox_circuit resource with basic and complete configurations

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

# Prerequisites
resource "netbox_provider" "test" {
  name = "Circuit Test Provider"
  slug = "circuit-test-provider"
}

resource "netbox_circuit_type" "test" {
  name = "Circuit Test Type"
  slug = "circuit-test-type"
}

resource "netbox_tenant" "test" {
  name = "Circuit Test Tenant"
  slug = "circuit-test-tenant"
}

# Basic Circuit with only required fields
resource "netbox_circuit" "basic" {
  cid              = "CKT-BASIC-001"
  circuit_provider = netbox_provider.test.slug
  type             = netbox_circuit_type.test.slug
}

# Complete Circuit with all optional fields
resource "netbox_circuit" "complete" {
  cid              = "CKT-COMPLETE-001"
  circuit_provider = netbox_provider.test.slug
  type             = netbox_circuit_type.test.slug
  status           = "active"
  description      = "Complete circuit for integration testing"
  comments         = "Created by terraform integration test"
  tenant           = netbox_tenant.test.id
  commit_rate      = 100000
}

# Active circuit for status testing
resource "netbox_circuit" "active" {
  cid              = "CKT-ACTIVE-001"
  circuit_provider = netbox_provider.test.slug
  type             = netbox_circuit_type.test.slug
  status           = "active"
  description      = "Active circuit test"
}

# Planned circuit for status testing
resource "netbox_circuit" "planned" {
  cid              = "CKT-PLANNED-001"
  circuit_provider = netbox_provider.test.slug
  type             = netbox_circuit_type.test.slug
  status           = "planned"
  description      = "Planned circuit test"
}
