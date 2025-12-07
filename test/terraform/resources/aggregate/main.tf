# Aggregate Resource Test

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

# First create the required RIR
resource "netbox_rir" "test" {
  name = "Test RIR for Aggregate"
  slug = "test-rir-aggregate"
}

# Test 1: Basic aggregate creation
resource "netbox_aggregate" "basic" {
  prefix = "10.0.0.0/8"
  rir    = netbox_rir.test.id
}

# Test 2: Aggregate with all optional fields
resource "netbox_aggregate" "complete" {
  prefix      = "172.16.0.0/12"
  rir         = netbox_rir.test.id
  description = "Private network aggregate for testing"
  comments    = "This aggregate was created for integration testing."
}
