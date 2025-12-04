# Device Integration Test
# Tests the netbox_device resource with basic and complete configurations

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
resource "netbox_site" "test" {
  name        = "Device Test Site"
  slug        = "device-test-site"
  status      = "active"
  description = "Site for device testing"
}

resource "netbox_manufacturer" "test" {
  name        = "Device Test Manufacturer"
  slug        = "device-test-manufacturer"
  description = "Manufacturer for device testing"
}

resource "netbox_device_type" "test" {
  model        = "Test Device Type"
  slug         = "test-device-type"
  manufacturer = netbox_manufacturer.test.id
  u_height     = 1
  description  = "Device type for device testing"
}

resource "netbox_device_role" "test" {
  name        = "Device Test Role"
  slug        = "device-test-role"
  color       = "00ff00"
  description = "Device role for device testing"
}

resource "netbox_tenant" "test" {
  name = "Device Test Tenant"
  slug = "device-test-tenant"
}

resource "netbox_platform" "test" {
  name        = "Device Test Platform"
  slug        = "device-test-platform"
  description = "Platform for device testing"
}

resource "netbox_rack" "test" {
  name        = "Device Test Rack"
  site        = netbox_site.test.id
  status      = "active"
  u_height    = 42
  description = "Rack for device testing"
}

# Basic Device with only required fields
resource "netbox_device" "basic" {
  name        = "Basic Test Device"
  site        = netbox_site.test.id
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
}

# Complete Device with all optional fields
resource "netbox_device" "complete" {
  name        = "Complete Test Device"
  site        = netbox_site.test.id
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  status      = "active"
  description = "Complete device for integration testing"
  comments    = "Created by terraform integration test"
  tenant      = netbox_tenant.test.id
  platform    = netbox_platform.test.id
  serial      = "SN-COMPLETE-001"
  asset_tag   = "ASSET-001"
}

# Device in rack
resource "netbox_device" "racked" {
  name        = "Racked Test Device"
  site        = netbox_site.test.id
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  status      = "active"
  description = "Device in rack"
  rack        = netbox_rack.test.id
  position    = 1
  face        = "front"
}

# Planned Device
resource "netbox_device" "planned" {
  name        = "Planned Test Device"
  site        = netbox_site.test.id
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  status      = "planned"
  description = "Planned device test"
}
