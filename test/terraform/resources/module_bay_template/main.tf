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
  name = "Test Manufacturer for Module Bay Template"
  slug = "test-mfg-mbt"
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = "Test Device Type Module Bay Template"
  slug         = "test-dt-mbt"
}

# Test 1: Basic module bay template
resource "netbox_module_bay_template" "basic" {
  device_type = netbox_device_type.test.id
  name        = "Module Bay 1"
}

# Test 2: Module bay template with all optional fields
resource "netbox_module_bay_template" "complete" {
  device_type = netbox_device_type.test.id
  name        = "Module Bay 2"
  label       = "MB-2"
  position    = "A"
  description = "Test module bay template with full details"
}
