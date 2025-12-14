# Device Bay Template Data Source Test
# Tests retrieving device bay templates from Netbox

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

# Create manufacturer for device type
resource "netbox_manufacturer" "test" {
  name        = "Test Manufacturer for Device Bay Template DS"
  slug        = "test-mfr-device-bay-template-ds"
  description = "Manufacturer for device bay template data source testing"
}

# Create device type
resource "netbox_device_type" "test" {
  manufacturer   = netbox_manufacturer.test.id
  model          = "Test Chassis for Device Bay DS"
  slug           = "test-chassis-device-bay-ds"
  u_height       = 4
  subdevice_role = "parent"
}

# Create device bay template to lookup
resource "netbox_device_bay_template" "test" {
  device_type = netbox_device_type.test.id
  name        = "Test Bay for DS"
  label       = "BAY-DS"
  description = "Device bay template for data source testing"
}

# Look up by ID
data "netbox_device_bay_template" "by_id" {
  id = netbox_device_bay_template.test.id
}

# Look up by name with device_type
data "netbox_device_bay_template" "by_name" {
  name        = netbox_device_bay_template.test.name
  device_type = netbox_device_type.test.id
}
