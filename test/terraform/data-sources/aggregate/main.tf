# Aggregate Data Source Test

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

# Dependencies - create resources to test data sources
resource "netbox_rir" "test" {
  name = "Test RIR for Aggregate DS"
  slug = "test-rir-aggregate-ds"
}

resource "netbox_aggregate" "test" {
  prefix      = "192.168.0.0/16"
  rir         = netbox_rir.test.id
  description = "Test aggregate for data source"
}

# Test: Lookup aggregate by ID
data "netbox_aggregate" "by_id" {
  id = netbox_aggregate.test.id
}

# Test: Lookup aggregate by prefix
data "netbox_aggregate" "by_prefix" {
  prefix = netbox_aggregate.test.prefix

  depends_on = [netbox_aggregate.test]
}
