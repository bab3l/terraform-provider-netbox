# Provider Data Source Test

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
resource "netbox_provider" "test" {
  name        = "Test Provider DS"
  slug        = "test-provider-ds"
  description = "Test provider for data source"
}

# Test: Lookup provider by ID
data "netbox_provider" "by_id" {
  id = netbox_provider.test.id
}

# Test: Lookup provider by name
data "netbox_provider" "by_name" {
  name = netbox_provider.test.name

  depends_on = [netbox_provider.test]
}
