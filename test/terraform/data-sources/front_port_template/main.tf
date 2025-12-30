# Front Port Template Data Source Test

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
  name = "Test Manufacturer for Front Port Template DS"
  slug = "test-manufacturer-front-port-template-ds"
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = "Test Model for Front Port Template DS"
  slug         = "test-model-front-port-template-ds"
}

resource "netbox_rear_port_template" "test" {
  name        = "RearPortTpl0-DS"
  device_type = netbox_device_type.test.id
  type        = "8p8c"
  positions   = 2
}

resource "netbox_front_port_template" "test" {
  name        = "FrontPortTpl0-DS"
  device_type = netbox_device_type.test.id
  type        = "8p8c"
  rear_port   = netbox_rear_port_template.test.name
  description = "Test front port template for data source"
}

# Test: Lookup front port template by ID
data "netbox_front_port_template" "by_id" {
  id = netbox_front_port_template.test.id
}
