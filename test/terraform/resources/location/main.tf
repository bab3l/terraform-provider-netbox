# Location Resource Test

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

# First, create required parent resources (site is required for location)
resource "netbox_site" "location_test_site" {
  name   = "Location Test Site"
  slug   = "location-test-site"
  status = "active"
}

resource "netbox_tenant" "location_test_tenant" {
  name = "Location Test Tenant"
  slug = "location-test-tenant"
}

# Basic location with only required fields
resource "netbox_location" "basic" {
  name = "Basic Test Location"
  slug = "basic-test-location"
  site = netbox_site.location_test_site.id
}

# Complete location with all optional fields
resource "netbox_location" "complete" {
  name        = "Complete Test Location"
  slug        = "complete-test-location"
  site        = netbox_site.location_test_site.id
  status      = "active"
  tenant      = netbox_tenant.location_test_tenant.id
  description = "A complete location with all fields configured"

  depends_on = [
    netbox_site.location_test_site,
    netbox_tenant.location_test_tenant
  ]
}

# Parent location for testing nested hierarchy
resource "netbox_location" "parent" {
  name = "Parent Location"
  slug = "parent-location"
  site = netbox_site.location_test_site.id
}

# Child location with parent
resource "netbox_location" "child" {
  name   = "Child Location"
  slug   = "child-location"
  site   = netbox_site.location_test_site.id
  parent = netbox_location.parent.id
}

# Grandchild location - testing nested hierarchy
resource "netbox_location" "grandchild" {
  name        = "Grandchild Location"
  slug        = "grandchild-location"
  site        = netbox_site.location_test_site.id
  parent      = netbox_location.child.id
  description = "Three levels deep in location hierarchy"
}

# Location with planned status
resource "netbox_location" "planned" {
  name   = "Planned Location"
  slug   = "planned-location"
  site   = netbox_site.location_test_site.id
  status = "planned"
}

# Location with staging status
resource "netbox_location" "staging" {
  name   = "Staging Location"
  slug   = "staging-location"
  site   = netbox_site.location_test_site.id
  status = "staging"
}
