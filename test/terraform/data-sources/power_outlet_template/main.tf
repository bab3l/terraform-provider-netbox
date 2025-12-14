# Power Outlet Template Data Source Test

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
resource "netbox_manufacturer" "test" {
  name = "Test Manufacturer for Power Outlet Template DS"
  slug = "test-manufacturer-power-outlet-template-ds"
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = "Test Model for Power Outlet Template DS"
  slug         = "test-model-power-outlet-template-ds"
}

resource "netbox_power_outlet_template" "test" {
  name        = "OutletTpl1-DS"
  device_type = netbox_device_type.test.id
  description = "Test power outlet template for data source"
}

# Test: Lookup power outlet template by ID
data "netbox_power_outlet_template" "by_id" {
  id = netbox_power_outlet_template.test.id
}
