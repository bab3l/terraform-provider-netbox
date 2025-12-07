# Provider Network Resource Test

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
resource "netbox_provider" "test" {
  name = "Test Provider for Network"
  slug = "test-provider-network"
}

# Test 1: Basic provider network creation
resource "netbox_provider_network" "basic" {
  name             = "Test Provider Network Basic"
  circuit_provider = netbox_provider.test.id
}

# Test 2: Provider network with all optional fields
resource "netbox_provider_network" "complete" {
  name             = "Test Provider Network Complete"
  circuit_provider = netbox_provider.test.id
  service_id       = "SVC-12345"
  description      = "A provider network for testing"
  comments         = "This provider network was created for integration testing."
}
