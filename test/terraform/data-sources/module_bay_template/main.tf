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
  name = "Test Manufacturer for Module Bay Template DS"
  slug = "test-mfg-mbt-ds"
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = "Test Device Type Module Bay Template DS"
  slug         = "test-dt-mbt-ds"
}

resource "netbox_module_bay_template" "test" {
  device_type = netbox_device_type.test.id
  name        = "Module Bay DS"
  label       = "MB-DS"
  position    = "B"
  description = "Module bay template for data source test"
}

data "netbox_module_bay_template" "by_id" {
  id = netbox_module_bay_template.test.id
}
