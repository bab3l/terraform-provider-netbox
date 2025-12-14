# IP Range Data Source Test

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
resource "netbox_ip_range" "test" {
  start_address = "192.168.200.10/24"
  end_address   = "192.168.200.50/24"
  status        = "active"
  description   = "Test IP range for data source"
}

# Test: Lookup IP range by ID
data "netbox_ip_range" "by_id" {
  id = netbox_ip_range.test.id
}
