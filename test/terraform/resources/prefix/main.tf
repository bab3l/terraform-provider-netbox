# Prefix Integration Test
# Tests the netbox_prefix resource with basic and complete configurations

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
resource "netbox_vrf" "test" {
  name        = "Prefix Test VRF"
  rd          = "65000:300"
  description = "VRF for prefix testing"
}

resource "netbox_site" "test" {
  name        = "Prefix Test Site"
  slug        = "prefix-test-site"
  status      = "active"
  description = "Site for prefix testing"
}

resource "netbox_tenant" "test" {
  name = "Prefix Test Tenant"
  slug = "prefix-test-tenant"
}

resource "netbox_vlan" "test" {
  vid         = 500
  name        = "Prefix Test VLAN"
  description = "VLAN for prefix testing"
}

# Basic Prefix with only required fields
resource "netbox_prefix" "basic" {
  prefix = "10.0.0.0/24"
}

# Complete Prefix with all optional fields
resource "netbox_prefix" "complete" {
  prefix      = "10.1.0.0/24"
  status      = "active"
  description = "Complete prefix for integration testing"
  comments    = "Created by terraform integration test"
  vrf         = netbox_vrf.test.id
  site        = netbox_site.test.id
  tenant      = netbox_tenant.test.id
  vlan        = netbox_vlan.test.id
  is_pool     = true
}

# IPv6 Prefix
resource "netbox_prefix" "ipv6" {
  prefix      = "2001:db8::/32"
  status      = "active"
  description = "IPv6 prefix test"
}

# Container Prefix
resource "netbox_prefix" "container" {
  prefix      = "172.16.0.0/16"
  status      = "container"
  description = "Container prefix for subnetting"
}

# Reserved Prefix
resource "netbox_prefix" "reserved" {
  prefix      = "192.168.100.0/24"
  status      = "reserved"
  description = "Reserved prefix test"
}
