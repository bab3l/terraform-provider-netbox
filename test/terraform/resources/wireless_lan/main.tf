# Wireless LAN Resource Test

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
resource "netbox_wireless_lan_group" "test" {
  name = "Test WLAN Group for Wireless LAN"
  slug = "test-wlan-group-wireless-lan"
}

# Test 1: Basic wireless LAN creation
resource "netbox_wireless_lan" "basic" {
  ssid   = "TestWLAN-Basic"
  status = "active"
}

# Test 2: Wireless LAN with all optional fields
resource "netbox_wireless_lan" "complete" {
  ssid        = "TestWLAN-Complete"
  group       = netbox_wireless_lan_group.test.id
  status      = "active"
  description = "A wireless LAN for testing"
  comments    = "This wireless LAN was created for integration testing."
  auth_type   = "wpa-personal"
  auth_cipher = "aes"
  auth_psk    = "testpassword123"
}
