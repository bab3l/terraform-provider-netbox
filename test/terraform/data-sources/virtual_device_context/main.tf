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

# Dependencies
resource "netbox_site" "test" {
  name   = "Test Site for VDC DS"
  slug   = "test-site-vdc-ds"
  status = "active"
}

resource "netbox_manufacturer" "test" {
  name = "Test Manufacturer for VDC DS"
  slug = "test-mfg-vdc-ds"
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = "Test Device Type VDC DS"
  slug         = "test-dt-vdc-ds"
}

resource "netbox_device_role" "test" {
  name  = "Test Role for VDC DS"
  slug  = "test-role-vdc-ds"
  color = "aa1409"
}

resource "netbox_device" "test" {
  name        = "test-device-vdc-ds"
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
  status      = "active"
}

resource "netbox_virtual_device_context" "test" {
  name        = "test-vdc-ds"
  device      = netbox_device.test.id
  status      = "active"
  description = "Virtual device context for data source test"
}

data "netbox_virtual_device_context" "by_id" {
  id = netbox_virtual_device_context.test.id
}
