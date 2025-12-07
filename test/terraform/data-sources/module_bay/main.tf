# Module Bay Data Source Test

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
  name   = "Test Site for Module Bay DS"
  slug   = "test-site-module-bay-ds"
  status = "active"
}

resource "netbox_manufacturer" "test" {
  name = "Test Manufacturer for Module Bay DS"
  slug = "test-manufacturer-module-bay-ds"
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = "Test Model for Module Bay DS"
  slug         = "test-model-module-bay-ds"
}

resource "netbox_device_role" "test" {
  name  = "Test Role for Module Bay DS"
  slug  = "test-role-module-bay-ds"
  color = "aabbcc"
}

resource "netbox_device" "test" {
  name        = "Test Device for Module Bay DS"
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
  status      = "active"
}

resource "netbox_module_bay" "test" {
  name        = "Slot1-DS"
  device      = netbox_device.test.id
  description = "Test module bay for data source"
}

# Test: Lookup module bay by ID
data "netbox_module_bay" "by_id" {
  id = netbox_module_bay.test.id
}
