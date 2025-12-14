# Device Bay Template Resource Test
# Tests creation and management of device bay templates in Netbox

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
  name        = "Test Manufacturer for Device Bay Template"
  slug        = "test-mfr-device-bay-template"
  description = "Manufacturer for device bay template testing"
}

# Create device type that will have device bays
resource "netbox_device_type" "test" {
  manufacturer    = netbox_manufacturer.test.id
  model           = "Test Chassis for Device Bay Templates"
  slug            = "test-chassis-device-bay-template"
  description     = "Test chassis device type"
  u_height        = 4
  is_full_depth   = true
  subdevice_role  = "parent"
}

# Create basic device bay template
resource "netbox_device_bay_template" "basic" {
  device_type = netbox_device_type.test.id
  name        = "Bay 1"
}

# Create device bay template with all optional fields
resource "netbox_device_bay_template" "full" {
  device_type = netbox_device_type.test.id
  name        = "Bay 2"
  label       = "SLOT-2"
  description = "Second device bay slot with all options"
}

# Create multiple device bay templates for a chassis
resource "netbox_device_bay_template" "multi" {
  count       = 4
  device_type = netbox_device_type.test.id
  name        = "Expansion Bay ${count.index + 1}"
  label       = "EXP-${count.index + 1}"
  description = "Expansion bay ${count.index + 1}"
}
