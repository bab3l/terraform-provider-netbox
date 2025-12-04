# VLAN Integration Test
# Tests the netbox_vlan resource with basic and complete configurations

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
resource "netbox_vlan_group" "test" {
  name        = "VLAN Test Group"
  slug        = "vlan-test-group"
  description = "VLAN group for VLAN testing"
}

resource "netbox_site" "test" {
  name        = "VLAN Test Site"
  slug        = "vlan-test-site"
  status      = "active"
  description = "Site for VLAN testing"
}

resource "netbox_tenant" "test" {
  name = "VLAN Test Tenant"
  slug = "vlan-test-tenant"
}

# Basic VLAN with only required fields
resource "netbox_vlan" "basic" {
  vid  = 100
  name = "Basic VLAN"
}

# Complete VLAN with group (note: group and site are mutually exclusive scopes)
resource "netbox_vlan" "complete" {
  vid         = 200
  name        = "Complete VLAN"
  status      = "active"
  description = "Complete VLAN for integration testing"
  comments    = "Created by terraform integration test"
  group       = netbox_vlan_group.test.id
  tenant      = netbox_tenant.test.id
}

# VLAN scoped to site (without group)
resource "netbox_vlan" "with_site" {
  vid         = 250
  name        = "Site Scoped VLAN"
  status      = "active"
  description = "VLAN scoped to a site"
  site        = netbox_site.test.id
}

# VLAN with different status values
resource "netbox_vlan" "reserved" {
  vid         = 300
  name        = "Reserved VLAN"
  status      = "reserved"
  description = "Reserved VLAN test"
}

resource "netbox_vlan" "deprecated" {
  vid         = 400
  name        = "Deprecated VLAN"
  status      = "deprecated"
  description = "Deprecated VLAN test"
}
