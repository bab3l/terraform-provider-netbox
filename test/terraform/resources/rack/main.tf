# Rack Resource Test

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

# First, create required parent resources (site is required for rack)
resource "netbox_site" "rack_test_site" {
  name   = "Rack Test Site"
  slug   = "rack-test-site"
  status = "active"
}

resource "netbox_location" "rack_test_location" {
  name = "Rack Test Location"
  slug = "rack-test-location"
  site = netbox_site.rack_test_site.id
}

resource "netbox_tenant" "rack_test_tenant" {
  name = "Rack Test Tenant"
  slug = "rack-test-tenant"
}

# Basic rack with only required fields
resource "netbox_rack" "basic" {
  name   = "Basic Test Rack"
  site   = netbox_site.rack_test_site.id
  status = "active"
}

# Complete rack with all optional fields
resource "netbox_rack" "complete" {
  name        = "Complete Test Rack"
  site        = netbox_site.rack_test_site.id
  location    = netbox_location.rack_test_location.id
  tenant      = netbox_tenant.rack_test_tenant.id
  status      = "active"
  serial      = "RACK-SERIAL-001"
  asset_tag   = "ASSET-RACK-001"
  u_height    = 48
  desc_units  = false
  outer_width = 600
  outer_depth = 1200
  outer_unit  = "mm"
  mounting_depth = 19
  weight      = 150.5
  max_weight  = 1000
  weight_unit = "kg"
  description = "A complete rack with all fields configured"

  depends_on = [
    netbox_site.rack_test_site,
    netbox_location.rack_test_location,
    netbox_tenant.rack_test_tenant
  ]
}

# Rack with descending units (numbered from top)
resource "netbox_rack" "descending" {
  name       = "Descending Units Rack"
  site       = netbox_site.rack_test_site.id
  status     = "active"
  u_height   = 42
  desc_units = true
}

# Rack with reserved status
resource "netbox_rack" "reserved" {
  name   = "Reserved Rack"
  site   = netbox_site.rack_test_site.id
  status = "reserved"
}

# Rack with planned status
resource "netbox_rack" "planned" {
  name   = "Planned Rack"
  site   = netbox_site.rack_test_site.id
  status = "planned"
}

# Rack with deprecated status
resource "netbox_rack" "deprecated" {
  name   = "Deprecated Rack"
  site   = netbox_site.rack_test_site.id
  status = "deprecated"
}

# Rack with outer dimensions in inches
resource "netbox_rack" "imperial" {
  name        = "Imperial Rack"
  site        = netbox_site.rack_test_site.id
  status      = "active"
  outer_width = 24
  outer_depth = 48
  outer_unit  = "in"
}

# Rack with weight in pounds
resource "netbox_rack" "weight_lb" {
  name        = "Weight LB Rack"
  site        = netbox_site.rack_test_site.id
  status      = "active"
  weight      = 330.0
  max_weight  = 2200
  weight_unit = "lb"
}

# Small rack (12U - often used for networking)
resource "netbox_rack" "small" {
  name     = "Small Network Rack"
  site     = netbox_site.rack_test_site.id
  status   = "active"
  u_height = 12
}

# Large rack (52U)
resource "netbox_rack" "large" {
  name     = "Large Data Center Rack"
  site     = netbox_site.rack_test_site.id
  status   = "active"
  u_height = 52
}
