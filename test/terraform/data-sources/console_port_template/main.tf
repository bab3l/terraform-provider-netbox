# Console Port Template Data Source Test

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
  name = "Test Manufacturer for Console Port Template DS"
  slug = "test-manufacturer-console-port-template-ds"
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = "Test Model for Console Port Template DS"
  slug         = "test-model-console-port-template-ds"
}

resource "netbox_console_port_template" "test" {
  name        = "ConsoleTpl0-DS"
  device_type = netbox_device_type.test.id
  description = "Test console port template for data source"
}

# Test: Lookup console port template by ID
data "netbox_console_port_template" "by_id" {
  id = netbox_console_port_template.test.id
}
