# Module Type Resource Test

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
  name = "Test Manufacturer for Module Type"
  slug = "test-mfg-module-type"
}

# Test 1: Basic module type creation
resource "netbox_module_type" "basic" {
  manufacturer = netbox_manufacturer.test.id
  model        = "Test Module Type Basic"
}

# Test 2: Module type with all optional fields
resource "netbox_module_type" "complete" {
  manufacturer = netbox_manufacturer.test.id
  model        = "Test Module Type Complete"
  part_number  = "MT-COMPLETE-001"
  description  = "A module type for testing"
  comments     = "This module type was created for integration testing."
}
