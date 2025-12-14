// Manufacturer Resource Test

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

resource "netbox_manufacturer" "basic" {
  name        = "Test Manufacturer"
  slug        = "test-manufacturer"
  description = "Manufacturer created for Terraform integration tests"
}
