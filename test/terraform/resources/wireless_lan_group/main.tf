# Wireless LAN Group Resource Test

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

# Test 1: Basic wireless LAN group creation
resource "netbox_wireless_lan_group" "basic" {
  name = "Test Wireless LAN Group Basic"
  slug = "test-wireless-lan-group-basic"
}

# Test 2: Wireless LAN group with parent (nested group)
resource "netbox_wireless_lan_group" "parent" {
  name = "Test Wireless LAN Group Parent"
  slug = "test-wireless-lan-group-parent"
}

resource "netbox_wireless_lan_group" "child" {
  name        = "Test Wireless LAN Group Child"
  slug        = "test-wireless-lan-group-child"
  parent      = netbox_wireless_lan_group.parent.id
  description = "A wireless LAN group for testing"
}
