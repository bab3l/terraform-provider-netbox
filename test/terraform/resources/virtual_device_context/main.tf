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
resource "netbox_site" "test" {
  name   = "Test Site for Virtual Device Context"
  slug   = "test-site-vdc"
  status = "active"
}

resource "netbox_manufacturer" "test" {
  name = "Test Manufacturer for VDC"
  slug = "test-mfg-vdc"
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = "Test Device Type VDC"
  slug         = "test-dt-vdc"
}

resource "netbox_device_role" "test" {
  name  = "Test Role for VDC"
  slug  = "test-role-vdc"
  color = "aa1409"
}

resource "netbox_device" "test" {
  name        = "test-device-vdc"
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
  status      = "active"
}

# Test 1: Basic virtual device context
resource "netbox_virtual_device_context" "basic" {
  name   = "test-vdc-basic"
  device = netbox_device.test.id
  status = "active"
}

# Test 2: Virtual device context with all fields
resource "netbox_virtual_device_context" "complete" {
  name        = "test-vdc-complete"
  device      = netbox_device.test.id
  status      = "active"
  description = "Test VDC with full details"
}
