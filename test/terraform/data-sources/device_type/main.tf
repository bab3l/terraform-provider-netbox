# Device Type Data Source Test

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
resource "netbox_manufacturer" "test" {
  name = "Test Manufacturer for Device Type DS"
  slug = "test-manufacturer-device-type-ds"
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = "Test Model DS"
  slug         = "test-model-ds"
  description  = "Test device type for data source"
}

# Test: Lookup device type by ID
data "netbox_device_type" "by_id" {
  id = netbox_device_type.test.id
}

# Test: Lookup device type by model (slug)
data "netbox_device_type" "by_slug" {
  slug = netbox_device_type.test.slug

  depends_on = [netbox_device_type.test]
}
