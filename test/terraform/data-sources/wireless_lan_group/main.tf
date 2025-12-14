# Wireless LAN Group Data Source Test

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
resource "netbox_wireless_lan_group" "test" {
  name        = "Test WLAN Group DS"
  slug        = "test-wlan-group-ds"
  description = "Test wireless LAN group for data source"
}

# Test: Lookup wireless LAN group by ID
data "netbox_wireless_lan_group" "by_id" {
  id = netbox_wireless_lan_group.test.id
}

# Test: Lookup wireless LAN group by name
data "netbox_wireless_lan_group" "by_name" {
  name = netbox_wireless_lan_group.test.name

  depends_on = [netbox_wireless_lan_group.test]
}
