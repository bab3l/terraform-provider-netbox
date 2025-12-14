# Location Data Source Test

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

# First, create required parent resources
resource "netbox_site" "test_site" {
  name   = "Location DS Test Site"
  slug   = "location-ds-test-site"
  status = "active"
}

resource "netbox_tenant" "test_tenant" {
  name = "Location DS Test Tenant"
  slug = "location-ds-test-tenant"
}

# Create a location to look up
resource "netbox_location" "test_location" {
  name        = "Data Source Test Location"
  slug        = "data-source-test-location"
  site        = netbox_site.test_site.id
  status      = "active"
  tenant      = netbox_tenant.test_tenant.id
  description = "A location for testing data source lookups"
}

# Create a child location for hierarchy testing
resource "netbox_location" "child_location" {
  name        = "Data Source Child Location"
  slug        = "data-source-child-location"
  site        = netbox_site.test_site.id
  parent      = netbox_location.test_location.id
  description = "A child location for testing data source lookups"
}

# Data source: look up location by ID
data "netbox_location" "by_id" {
  id = netbox_location.test_location.id

  depends_on = [netbox_location.test_location]
}

# Data source: look up location by name
data "netbox_location" "by_name" {
  name = netbox_location.test_location.name

  depends_on = [netbox_location.test_location]
}

# Data source: look up location by slug
data "netbox_location" "by_slug" {
  slug = netbox_location.test_location.slug

  depends_on = [netbox_location.test_location]
}

# Data source: look up child location to verify parent relationship
data "netbox_location" "child_by_id" {
  id = netbox_location.child_location.id

  depends_on = [netbox_location.child_location]
}
