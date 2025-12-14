# Wireless Link Data Source Test
# Tests the netbox_wireless_link data source

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

# Prerequisites: Create a wireless link to look up
resource "netbox_site" "test" {
  name        = "Wireless Link DS Test Site"
  slug        = "wireless-link-ds-test-site"
  description = "Site for wireless link data source testing"
}

resource "netbox_manufacturer" "test" {
  name        = "Wireless Link DS Test Manufacturer"
  slug        = "wireless-link-ds-test-manufacturer"
  description = "Manufacturer for wireless link data source testing"
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = "Wireless Link DS Test Model"
  slug         = "wireless-link-ds-test-model"
  description  = "Device type for wireless link data source testing"
}

resource "netbox_device_role" "test" {
  name        = "Wireless Link DS Test Device Role"
  slug        = "wireless-link-ds-test-device-role"
  description = "Device role for wireless link data source testing"
}

resource "netbox_device" "device_a" {
  name        = "wireless-link-ds-test-device-a"
  site        = netbox_site.test.id
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  description = "Device A for wireless link data source testing"
}

resource "netbox_device" "device_b" {
  name        = "wireless-link-ds-test-device-b"
  site        = netbox_site.test.id
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  description = "Device B for wireless link data source testing"
}

resource "netbox_interface" "device_a_wlan0" {
  device      = netbox_device.device_a.id
  name        = "wlan0"
  type        = "ieee802.11ax"
  description = "Device A wlan0"
}

resource "netbox_interface" "device_b_wlan0" {
  device      = netbox_device.device_b.id
  name        = "wlan0"
  type        = "ieee802.11ax"
  description = "Device B wlan0"
}

resource "netbox_wireless_link" "test" {
  interface_a   = netbox_interface.device_a_wlan0.id
  interface_b   = netbox_interface.device_b_wlan0.id
  ssid          = "TestLink-DataSource"
  status        = "connected"
  description   = "Wireless link for data source testing"
  comments      = "This wireless link was created for data source integration testing."
}

# Data source test: lookup by ID
data "netbox_wireless_link" "by_id" {
  id = netbox_wireless_link.test.id
}
