# Provider Account Data Source Test

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
  name = "Test Provider for Account DS"
  slug = "test-provider-account-ds"
}

resource "netbox_provider_account" "test" {
  name             = "Test Provider Account DS"
  circuit_provider = netbox_provider.test.id
  account          = "ACCT-DS-001"
  description      = "Test provider account for data source"
}

# Test: Lookup provider account by ID
data "netbox_provider_account" "by_id" {
  id = netbox_provider_account.test.id
}
