# IP Range Resource Test

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

# Test 1: Basic IP range creation
resource "netbox_ip_range" "basic" {
  start_address = "10.0.0.1/24"
  end_address   = "10.0.0.100/24"
}

# Test 2: IP range with all optional fields
resource "netbox_ip_range" "complete" {
  start_address = "192.168.1.100/24"
  end_address   = "192.168.1.200/24"
  status        = "active"
  description   = "An IP range for testing"
  comments      = "This IP range was created for integration testing."
}
