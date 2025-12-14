# Device Type Resource Test

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

# First, create required manufacturer
resource "netbox_manufacturer" "test" {
  name = "Device Type Test Manufacturer"
  slug = "device-type-test-manufacturer"
}

# Basic device type with required fields only
resource "netbox_device_type" "basic" {
  manufacturer = netbox_manufacturer.test.id
  model        = "Basic Test Model"
  slug         = "basic-test-model"
}

# Complete device type with all fields
resource "netbox_device_type" "complete" {
  manufacturer   = netbox_manufacturer.test.id
  model          = "Complete Test Model"
  slug           = "complete-test-model"
  part_number    = "CTM-001"
  u_height       = 2
  is_full_depth  = true
  description    = "A complete device type with all fields"

  depends_on = [netbox_manufacturer.test]
}

# 1U device type
resource "netbox_device_type" "server_1u" {
  manufacturer  = netbox_manufacturer.test.id
  model         = "1U Server"
  slug          = "1u-server"
  u_height      = 1
  is_full_depth = false
}
