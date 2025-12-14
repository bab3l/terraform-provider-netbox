# VRF Integration Test
# Tests the netbox_vrf resource with basic and complete configurations

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

# Basic VRF with only required fields
resource "netbox_vrf" "basic" {
  name = "Basic Test VRF"
}

# Complete VRF with all optional fields
resource "netbox_vrf" "complete" {
  name        = "Complete Test VRF"
  rd          = "65000:100"
  description = "Complete VRF for integration testing"
  comments    = "Created by terraform integration test"
}

# Test tenant for VRF association
resource "netbox_tenant" "test" {
  name        = "VRF Test Tenant"
  slug        = "vrf-test-tenant"
  description = "Tenant for VRF testing"
}

# VRF with tenant association
resource "netbox_vrf" "with_tenant" {
  name        = "VRF With Tenant"
  rd          = "65000:200"
  description = "VRF with tenant association"
  tenant      = netbox_tenant.test.id
}
