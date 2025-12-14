# RIR Data Source Test

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
resource "netbox_rir" "test" {
  name        = "Test RIR DS"
  slug        = "test-rir-ds"
  description = "Test RIR for data source"
}

# Test: Lookup RIR by ID
data "netbox_rir" "by_id" {
  id = netbox_rir.test.id
}

# Test: Lookup RIR by name
data "netbox_rir" "by_name" {
  name = netbox_rir.test.name

  depends_on = [netbox_rir.test]
}
