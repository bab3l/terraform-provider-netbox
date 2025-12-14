# Virtual Chassis Resource Test

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

# Test 1: Basic virtual chassis creation
resource "netbox_virtual_chassis" "basic" {
  name = "Test Virtual Chassis Basic"
}

# Test 2: Virtual chassis with all optional fields
resource "netbox_virtual_chassis" "complete" {
  name        = "Test Virtual Chassis Complete"
  domain      = "example.local"
  description = "A virtual chassis for testing"
  comments    = "This virtual chassis was created for integration testing."
}
