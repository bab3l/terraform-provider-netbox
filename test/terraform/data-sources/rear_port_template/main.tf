# Rear Port Template Data Source Test

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
  name = "Test Manufacturer for Rear Port Template DS"
  slug = "test-manufacturer-rear-port-template-ds"
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = "Test Model for Rear Port Template DS"
  slug         = "test-model-rear-port-template-ds"
}

resource "netbox_rear_port_template" "test" {
  name        = "RearPortTpl0-DS"
  device_type = netbox_device_type.test.id
  type        = "8p8c"
  positions   = 2
  description = "Test rear port template for data source"
}

# Test: Lookup rear port template by ID
data "netbox_rear_port_template" "by_id" {
  id = netbox_rear_port_template.test.id
}
