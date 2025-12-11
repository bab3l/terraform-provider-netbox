# Front Port Template Resource Test

provider "netbox" {
  # Uses NETBOX_SERVER_URL and NETBOX_API_TOKEN environment variables
}

# Dependencies
resource "netbox_manufacturer" "test" {
  name = "Test Manufacturer for Front Port Template"
  slug = "test-mfg-front-port-tpl"
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = "Test Device Type Front Port Template"
  slug         = "test-dt-front-port-tpl"
}

resource "netbox_rear_port_template" "test" {
  device_type = netbox_device_type.test.id
  name        = "rear0"
  type        = "8p8c"
  positions   = 4
}

# Test 1: Basic front port template creation
resource "netbox_front_port_template" "basic" {
  device_type = netbox_device_type.test.id
  name        = "front0"
  type        = "8p8c"
  rear_port   = netbox_rear_port_template.test.name
}

# Test 2: Front port template with all optional fields
resource "netbox_front_port_template" "complete" {
  device_type        = netbox_device_type.test.id
  name               = "front1"
  type               = "8p8c"
  rear_port          = netbox_rear_port_template.test.name
  rear_port_position = 2
  label              = "Front Port Template 1"
  color              = "aa1409"
  description        = "Front port template for testing"
}
