# Rack Data Source Test

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
  name   = "Rack DS Test Site"
  slug   = "rack-ds-test-site"
  status = "active"
}

resource "netbox_location" "test_location" {
  name = "Rack DS Test Location"
  slug = "rack-ds-test-location"
  site = netbox_site.test_site.id
}

resource "netbox_tenant" "test_tenant" {
  name = "Rack DS Test Tenant"
  slug = "rack-ds-test-tenant"
}

# Create a rack to look up with all fields
resource "netbox_rack" "test_rack" {
  name        = "Data Source Test Rack"
  site        = netbox_site.test_site.id
  location    = netbox_location.test_location.id
  tenant      = netbox_tenant.test_tenant.id
  status      = "active"
  serial      = "DS-RACK-001"
  asset_tag   = "ASSET-DS-001"
  u_height    = 42
  desc_units  = false
  description = "A rack for testing data source lookups"
}

# Create a second rack with different status
resource "netbox_rack" "reserved_rack" {
  name     = "Reserved DS Rack"
  site     = netbox_site.test_site.id
  status   = "reserved"
  u_height = 24
}

# Data source: look up rack by ID
data "netbox_rack" "by_id" {
  id = netbox_rack.test_rack.id

  depends_on = [netbox_rack.test_rack]
}

# Data source: look up rack by name
data "netbox_rack" "by_name" {
  name = netbox_rack.test_rack.name

  depends_on = [netbox_rack.test_rack]
}

# Data source: look up reserved rack
data "netbox_rack" "reserved_by_id" {
  id = netbox_rack.reserved_rack.id

  depends_on = [netbox_rack.reserved_rack]
}
