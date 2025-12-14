// Platform Resource Test

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

# Create a manufacturer for the platform
resource "netbox_manufacturer" "for_platform" {
  name = "Platform Test Manufacturer"
  slug = "platform-test-manufacturer"
}

resource "netbox_platform" "basic" {
  name         = "Test Platform"
  slug         = "test-platform"
  manufacturer = netbox_manufacturer.for_platform.id
  description  = "Platform created for Terraform integration tests"
}
