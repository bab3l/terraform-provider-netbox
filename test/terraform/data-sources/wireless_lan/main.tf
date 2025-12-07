# Wireless LAN Data Source Test

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
resource "netbox_wireless_lan" "test" {
  ssid        = "TestWLAN-DS"
  status      = "active"
  description = "Test wireless LAN for data source"
}

# Test: Lookup wireless LAN by ID
data "netbox_wireless_lan" "by_id" {
  id = netbox_wireless_lan.test.id
}
