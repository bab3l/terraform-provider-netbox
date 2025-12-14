# Rack Type Resource Test

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

# Dependencies
resource "netbox_manufacturer" "test" {
  name = "Test Manufacturer for Rack Type"
  slug = "test-manufacturer-rack-type"
}

# Test 1: Basic rack type creation
resource "netbox_rack_type" "basic" {
  manufacturer = netbox_manufacturer.test.id
  model        = "Test Rack Type Basic"
  slug         = "test-rack-type-basic"
  form_factor  = "2-post-frame"
  width        = 19
  u_height     = 42
}

# Test 2: Rack type with all optional fields
resource "netbox_rack_type" "complete" {
  manufacturer       = netbox_manufacturer.test.id
  model              = "Test Rack Type Complete"
  slug               = "test-rack-type-complete"
  form_factor        = "4-post-cabinet"
  width              = 23
  u_height           = 48
  starting_unit      = 1
  outer_width        = 600
  outer_depth        = 1200
  outer_unit         = "mm"
  mounting_depth     = 800
  weight             = 150.5
  weight_unit        = "kg"
  max_weight         = 1500
  description        = "A rack type for testing"
  comments           = "This rack type was created for integration testing."
}
