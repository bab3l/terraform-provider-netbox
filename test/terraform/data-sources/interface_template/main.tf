# Interface Template Data Source Test

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
  name = "Test Manufacturer for Interface Template DS"
  slug = "test-manufacturer-interface-template-ds"
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = "Test Model for Interface Template DS"
  slug         = "test-model-interface-template-ds"
}

resource "netbox_interface_template" "test" {
  name        = "eth0-DS"
  device_type = netbox_device_type.test.id
  type        = "1000base-t"
  description = "Test interface template for data source"
}

# Test: Lookup interface template by ID
data "netbox_interface_template" "by_id" {
  id = netbox_interface_template.test.id
}
