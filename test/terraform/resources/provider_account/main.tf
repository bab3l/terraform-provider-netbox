# Provider Account Resource Test

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
  name = "Test Provider for Account"
  slug = "test-provider-account"
}

# Test 1: Basic provider account creation
resource "netbox_provider_account" "basic" {
  name             = "Test Provider Account Basic"
  circuit_provider = netbox_provider.test.id
  account          = "BASIC-001"
}

# Test 2: Provider account with all optional fields
resource "netbox_provider_account" "complete" {
  name             = "Test Provider Account Complete"
  circuit_provider = netbox_provider.test.id
  account          = "ACCT-12345"
  description      = "A provider account for testing"
  comments         = "This provider account was created for integration testing."
}
