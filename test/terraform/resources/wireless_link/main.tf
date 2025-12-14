# Wireless Link Resource Test
# Tests the netbox_wireless_link resource by creating wireless links between interfaces

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

# Prerequisites: Site, Manufacturer, Device Type, Device Role, Devices, and Wireless Interfaces
# Wireless links connect wireless interfaces on different devices

resource "netbox_site" "test" {
  name        = "Wireless Link Test Site"
  slug        = "wireless-link-test-site"
  description = "Site for wireless link testing"
}

resource "netbox_manufacturer" "test" {
  name        = "Wireless Link Test Manufacturer"
  slug        = "wireless-link-test-manufacturer"
  description = "Manufacturer for wireless link testing"
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = "Wireless Link Test Model"
  slug         = "wireless-link-test-model"
  description  = "Device type for wireless link testing"
}

resource "netbox_device_role" "test" {
  name        = "Wireless Link Test Device Role"
  slug        = "wireless-link-test-device-role"
  description = "Device role for wireless link testing"
}

resource "netbox_tenant" "test" {
  name        = "Wireless Link Test Tenant"
  slug        = "wireless-link-test-tenant"
  description = "Tenant for wireless link testing"
}

# Device A - source of wireless links
resource "netbox_device" "device_a" {
  name        = "wireless-link-test-device-a"
  site        = netbox_site.test.id
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  description = "Device A for wireless link testing"
}

# Device B - destination of wireless links
resource "netbox_device" "device_b" {
  name        = "wireless-link-test-device-b"
  site        = netbox_site.test.id
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  description = "Device B for wireless link testing"
}

# Wireless Interfaces on Device A
resource "netbox_interface" "device_a_wlan0" {
  device      = netbox_device.device_a.id
  name        = "wlan0"
  type        = "ieee802.11a"
  description = "Device A wlan0"
}

resource "netbox_interface" "device_a_wlan1" {
  device      = netbox_device.device_a.id
  name        = "wlan1"
  type        = "ieee802.11ac"
  description = "Device A wlan1"
}

# Wireless Interfaces on Device B
resource "netbox_interface" "device_b_wlan0" {
  device      = netbox_device.device_b.id
  name        = "wlan0"
  type        = "ieee802.11a"
  description = "Device B wlan0"
}

resource "netbox_interface" "device_b_wlan1" {
  device      = netbox_device.device_b.id
  name        = "wlan1"
  type        = "ieee802.11ac"
  description = "Device B wlan1"
}

# Test 1: Basic wireless link
resource "netbox_wireless_link" "basic" {
  interface_a = netbox_interface.device_a_wlan0.id
  interface_b = netbox_interface.device_b_wlan0.id
}

# Test 2: Wireless link with SSID and status
resource "netbox_wireless_link" "with_ssid" {
  interface_a = netbox_interface.device_a_wlan1.id
  interface_b = netbox_interface.device_b_wlan1.id
  ssid        = "TestLink-Network"
  status      = "connected"
  description = "Wireless link with SSID"
}

# Test 3: Wireless link with all optional fields (will need additional interfaces)
resource "netbox_interface" "device_a_wlan2" {
  device      = netbox_device.device_a.id
  name        = "wlan2"
  type        = "ieee802.11ax"
  description = "Device A wlan2"
}

resource "netbox_interface" "device_b_wlan2" {
  device      = netbox_device.device_b.id
  name        = "wlan2"
  type        = "ieee802.11ax"
  description = "Device B wlan2"
}

resource "netbox_wireless_link" "complete" {
  interface_a   = netbox_interface.device_a_wlan2.id
  interface_b   = netbox_interface.device_b_wlan2.id
  ssid          = "TestLink-Complete"
  status        = "planned"
  tenant        = netbox_tenant.test.id
  auth_type     = "wpa-personal"
  auth_cipher   = "aes"
  auth_psk      = "testpassword123"
  distance      = 1.5
  distance_unit = "km"
  description   = "A complete wireless link for testing"
  comments      = "This wireless link was created for integration testing."
}
