# Circuit Type Integration Test
# Tests the netbox_circuit_type resource with basic and complete configurations

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

# Basic Circuit Type with only required fields
resource "netbox_circuit_type" "basic" {
  name = "Basic Test Circuit Type"
  slug = "basic-test-circuit-type"
}

# Complete Circuit Type with all optional fields
resource "netbox_circuit_type" "complete" {
  name        = "Complete Test Circuit Type"
  slug        = "complete-test-circuit-type"
  description = "Complete circuit type for integration testing"
  color       = "ff5722"
}

# Common circuit types for testing
resource "netbox_circuit_type" "internet" {
  name        = "Internet"
  slug        = "internet"
  description = "Internet connectivity circuit"
}

resource "netbox_circuit_type" "mpls" {
  name        = "MPLS"
  slug        = "mpls"
  description = "MPLS VPN circuit"
}
