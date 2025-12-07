# Circuit Termination Data Source Test

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
resource "netbox_site" "test" {
  name   = "Test Site for Circuit Term DS"
  slug   = "test-site-circuit-term-ds"
  status = "active"
}

resource "netbox_provider" "test" {
  name = "Test Provider for Circuit Term DS"
  slug = "test-provider-circuit-term-ds"
}

resource "netbox_circuit_type" "test" {
  name = "Test Circuit Type for Circuit Term DS"
  slug = "test-circuit-type-circuit-term-ds"
}

resource "netbox_circuit" "test" {
  cid              = "CKT-TERM-DS-001"
  circuit_provider = netbox_provider.test.id
  type             = netbox_circuit_type.test.id
  status           = "active"
}

resource "netbox_circuit_termination" "test" {
  circuit     = netbox_circuit.test.id
  term_side   = "A"
  site        = netbox_site.test.id
  description = "Test circuit termination for data source"
}

# Test: Lookup circuit termination by ID
data "netbox_circuit_termination" "by_id" {
  id = netbox_circuit_termination.test.id
}
