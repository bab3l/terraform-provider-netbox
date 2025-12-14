# Role (IPAM Role) Resource Test

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

# Test 1: Basic role creation
resource "netbox_role" "basic" {
  name = "Test IPAM Role Basic"
  slug = "test-ipam-role-basic"
}

# Test 2: Role with all optional fields
resource "netbox_role" "complete" {
  name        = "Test IPAM Role Complete"
  slug        = "test-ipam-role-complete"
  weight      = 500
  description = "An IPAM role for testing"
}
