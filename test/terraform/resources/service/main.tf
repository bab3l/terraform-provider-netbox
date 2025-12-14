# Service Resource Test

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
  name   = "Test Site for Service"
  slug   = "test-site-service"
  status = "active"
}

resource "netbox_manufacturer" "test" {
  name = "Test Manufacturer for Service"
  slug = "test-manufacturer-service"
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = "Test Model for Service"
  slug         = "test-model-service"
}

resource "netbox_device_role" "test" {
  name  = "Test Role for Service"
  slug  = "test-role-service"
  color = "aabbcc"
}

resource "netbox_device" "test" {
  name        = "Test Device for Service"
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
  status      = "active"
}

# Test 1: Basic service creation on device
resource "netbox_service" "basic" {
  name     = "Test Service Basic"
  device   = netbox_device.test.id
  ports    = [80]
  protocol = "tcp"
}

# Test 2: Service with all optional fields
resource "netbox_service" "complete" {
  name        = "Test Service Complete"
  device      = netbox_device.test.id
  ports       = [443, 8443]
  protocol    = "tcp"
  description = "A service for testing"
  comments    = "This service was created for integration testing."
}
