# IP Address Data Source Test

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
resource "netbox_ip_address" "test" {
  address     = "192.168.100.1/24"
  status      = "active"
  description = "Test IP address for data source"
}

# Test: Lookup IP address by ID
data "netbox_ip_address" "by_id" {
  id = netbox_ip_address.test.id
}

# Test: Lookup IP address by address
data "netbox_ip_address" "by_address" {
  address = netbox_ip_address.test.address

  depends_on = [netbox_ip_address.test]
}
