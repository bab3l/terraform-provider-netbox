# Power Port Template Data Source Test

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
  name = "Test Manufacturer for Power Port Template DS"
  slug = "test-manufacturer-power-port-template-ds"
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = "Test Model for Power Port Template DS"
  slug         = "test-model-power-port-template-ds"
}

resource "netbox_power_port_template" "test" {
  name        = "PSU-Tpl1-DS"
  device_type = netbox_device_type.test.id
  description = "Test power port template for data source"
}

# Test: Lookup power port template by ID
data "netbox_power_port_template" "by_id" {
  id = netbox_power_port_template.test.id
}
