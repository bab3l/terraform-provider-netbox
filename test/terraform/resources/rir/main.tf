# RIR (Regional Internet Registry) Resource Test

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

# Test 1: Basic RIR creation
resource "netbox_rir" "basic" {
  name = "Test RIR Basic"
  slug = "test-rir-basic"
}

# Test 2: RIR with all optional fields
resource "netbox_rir" "complete" {
  name        = "Test RIR Complete"
  slug        = "test-rir-complete"
  is_private  = true
  description = "A RIR for testing"
}
