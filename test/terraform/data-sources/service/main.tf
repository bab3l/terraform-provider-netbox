# Service Data Source Test

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
  name   = "Test Site for Service DS"
  slug   = "test-site-service-ds"
  status = "active"
}

resource "netbox_manufacturer" "test" {
  name = "Test Manufacturer for Service DS"
  slug = "test-manufacturer-service-ds"
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = "Test Model for Service DS"
  slug         = "test-model-service-ds"
}

resource "netbox_device_role" "test" {
  name  = "Test Role for Service DS"
  slug  = "test-role-service-ds"
  color = "aabbcc"
}

resource "netbox_device" "test" {
  name        = "Test Device for Service DS"
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
  status      = "active"
}

resource "netbox_service" "test" {
  name        = "Test Service DS"
  device      = netbox_device.test.id
  ports       = [80]
  protocol    = "tcp"
  description = "Test service for data source"
}

# Test: Lookup service by ID
data "netbox_service" "by_id" {
  id = netbox_service.test.id
}
