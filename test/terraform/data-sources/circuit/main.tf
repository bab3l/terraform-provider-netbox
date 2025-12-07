# Circuit Data Source Test

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
  name = "Test Provider for Circuit DS"
  slug = "test-provider-circuit-ds"
}

resource "netbox_circuit_type" "test" {
  name = "Test Circuit Type for Circuit DS"
  slug = "test-circuit-type-circuit-ds"
}

resource "netbox_circuit" "test" {
  cid              = "CKT-DS-001"
  circuit_provider = netbox_provider.test.id
  type             = netbox_circuit_type.test.id
  status           = "active"
  description      = "Test circuit for data source"
}

# Test: Lookup circuit by ID
data "netbox_circuit" "by_id" {
  id = netbox_circuit.test.id
}

# Test: Lookup circuit by CID
data "netbox_circuit" "by_cid" {
  cid = netbox_circuit.test.cid

  depends_on = [netbox_circuit.test]
}
