# Provider Network Data Source Test

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
  name = "Test Provider for Network DS"
  slug = "test-provider-network-ds"
}

resource "netbox_provider_network" "test" {
  name             = "Test Provider Network DS"
  circuit_provider = netbox_provider.test.id
  description      = "Test provider network for data source"
}

# Test: Lookup provider network by ID
data "netbox_provider_network" "by_id" {
  id = netbox_provider_network.test.id
}
