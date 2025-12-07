# Circuit Termination Resource Test

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
  name = "Test Provider for Circuit Term"
  slug = "test-provider-circuit-term"
}

resource "netbox_circuit_type" "test" {
  name = "Test Circuit Type"
  slug = "test-circuit-type-term"
}

resource "netbox_circuit" "test" {
  cid              = "TEST-CIRCUIT-001"
  circuit_provider = netbox_provider.test.id
  type             = netbox_circuit_type.test.id
}

resource "netbox_site" "test" {
  name   = "Test Site for Circuit Term"
  slug   = "test-site-circuit-term"
  status = "active"
}

# Test 1: Basic circuit termination (Side A)
resource "netbox_circuit_termination" "basic" {
  circuit   = netbox_circuit.test.id
  term_side = "A"
  site      = netbox_site.test.id
}

# Test 2: Circuit termination with all optional fields (Side Z)
resource "netbox_circuit_termination" "complete" {
  circuit     = netbox_circuit.test.id
  term_side   = "Z"
  site        = netbox_site.test.id
  port_speed  = 1000000
  description = "Z-side termination for testing"
}
