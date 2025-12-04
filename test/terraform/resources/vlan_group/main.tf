# VLAN Group Integration Test
# Tests the netbox_vlan_group resource with basic and complete configurations

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

# Basic VLAN Group with only required fields
resource "netbox_vlan_group" "basic" {
  name = "Basic VLAN Group"
  slug = "basic-vlan-group"
}

# Complete VLAN Group with all optional fields
resource "netbox_vlan_group" "complete" {
  name        = "Complete VLAN Group"
  slug        = "complete-vlan-group"
  description = "Complete VLAN group for integration testing"
}

# Site for VLAN Group association
resource "netbox_site" "test" {
  name        = "VLAN Group Test Site"
  slug        = "vlan-group-test-site"
  status      = "active"
  description = "Site for VLAN group testing"
}

# VLAN Group scoped to site
resource "netbox_vlan_group" "site_scoped" {
  name        = "Site Scoped VLAN Group"
  slug        = "site-scoped-vlan-group"
  description = "VLAN group scoped to a specific site"
  scope_type  = "dcim.site"
  scope_id    = netbox_site.test.id
}
